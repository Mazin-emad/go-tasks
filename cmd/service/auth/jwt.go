package auth

import (
	"fmt"
	"time"

	"github.com/Mazin-emad/todo-backend/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(secret []byte,userID int) (string, error) {

	expiresAt := time.Second * time.Duration(config.ConfigAmigoo.JWTExpiresIn)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp": time.Now().Add(expiresAt).Unix(),
	})

	stringToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return stringToken, nil
}


func VerifyToken(secret []byte, tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return 0, err
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}
	
	userId, ok := claims["userId"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID in token")
	}
	
	return int(userId), nil
}

func GetUserIDFromToken(secret []byte, tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return 0, err
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}
	
	userId, ok := claims["userId"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID in token")
	}
	
	return int(userId), nil
}
