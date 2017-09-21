package front

import (
	"github.com/labstack/echo"
	//"fmt"
	"net/http"
	//"io/ioutil"
	//"bytes"
	"VCMS/apps/handler"
	"github.com/Sirupsen/logrus"
	"VCMS/apps/db"
)

type NewsLetterAPI struct {
	request NewsLetterRequest
}

type NewsLetterRequest struct {
	Name	string 		`validate:"required" json:"name"`
	Email	string 		`validate:"required" json:"email"`
}

type NewsLetterResult struct {
	Name	string 		`json:"name"`
	Email	string 		`json:"email"`
}

// 뉴스 레터 등록
func (api NewsLetterAPI) PostNewsLetter(c echo.Context) error {

	payload := &api.request
	c.Bind(payload)

	logrus.Debug(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx := db.GetHelloTDB()
	tx.Begin()
	defer tx.Close()

	newsletter_results := []NewsLetterResult{}

	tx.Raw(`
		SELECT name, email FROM apply_newsletter
		WHERE email = ?
		`, payload.Email).Scan(&newsletter_results)

	if len(newsletter_results) > 0 {
		// 이미 존재한다.
		tx.Rollback()
		return echo.NewHTTPError(http.StatusConflict, "EXISTS_EMAIL")
	}


	tx.Exec(`
		INSERT INTO apply_newsletter
		(
			name,
			email,
			receive_flag,
			reg_date,
			reg_ip,
			apply_site
		)
		VALUES
		(
			?,
			?,
			?,
			NOW(),
			?,
			?
		)`, payload.Name, payload.Email, "Y", c.Request().RemoteAddr, "webinar.hellot.net")


	tx.Commit()

	return handler.APIResultHandler(c, true, http.StatusOK, "SUCCESS")
}
