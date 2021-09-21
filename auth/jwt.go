package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/revzim/nano/internal/log"
)

type (
	JWT struct {
		algo          string
		Parse         func(tokenString string) jwt.MapClaims
		GenerateToken func(id string, name string, duration int64) (string, error)
	}
)

var (
	jwtSigningKey []byte
	jwtAlgo       string
)

// func init() {}

func NewJWT(signKey, algo string, genTokenFunc func(id, name string, duration int64) (string, error)) *JWT {
	initJWT(signKey, algo)
	if genTokenFunc == nil {
		genTokenFunc = generateJWTToken
	}
	return &JWT{
		algo:          jwtAlgo,
		Parse:         parseJWTToken,
		GenerateToken: genTokenFunc,
	}
}

func initJWT(signKey, algo string) {
	jwtSigningKey = []byte(signKey)
	jwtAlgo = algo
}

func jwtKeyFunc(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != jwtAlgo {
		return nil, fmt.Errorf("bad signing method: %v\n", t.Header["alg"])
	}
	return jwtSigningKey, nil
}

func parseJWTToken(tokenString string) jwt.MapClaims {
	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		log.Println("auth err:", err)
		return jwt.MapClaims{
			"error": err.Error(),
		}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	}

	return jwt.MapClaims{
		"error": err.Error(),
	}
}

func generateJWTToken(id, name string, duration int64) (string, error) {
	nowTime := time.Now().Unix()
	claims := &jwt.MapClaims{
		"id":   id,
		"name": name,
		"iat":  nowTime,
		"nbf":  nowTime - 10,
		"exp":  nowTime + duration,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSigningKey)
	if err != nil {
		log.Println(err)
	}
	return tokenString, err
}

func generateJWTTokenWithClaims(id, name string, duration int64, claims *jwt.MapClaims) (string, error) {
	if claims == nil {
		nowTime := time.Now().Unix()
		claims = &jwt.MapClaims{
			"id":   id,
			"name": name,
			"iat":  nowTime,
			"nbf":  nowTime - 10,
			"exp":  nowTime + duration,
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSigningKey)
	if err != nil {
		log.Println(err)
	}
	return tokenString, err
}
