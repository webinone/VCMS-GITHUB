package api

import (
	appConfig "VCMS/apps/config"
	"VCMS/apps/handler"
	"VCMS/apps/libs"
	apiModel "VCMS/apps/models/api"
	rdbModel "VCMS/apps/models/rdb"
	"fmt"
	"net/http"
	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type AuthAPI struct {
	request AuthRequest
}

type AuthRequest struct {
	TenantId string `validate:"required" json:"tenant_id"`
	UserId   string `validate:"required" json:"user_id"`
	Password string `validate:"required" json:"password"`
}

// 사용자 로그인 JWT 토큰 발급
func (api AuthAPI) PostLogin(c echo.Context) error {

	payload := &api.request
	c.Bind(payload)

	fmt.Println(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 암호전 패스워드
	logrus.Debug("general Id : ", payload.UserId)
	logrus.Debug("general password : ", payload.Password)

	enc_pass := libs.Sha256Encoding(payload.Password)

	logrus.Debug("enc_pass : ", enc_pass)

	user := &rdbModel.User{
		TenantID: payload.TenantId,
		UserId:   payload.UserId,
		Password: enc_pass,
	}

	var count = 0
	tx := c.Get("Tx").(*gorm.DB)
	tx.Where("tenant_id = ? and user_id = ? and password = ? and use_yn = ?",
		user.TenantID, user.UserId, user.Password, "Y").Find(user).Count(&count)

	logrus.Debug(">>>> count : ", count)

	if count == 0 {
		// 존재하지 않는다면...
		return echo.NewHTTPError(http.StatusUnauthorized, "id or password not matched")
	}

	// JWT Key Get
	jwtKey := []byte(appConfig.Config.AUTH.JwtKey)
	token := jwt.New(jwt.SigningMethodHS256)

	claims := &apiModel.JWTClaims{
		TenantId: user.TenantID,
		UserId:   user.UserId,
		UserName: user.Name,
		UserRole: user.Role,
	}

	token.Claims = claims
	tokenString, _ := token.SignedString(jwtKey)

	fmt.Println("tokenString : ", tokenString)

	return handler.APIResultHandler(c, true, http.StatusOK, map[string]interface{}{"token_key": tokenString})
}
