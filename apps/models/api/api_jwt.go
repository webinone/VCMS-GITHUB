package api

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type JWTClaims struct {
	jwt.StandardClaims
	TenantId	string			`json:"tenant_id"`
	UserId          string		`json:"user_id"`
	UserName	string			`json:"user_name"`
	UserRole	string			`json:"user_role"`
}

func (c JWTClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	return nil
}

func GetJWTClaims (c echo.Context) *JWTClaims {
	claims 		:= c.Get("jwt").(*jwt.Token).Claims.(*JWTClaims)
	return claims
}


type FrontJWTClaims struct {
	jwt.StandardClaims
	MemberId      	string		`json:"member_id"`
	MemberName		string		`json:"member_name"`
	SessionId		string		`json:"session_id"`
}

func (c FrontJWTClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	return nil
}

func GetFrontJWTClaims (c echo.Context) *FrontJWTClaims {
	claims 		:= c.Get("jwt-front").(*jwt.Token).Claims.(*FrontJWTClaims)
	return claims
}