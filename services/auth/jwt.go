package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Ricardoarsv/E-commerce_REST-API/config"
	"github.com/golang-jwt/jwt"
)

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
