package cognito

import (
	"log"
	"context"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/auth"
	"github.com/golang-jwt/jwt"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/hash"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
    //"github.com/casbin/casbin/v2"

)


type options struct {
	rootUser      auth.RootUser
}
type Option func(*options)

var defaultOptions = options{
	rootUser: auth.RootUser{UserName: "root", Password: "rootpwd"},
}



type CongnitoRepository struct {
	opts    *options
}


func NewRepository() *CongnitoRepository {
	return &CongnitoRepository{
	}
}

func (r *CongnitoRepository) FindRootUser(ctx context.Context, userName string) *auth.RootUser {
	log.Println("cognito repo: FindRootUser()")

	return nil
}



func (r *CongnitoRepository) GenerateToken(ctx context.Context, userID string) (*auth.Auth, error) {
	log.Println("DEVMODE: cognito repo: GenerateToken()")


	userPoolID := "your_user_pool_id"
	clientID   := "your_client_id"
	region     := "your_region"


	// svc := cognitoidentityprovider.New(session.New(), &aws.Config{Region: aws.String("us-west-2")})

	// authInput := &cognitoidentityprovider.InitiateAuthInput{
	// 	AuthFlow: aws.String("USER_PASSWORD_AUTH"),
	// 	AuthParameters: map[string]*string{
	// 		"USERNAME": aws.String(username),
	// 		"PASSWORD": aws.String(password),
	// 	},
	// 	ClientId: aws.String("your_cognito_app_client_id"),
	// }


	// TODO: Cognito
	sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String(region),
    }))
    cognitoClient := cognitoidentityprovider.New(sess)
	_ = cognitoClient

	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(username),
			"PASSWORD": aws.String(password),
		},
		ClientId: aws.String("your_cognito_app_client_id"),
	}


	authOutput, err := svc.InitiateAuth(authInput)
	if err != nil {
		return "", err
	}



	// TODO: Cognito
	sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String(region),
    }))
    cognitoClient := cognitoidentityprovider.New(sess)
	_ = cognitoClient



	// // TODO: Cognito
	// sess := session.Must(session.NewSession(&aws.Config{
    //     Region: aws.String(region),
    // }))
    // cognitoClient := cognitoidentityprovider.New(sess)
	// _ = cognitoClient


    // // Initiate authentication with USER_PASSWORD_AUTH flow
    // input := &cognitoidentityprovider.InitiateAuthInput{
    //     AuthFlow:       cognitoidentityprovider.AuthFlowTypeUSER_PASSWORD_AUTH,
    //     AuthParameters: map[string]string{
    //         "USERNAME": username,
    //         "PASSWORD": password,
    //     },
    //     ClientId:       aws.String(clientID),
    //     UserPoolId:     aws.String(userPoolID),
    // }




	log.Println("DEVMODE: cognito repo: GenerateToken() DONE")
	return nil, nil
}



func (r *CongnitoRepository) parseToken(tokenString string) (*jwt.StandardClaims, error) {
	log.Println("cognito repo: parseToken()")
	return nil, nil
}


func (r *CongnitoRepository) DestroyToken(ctx context.Context, tokenString string) error {
	log.Println("cognito repo: DestroyToken()")

	return nil
}

func (r *CongnitoRepository) ParseUserID(ctx context.Context, tokenString string) (string, error) {
	log.Println("cognito repo: ParseUserID()")

	return "", nil
}

func (r *CongnitoRepository) Release() error {
	return nil
}


func SetRootUser(id, password string) Option {
	return func(o *options) {
		o.rootUser = auth.RootUser{
			UserName: id,
			Password: hash.MD5String(password),
		}
	}
}



func SetSigningMethod(method jwt.SigningMethod) Option {
	log.Println("cognito repo: SetSigningMethod()")
	return func(o *options) {}
}

func SetSigningKey(key interface{}) Option {
	log.Println("cognito repo: SetSigningKey()")
	return func(o *options) {}
}


func SetKeyFunc(keyFunc jwt.Keyfunc) Option {
	log.Println("cognito repo: SetKeyFunc()")
	return func(o *options) {}
}

func SetExpired(expired int) Option {
	log.Println("cognito repo: SetExpired()")
	return func(o *options) {}	
}


























