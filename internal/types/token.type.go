package types

import "github.com/golang-jwt/jwt"

type JwtToken struct {
	Token     string
	ExpiresIn int
}

type MyCustomClaims struct {
	Payload interface{} `json:"payload"`
	jwt.StandardClaims
}
