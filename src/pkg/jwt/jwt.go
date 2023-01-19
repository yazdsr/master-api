package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type Jwt interface {
	GenerateToken(claims map[string]interface{}) (string, error)
	ParseToken(token string) (map[string]interface{}, error)
}

type jwtpkg struct {
	secret string
}

func New(secret string) Jwt {
	return &jwtpkg{
		secret: secret,
	}
}

func (j *jwtpkg) GenerateToken(claims map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": claims["exp"],
		"iat": claims["iat"],
		"sub": claims["sub"],
	})
	return token.SignedString([]byte(j.secret))
}

func (j *jwtpkg) ParseToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("Invalid Token")
}
