package util

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

const JWT_LENGTH_KEY = 32

var JWT_SIGNATURE_KEY []byte
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

type JWTStruct struct {
	Username string `json:"username"`
}

func init() {
	JWT_SIGNATURE_KEY = []byte(GenRandomString(JWT_LENGTH_KEY))
	print("JWT Signature Key: ", string(JWT_SIGNATURE_KEY))
}

func CreateJWTSign(data *JWTStruct) (string, error) {
	return jwt.NewWithClaims(JWT_SIGNING_METHOD, jwt.MapClaims{
		"username": data.Username,
	}).SignedString(JWT_SIGNATURE_KEY)
}

func ValidateJWTSign(token string) (*JWTStruct, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return JWT_SIGNATURE_KEY, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		return &JWTStruct{
			Username: claims["username"].(string),
		}, nil
	} else {
		return nil, err
	}
}
