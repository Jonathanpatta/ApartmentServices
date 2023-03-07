package Middleware

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"net/http"
	"strings"
)

type MiddlwareService struct {
	auth *auth.Client
}

type FirebaseUser struct {
	Name    string
	Email   string
	Picture string
	UserId  string
}

func GetFirebaseUser(ctx context.Context) *FirebaseUser {
	user := FirebaseUser{}
	ctxUser := ctx.Value("user")
	if ctxUser != nil {
		user = ctxUser.(FirebaseUser)
	}
	return &user
}

func GetFirebaseUserFromToken(token *auth.Token) FirebaseUser {
	user := FirebaseUser{}
	name := token.Claims["name"]
	if name != nil {
		user.Name = name.(string)
	}
	userId := token.Claims["user_id"]
	if userId != nil {
		user.UserId = userId.(string)
	}
	picture := token.Claims["picture"]
	if picture != nil {
		user.Picture = picture.(string)
	}

	email := token.Firebase.Identities["email"]
	if email != nil {
		emailsInterface := email.([]interface{})
		if emailsInterface != nil {
			if len(emailsInterface) > 0 {
				myEmail := emailsInterface[0]
				if myEmail != nil {
					user.Email = myEmail.(string)
				}
			}
		}
	}
	return user
}

func NewMiddlwareService(auth *auth.Client) (*MiddlwareService, error) {
	return &MiddlwareService{
		auth: auth,
	}, nil
}

const TokenName = "Authorization"

func AuthError(w http.ResponseWriter, r *http.Request, err error) {
	errStr := fmt.Sprintf("Auth Error:%v", err)
	http.Error(w, errStr, http.StatusUnauthorized)
}

func (s *MiddlwareService) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get(TokenName)
		val = strings.ReplaceAll(val, "Bearer ", "")
		if val == "" {
			AuthError(w, r, errors.New("invalid token"))
			return
		}
		token, err := s.auth.VerifyIDToken(context.Background(), val)
		if err != nil {
			AuthError(w, r, err)
			return
		}

		user := GetFirebaseUserFromToken(token)

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (s *MiddlwareService) ValidateTokenHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get(TokenName)
		val = strings.ReplaceAll(val, "Bearer ", "")
		if val == "" {
			AuthError(w, r, errors.New("invalid token"))
			return
		}
		token, err := s.auth.VerifyIDToken(context.Background(), val)
		if err != nil {
			AuthError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), "Claims", token.Claims)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
