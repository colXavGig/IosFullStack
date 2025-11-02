package web

import (
	"context"
	"errors"
	"fmt"
	"ios_full_stack/dto"
	"ios_full_stack/models"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	userId uint
	jwt.RegisteredClaims
}

var (
	ctxUserKey string = "APP_USER"

	signingKeySecret string = os.Getenv("JWT_SECRET")
	cookieAuthKey    string = "auth-token"
)

var (
	ErrUnauthorizedAccess error = errors.New("Not authorized")
)

func GenerateToken(user dto.User) (*http.Cookie, error) {
	var (
		expiration time.Time  = time.Now().Add(time.Hour * 48)
		userClaims UserClaims = UserClaims{
			userId: user.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiration),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "ios_full_stack",
			},
		}
		token jwt.Token = *jwt.NewWithClaims(jwt.SigningMethodES256, &userClaims)
	)

	signedToken, err := token.SignedString(signingKeySecret)
	if err != nil {
		return nil, fmt.Errorf("could not sign token: %w", err)
	}

	return &http.Cookie{
		Name:    cookieAuthKey,
		Value:   signedToken,
		Expires: expiration,
	}, nil
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx    context.Context = r.Context()
			claims *UserClaims
			user   dto.User

			err    error
			status int = http.StatusAccepted
		)
		defer func() {
			if err != nil {
				_ = JsonResponseWriter(w, &HttpError{
					Code:    status,
					Message: err.Error(),
					Details: err.Error(),
				})
			}
		}()

		tokenCookie, err := r.Cookie(cookieAuthKey)
		if err != nil || tokenCookie.Expires.Before(time.Now()) {
			err = ErrUnauthorizedAccess
			status = http.StatusUnauthorized
			return
		}

		tokenStr := tokenCookie.Value

		token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (any, error) {
			return signingKeySecret, nil
		})
		if err != nil || !token.Valid {
			err = errors.New("invalid token")
			status = http.StatusInternalServerError
			return
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			err = errors.New("internal error while claiming user from token")
			status = http.StatusInternalServerError
			return
		}

		appUser, err := models.TryFindUserById(ctx, claims.userId)
		if errors.Is(err, models.ErrUserNotFound) {
			status = http.StatusUnauthorized
			return
		} else if err != nil {
			log.Panicln(err)
			return
		}

		ctx = context.WithValue(ctx, ctxUserKey, appUser)

		next.ServeHTTP(w, r.WithContext(ctx))

		refreshedCookie, err := GenerateToken(user)
		if err != nil {
			err = errors.New("error refreshing token")
			status = http.StatusInternalServerError
			return
		}

		http.SetCookie(w, refreshedCookie)
	})
}

func GetUserFromContext(ctx context.Context) *models.AppUser {
	user, ok := ctx.Value(ctxUserKey).(*models.AppUser)
	if !ok {
		log.Panicln(errors.New("could not find user in context"))
	}

	return user
}
