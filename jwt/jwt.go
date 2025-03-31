package jwt

import (
	model "file-sharing/model/error"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserId uint `json:"UserId"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("1031092qijdljdqkwdhiu2hdq98d9w87d9w8s98798f79f79d8f7g9sd8fg7")

func JwtParse(token string) (string, error) {

	var tokenString string

	if token == "" {
		return "", &model.InvalidJWTToken
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		tokenString = token[7:]
	} else {
		return "", &model.InvalidJWTToken
	}

	tokenParse, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return jwtSecret, nil
	})

	if err != nil || !tokenParse.Valid {
		return "", &model.InvalidJWTToken
	}

	claims, ok := tokenParse.Claims.(jwt.MapClaims)
	if !ok {
		return "", &model.InvalidJWTToken
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", &model.InvalidJWTToken
	}

	return username, nil
}

func Create(userId uint) (string, error) {

	expireTime := time.Now().Add(6 * time.Hour) // 6 saat

	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)

}

//func Verify(token string) (string, error) {
//
//	parsedToken := token[len("Bearer "):]
//
//	claims, err := jwt.ParseWithClaims(parsedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
//		return jwtSecret, nil
//	})
//
//	if err != nil {
//		return "", err
//	}
//
//	return claims.Claims.(*Claims).Username, nil
//
//}
