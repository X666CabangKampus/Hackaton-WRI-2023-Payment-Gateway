package util

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

var JWT_SIGNATURE_KEY = []byte("_22Y~!dOLlHEv8I@\",Wr(y(iar?RaBdt$j<@2u;^v2HymCC!H-$s1wt8S%^8CG6")
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

type JWTStruct struct {
	Username string `json:"username"`
}

func CreateJWTSign(data *JWTStruct) (string, error) {
	return jwt.NewWithClaims(JWT_SIGNING_METHOD, jwt.MapClaims{
		"username": data.Username,
	}).SignedString(JWT_SIGNATURE_KEY)
}

func ValidateJWTSign(token string) (*JWTStruct, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return JWT_SIGNATURE_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	return &JWTStruct{
		Username: claims["username"].(string),
	}, nil
}
