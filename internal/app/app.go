package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/andriykutsevol/DDDCasbinExample/configs"

	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/middleware"

	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/handler"
	"github.com/andriykutsevol/DDDCasbinExample/internal/presentation/http/router"

	"github.com/andriykutsevol/DDDCasbinExample/internal/app/application"

	//services
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/user"

	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/cognito"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/mongo"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/redis"

	//"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/cognito"

	pgstorage "github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/storage"
)


// Dependency Inversion Principle (DIP):
// This principle states that high-level modules should not depend on low-level modules. (The domain does not depend on anything - higest level)
// Instead, both should depend on abstractions (repositories, for example).
// If we violate this principle, we may have code that is difficult to test and reuse, and that is tightly coupled.
// This can also result in code that is difficult to maintain and extend.



func Run(configPath string) {

	var wg sync.WaitGroup

	//=========================================================================
	// Config
	//=========================================================================

	configs.MustLoad("../../configs/config.toml")
	//configs.PrintWithJSON()

	//=========================================================================
	// Repo Builder
	// Dependency injection for Infrastructure layer
	// Driven Adapters (Hexagonal Architecture )
	// Driven Adapter: Databases, Message Brokers, GRPC.
	// Repository implementation - adapter.
	//=========================================================================

	var infrarepos *infrastructure.InfraRepos
	infrarepos = &infrastructure.InfraRepos{}

	login_provider, ok := os.LookupEnv("LOGINPROVIDER")
	if !ok {
		log.Fatal("LOGINPROVIDER environment variable is not set")
	}

	if login_provider == "cognito"{
		var err error
		infrarepos, err = cognito.BuildRespositories(infrarepos)
		if err != nil {	
			log.Fatal("Cannot cognito.BuildRespositories():", err)
		}		
	} else if login_provider == "redis"{
		var err error
		infrarepos, err = redis.BuildRespositories(infrarepos)
		if err != nil {	
			log.Fatal("Cannot redis.BuildRespositories():", err)
		}
	} else {
		log.Fatal("Wrong login provider.")
	}
	
	db_type, ok := os.LookupEnv("DBTYPE")
	if !ok {
		log.Fatal("DBTYPE environment variable is not set")
	}	

	if db_type == "mongo" {
		var err error
		infrarepos, err = mongo.BuildRespositories(infrarepos)
		if err != nil {
			log.Fatal("Cannot mongo.BuildRespositories():", err)
		}

	} else if db_type == "pg" {

		var err error
		// https://medium.com/geekculture/work-with-go-postgresql-using-pgx-caee4573672

		// if we use docker run with --nerwork, we have to use container name (@template_go_react_pg)
		// if we use docker-compose, we have to use service name (@postgres).

		casbinDatabaseUri, ok := os.LookupEnv("PGCASBINURI")
		if !ok {
			log.Fatal("PGCASBINURI environment variable is not set")
		}
		casbinDbService, err := pgstorage.NewDatabaseService(casbinDatabaseUri)
		if err != nil {
			log.Fatal("Unable to connect to casbin database:", err)
		}
		defer casbinDbService.Pool.Close()

		weatherDatabaseUri, ok := os.LookupEnv("PGWEATHERURI")
		if !ok {
			log.Fatal("PGWEATHERURI environment variable is not set")
		}
		weatherDbService, err := pgstorage.NewDatabaseService(weatherDatabaseUri)
		if err != nil {
			log.Fatal("Unable to connect to weather database:", err)
		}
		defer casbinDbService.Pool.Close()

		infrarepos, err = pg.BuildRespositories(infrarepos, casbinDbService, weatherDbService)
		if err != nil {
			log.Fatal("Cannot pg.BuildRespositories():", err)
		}

	} else {
		log.Fatal("Error: db_type is wrong. It have to be either 'mongo' or 'pg'. ")
	}


	//=========================================================================
	// Dependency injection for Domain Services
	//=========================================================================

	userservice := user.NewService(
		infrarepos.AuthRepository,
		infrarepos.UserRepository,
		infrarepos.RoleRepository,
		infrarepos.UserRoleRepository)

	menuService := menu.NewService(
		infrarepos.MenuRepository,
		infrarepos.MenuActionRepository,
		infrarepos.MenuActionResourceRepository)


	//=========================================================================
	//Dependency injection for Application layer
	//=========================================================================

	// In Domain-Driven Design (DDD), the application layer is typically not considered an adapter.
	// Instead, it serves as a distinct layer responsible for coordinating interactions between the domain model and external systems or user interfaces.

	seedApplication := application.NewSeed(
		menuService,
		infrarepos.MenuActionRepository,
		infrarepos.RoleRepository,
		infrarepos.RoleMeuRepository,
		infrarepos.UserRepository,
		infrarepos.UserRoleRepository,
		infrarepos.WeatherRepository,
		infrarepos.AuthRepository,
	)
	// TODO: Whether to launch or not should be decided by the config file.
	err := seedApplication.Execute(context.TODO())
	if err != nil {
		log.Panic("ERROR: seedApplication.Execute(): ", err)
	}

	rbacAdapter := application.NewRbacAdapter(infrarepos.RbacRepository)

	syncedEnforcer, cleanup3, err := InitCasbin(rbacAdapter)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	_ = cleanup3

	roleApplication := application.NewRole(
		rbacAdapter,
		syncedEnforcer,
		infrarepos.RoleRepository,
		infrarepos.RoleMeuRepository,
		infrarepos.UserRepository)

	userApplication := application.NewUser(
		infrarepos.AuthRepository,
		rbacAdapter, syncedEnforcer,
		infrarepos.UserRepository,
		infrarepos.UserRoleRepository,
		infrarepos.RoleRepository)

	loginApplication := application.NewLogin(
		infrarepos.AuthRepository,
		infrarepos.UserRepository,
		infrarepos.RoleRepository,
		infrarepos.UserRoleRepository,
		userservice,
		infrarepos.MenuRepository,
		infrarepos.MenuActionRepository,
		infrarepos.RoleMeuRepository)

	menuApplication := application.NewMenu(menuService)

	demosApplication := application.NewDemos(
		infrarepos.UserRepository,
		infrarepos.UserRoleRepository)

	weatherApplication := application.NewWeather(infrarepos.WeatherRepository)


	//=========================================================================
	// Driving Adapters  (Hexagonal Architecture )
	// Dependency injection for Presentation layer
	// REST API
	// GRPC
	//=========================================================================

	// The handler recieves an application interface (driving port), but returns handler interface (which is visible in router)
	// and the handler is a driving adapter in this case (because it implements application interface)
	// But that driving port could also be implemented with GRPC (and that would be an another adapter)
	roleHandler := handler.NewRole(roleApplication)

	userHandler := handler.NewUser(userApplication)
	menuHandler := handler.NewMenu(menuApplication)
	loginHandler := handler.NewLogin(loginApplication)

	demosHandler := handler.NewDemos(demosApplication)

	weatherHandler := handler.NewWeather(weatherApplication)

	routerRouter := router.NewRouter(
		infrarepos.AuthRepository,
		syncedEnforcer,
		loginHandler,
		menuHandler,
		roleHandler,
		userHandler,
		demosHandler,
		weatherHandler,
	)

	//=========================================================================
	// Init Gin Engine
	//=========================================================================

	//gin.SetMode("debug")
	//gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()

	// CORS
	if configs.C.CORS.Enable {
		ginEngine.Use(middleware.CORSMiddleware())
	}

	//--------------------------------------------------------------
	ginEngine.GET("/swagger.yaml", func(c *gin.Context) {
		c.File("../../internal/presentation/swagger/swagger.yaml")
	})
	ginEngine.Static("/swaggerui/", "../../swaggerui")
	//--------------------------------------------------------------

	err = routerRouter.Register(ginEngine)
	if err != nil {
		log.Panic("ERROR: routerRouter.Register: ", err)
	}

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      ginEngine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting HTTP server on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Error starting server:", err)
		}

	}()

	wg.Wait()
	log.Println("Server stopped")

}
