package configs

import (
	"os"
	"strings"
	"sync"

	"github.com/koding/multiconfig"

	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/json"
)

var (
	C    = new(Config)
	once sync.Once
)

func MustLoad(fpaths ...string) {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
		}
		for _, fpath := range fpaths {
			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "json") {
				loaders = append(loaders, &multiconfig.JSONLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
			}
		}
		// Can be overwritten config with environment variable
		loaders = append(loaders, &multiconfig.EnvironmentLoader{})
		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}
		m.MustLoad(C)
	})
}

func PrintWithJSON() {
	if C.PrintConfig {
		b, err := json.MarshalIndent(C, "", " ")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

type Config struct {
	RunMode     string
	WWW         string
	Swagger     bool
	PrintConfig bool
	HTTP        HTTP
	Casbin      Casbin
	Root        Root
	JWTAuth     JWTAuth
	CORS        CORS
	Redis       Redis
}

type HTTP struct {
	Host             string
	Port             int
	CertFile         string
	KeyFile          string
	ShutdownTimeout  int
	MaxContentLength int64
	MaxLoggerLength  int `default:"4096"`
}

type Casbin struct {
	Enable           bool
	Debug            bool
	Model            string
	AutoLoad         bool
	AutoLoadInternal int
}

type Root struct {
	UserName string
	Password string
	RealName string
}

type JWTAuth struct {
	Enable        bool
	SigningMethod string
	SigningKey    string
	Expired       int
	Store         string
	FilePath      string
	RedisDB       int
	RedisPrefix   string
}

type CORS struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

type Redis struct {
	Addr     string
	Password string
}
