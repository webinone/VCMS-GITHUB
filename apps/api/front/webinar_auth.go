package front

import (
	"github.com/labstack/echo"
	//"fmt"
	"net/http"
	//"io/ioutil"
	//"bytes"
	"VCMS/apps/handler"
	apiModel "VCMS/apps/models/api"
	appConfig "VCMS/apps/config"
	"github.com/Sirupsen/logrus"
	"fmt"
	"io/ioutil"
	"bytes"
)

type AuthAPI struct {
	request AuthRequest
}

type AuthRequest struct {
	Token	string 		`validate:"required" json:"token"`
}

// 세션 체크
func (api AuthAPI) PostSessionCheck(c echo.Context) error {

	claims := apiModel.GetFrontJWTClaims(c)
	payload := &api.request
	c.Bind(payload)

	logrus.Debug(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	logrus.Debug("MEMBER ID : ", claims.MemberId)
	logrus.Debug("MEMBER NAME : ", claims.MemberName)
	logrus.Debug("TOKEN : ", payload.Token)


	// Cookie는 존재하지만
	url := appConfig.Config.AUTH.SSOUrl + "/api/v1/sso/auth/active"

	reqBody := []byte(`
					{
						"token" : "`+payload.Token+`"
					}
				`)

	var jsonStr = []byte(reqBody)

	fmt.Println(string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Accept","application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusUnauthorized, "EXPIRED")

	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.StatusCode)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	if resp.StatusCode == http.StatusNonAuthoritativeInfo {
		return handler.APIResultHandler(c, false, http.StatusUnauthorized, "EXPIRED")

	} else {
		return handler.APIResultHandler(c, true, http.StatusOK, "SUCCESS")
	}
}


// 회원 가입 후 헬로우 T 사용자 가져오기
func (api AuthAPI) PostAddUser (c echo.Context) error {

	// 회원 DB에서 가져와서 넣는다.

	return handler.APIResultHandler(c, true, http.StatusOK, "SUCCESS")
}
