package jwt

import (
	model "file-sharing/model/error"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	UserId   uint   `json:"UserId"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("1031092qijdljdqkwdhiu2hdq98d9w87d9w8s98798f79f79d8f7g9sd8fg7")
var JwtShared = []byte("alsgkdalshlshdnvalssadklgh235klh45hkkbe2ebn5nlsjkflsdnlncln2")

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
	log.Info("part 1")
	tokenParse, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return jwtSecret, nil
	})
	log.Info("part 2")

	if err != nil || !tokenParse.Valid {
		return "", &model.InvalidJWTToken
	}

	claims, ok := tokenParse.Claims.(jwt.MapClaims)
	if !ok {
		return "", &model.InvalidJWTToken
	}
	log.Info("part3")

	username, ok := claims["username"].(string)
	//if username == "" {
	//	username = "anonymous" // Default value
	//}
	if !ok {
		log.Info("breakpoint")
		return "", &model.InvalidJWTToken
	}
	log.Info("part 4")

	return username, nil
}

func Create(username string, userId uint) (string, error) {

	expireTime := time.Now().Add(6 * time.Hour) // 6 saat

	claims := &Claims{
		Username: username,
		UserId:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)

}
func SharedToken(object string, expiresIn time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"object": object,
		"exp":    time.Now().Add(expiresIn).Unix(),
	}
	log.Info("test")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Info("token")
	return token.SignedString(JwtShared)
}

//func ValidateToken(token string) (*Claims, error) {
//	var claim *Claims
//	tkn, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
//		return jwtSecret, nil
//	})
//	if err != nil {
//		return nil, err
//	}
//	if claims, ok := tkn.Claims.(*Claims); ok && tkn.Valid {
//		return claims, nil
//	}
//	return nil, &model.InvalidJWTToken
//}
