package auth

import (
	"context"
	"log"

	"time"

	"github.com/golang-jwt/jwt"

	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/auth"
	"github.com/andriykutsevol/DDDCasbinExample/pkg/util/hash"
)

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyFunc       jwt.Keyfunc
	expired       int
	tokenType     string
	rootUser      auth.RootUser
}

type Option func(*options)

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       172000,
	signingMethod: jwt.SigningMethodHS512,
	signingKey:    []byte(defaultKey),
	keyFunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(defaultKey), nil
	},
	rootUser: auth.RootUser{UserName: "root", Password: "rootpwd"},
}

type Storage interface {
	Set(ctx context.Context, tokenString string, expiration time.Duration) error
	Check(ctx context.Context, tokenString string) (bool, error)
}

type AuthRepository struct {
	//storage AuthStorageBridge
	storage Storage
	opts    *options
}

func NewRepository(storage Storage) *AuthRepository {
	o := defaultOptions
	return &AuthRepository{
		storage: storage,
		opts:    &o,
	}
}

func (r *AuthRepository) FindRootUser(ctx context.Context, userName string) *auth.RootUser {
	if userName == r.opts.rootUser.UserName {
		return &auth.RootUser{
			UserName: userName,
			Password: r.opts.rootUser.Password,
		}
	}
	return nil
}

func (r *AuthRepository) GenerateToken(ctx context.Context, userID string) (*auth.Auth, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(r.opts.expired) * time.Second).Unix()

	token := jwt.NewWithClaims(r.opts.signingMethod, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		Subject:   userID,
	})

	tokenString, err := token.SignedString(r.opts.signingKey)
	if err != nil {
		return nil, err
	}

	tokenInfo := &auth.Auth{
		ExpiresAt:   expiresAt,
		TokenType:   r.opts.tokenType,
		AccessToken: tokenString,
	}
	return tokenInfo, nil
}

func (r *AuthRepository) parseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, r.opts.keyFunc)
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return token.Claims.(*jwt.StandardClaims), nil
}

func (r *AuthRepository) DestroyToken(ctx context.Context, tokenString string) error {

	//claims, err := r.parseToken(tokenString)
	//expired := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
	//store.Set(ctx, tokenString, expired)

	err := r.storage.Set(ctx, "tokenString", time.Duration(10)*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) ParseUserID(ctx context.Context, tokenString string) (string, error) {
	log.Println("redis auth repo: ParseUserID")

	if tokenString == "" {
		return "", auth.ErrInvalidToken
	}

	claims, err := r.parseToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

func (r *AuthRepository) Release() error {
	return nil
}

const defaultKey = "qq-dashboard"



func SetRootUser(id, password string) Option {
	return func(o *options) {
		o.rootUser = auth.RootUser{
			UserName: id,
			Password: hash.MD5String(password),
		}
	}
}

func SetSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

func SetSigningKey(key interface{}) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

func SetKeyFunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyFunc = keyFunc
	}
}

func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}
