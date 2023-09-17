package auth

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

func (j *JwtWrapper) GenerateToken(email string) (signedToken string, err error) {
	fmt.Println("time to expire", time.Now().Local().Add(time.Minute*time.Duration(j.ExpirationMinutes)).Unix())
	claims := &JwtClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(j.ExpirationMinutes)).Unix(),
			Issuer:    j.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}
	return
}

func (j *JwtWrapper) RefreshToken(email string) (signedToken string, err error) {
	claims := &JwtClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}
	return
}

func (j *JwtWrapper) ValidateToken(signedToken string) (err error, claims *JwtClaim) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	fmt.Println("Error", err)
	if err != nil {
		return err, nil
	}
	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return err, nil
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return err, nil
	}
	return
}

func (j *JwtWrapper) GetDetailsInToken(signedToken string) (claim *JwtClaim) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return
	}

	return claims
}
