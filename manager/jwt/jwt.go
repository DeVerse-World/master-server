package jwt

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	// "github.com/DCloudGaming/cloud-morph-host/pkg/env"
	// "github.com/DCloudGaming/cloud-morph-host/pkg/errors"
	"github.com/hyperjiang/gin-skeleton/model"
)

// jwt-cookie building and parsing
const cookieName = "deverse-jwt"

// tokens auto-refresh at the end of their lifetime,
// so long as the user hasn't been disabled in the interim
const tokenLifetime = time.Hour * 6

var hmacSecret []byte

func init() {
	hmacSecret = []byte(os.Getenv("TOKEN_SECRET"))
	if hmacSecret == nil {
		panic("No TOKEN_SECRET environment variable was found")
	}
}

type claims struct {
	User *model.Wallet
	jwt.StandardClaims
}

//
// // RequireAuth middleware makes sure the user exists based on their JWT
func RequireAuth(w http.ResponseWriter, r *http.Request) (*model.Wallet, bool) {
	u, err := HandleUserCookie(w, r)
	if err != nil || u.Address == "" {
		return u, false
	}

	return u, true
}

// WriteUserCookie encodes a user's JWT and sets it as an httpOnly & Secure cookie
func WriteUserCookie(w http.ResponseWriter, u *model.Wallet) {
	fmt.Println("Set cookie")
	fmt.Println(u)
	fmt.Println(EncodeUser(u))
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    EncodeUser(u),
		Path:     "/",
		SameSite: 4, //SameSiteNoneMode,
		HttpOnly: false,
		Secure:   true,
	})
}

// //// HandleUserCookie attempts to refresh an expired token if the user is still valid
func HandleUserCookie(w http.ResponseWriter, r *http.Request) (*model.Wallet, error) {
	u, err := userFromCookie(r)

	// attempt refresh of expired token:
	// if err == model.ErrExpiredToken {
	// 	user, fetchError := e.UserRepo().GetUser(u.WalletAddress)
	// 	if fetchError != nil {
	// 		return nil, err
	// 	}
	// 	if user.Status > 0 {
	// 		WriteUserCookie(w, user)
	// 		return user, nil
	// 	}
	// }

	return u, err
}

// userFromCookie builds a user object from a JWT, if it's valid
func userFromCookie(r *http.Request) (*model.Wallet, error) {
	cookie, _ := r.Cookie(cookieName)
	var tokenString string
	if cookie != nil {
		tokenString = cookie.Value
	} else {
		return nil, model.ErrInvalidToken
	}

	if tokenString == "" {
		return &model.Wallet{}, nil
	}

	return DecodeUser(tokenString)
}

// EncodeUser convert a user struct into a jwt
func EncodeUser(u *model.Wallet) (tokenString string) {
	claims := claims{
		u,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Add(-time.Second).Unix(),
			ExpiresAt: time.Now().Add(tokenLifetime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// unhandled err here
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		log.Println("Error signing token", err)
	}
	return
}

// DecodeUser converts a jwt into a user struct (or returns a zero-value user)
func DecodeUser(tokenString string) (*model.Wallet, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		// check for expired token
		if verr, ok := err.(*jwt.ValidationError); ok {
			if verr.Errors&jwt.ValidationErrorExpired != 0 {
				return getUserFromToken(token), model.ErrExpiredToken
			}
		}
	}

	if err != nil || !token.Valid {
		return nil, model.ErrInvalidToken
	}

	return getUserFromToken(token), nil
}

func getUserFromToken(token *jwt.Token) *model.Wallet {
	if claims, ok := token.Claims.(*claims); ok {
		return claims.User
	}

	return &model.Wallet{}
}
