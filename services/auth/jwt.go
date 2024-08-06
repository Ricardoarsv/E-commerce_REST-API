package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ricardoarsv/E-commerce_REST-API/config"
	"github.com/Ricardoarsv/E-commerce_REST-API/types"
	"github.com/Ricardoarsv/E-commerce_REST-API/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/net/context"
)

type contextKey string

const UserKey contextKey = "UserID"

func CreateJWT(secret []byte, userID, userROLE int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	// todo implement JWT creation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"userROLE":  strconv.Itoa(userROLE),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		return tokenAuth
	}

	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permision denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}
	return userID
}

func ValidateJwtTokenForCreateProducts(tokenString string) error {
	expectedRole := config.Envs.USERADMIN // Asegúrate de que este valor sea 2

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	userRoleStr, ok := claims["userROLE"].(string)
	if !ok {
		return fmt.Errorf("userROLE not found in token")
	}

	userRole, err := strconv.Atoi(userRoleStr)
	if err != nil {
		return fmt.Errorf("userROLE is not a valid integer")
	}

	if int64(userRole) != expectedRole {
		fmt.Printf("userRole: %d, expectedRole: %d\n", userRole, expectedRole)
		return fmt.Errorf("user does not have admin privileges")
	}

	return nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequest(r)

		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			permissionDenied(w)
			return
		}

		// Add the user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userID)
		r = r.WithContext(ctx)

		// Call the function if the token is valid
		handlerFunc(w, r)
	}
}
