package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/revzim/nano/internal/log"
)

type (
	JWTFunc func(claimsMap map[string]interface{}, duration int64) (string, error) // func(claimsMap map[string]interface{}, duration int64) (string, error) // func(id string, name string, duration int64, claims map[string]interface{}) (string, error) // func(id string, name string, duration int64) (string, error)

	JWT struct {
		algo          string
		Parse         func(tokenString string) jwt.MapClaims
		GenerateToken JWTFunc // func(claimsMap map[string]interface{}, duration int64) (string, error)
	}
)

var (
	jwtSigningKey []byte
	jwtAlgo       string
)

// func init() {}

func NewJWT(signKey, algo string, genTokenFunc JWTFunc) *JWT {
	initJWT(signKey, algo)
	if genTokenFunc == nil {
		// genTokenFunc = generateJWTToken
		genTokenFunc = generateJWTTokenWithClaims
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
	return token.SignedString(jwtSigningKey)
	// tokenString, err := token.SignedString(jwtSigningKey)
	// if err != nil {
	// 	log.Println(err)
	// }
	// return tokenString, err
}

func generateJWTTokenWithClaims(claimsMap map[string]interface{}, duration int64) (string, error) {
	nowTime := time.Now().Unix()
	if claimsMap == nil {
		return "", errors.New("no claims!")
	}
	claimsMap["iat"] = nowTime
	claimsMap["nbf"] = nowTime - 10
	claimsMap["exp"] = nowTime + duration

	b, _ := json.Marshal(claimsMap)
	var claims *jwt.MapClaims
	err := json.Unmarshal(b, &claims)
	if err != nil {
		log.Println(err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSigningKey)
}
