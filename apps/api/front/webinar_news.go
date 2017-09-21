package front

import (
	"github.com/labstack/echo"
	appConfig "VCMS/apps/config"
	"net/http"
	"fmt"
	"io/ioutil"
	//"encoding/json"
	"encoding/json"
	"VCMS/apps/handler"
)

type WebinarNewsAPI struct {
	//requestPost WebinarBannerMasterRequest
	//requestPut  ContentRequestPut
}

type HelloTNewsResult struct {
	Success 		string 			`json:"success"`
	ResultCode  	string  		`json:"resultCode"`
	ResultData 		interface{} 	`json:"resultData"`
}

 //배너 리스트 조회
func (api WebinarNewsAPI) GetWebinarNews (c echo.Context) error  {

	//claims 	 := apiModel.GetJWTClaims(c)

	offset 		 := c.QueryParam("offset")
	limit		 := c.QueryParam("limit")

	webinar_site_id	 := c.QueryParam("webinar_site_id")

	url := appConfig.Config.HELLOT.ApiUrl + "/webinar_article.php?code="+webinar_site_id+"&offset="+offset+"&limit="+limit
	resp, err := http.Get(url)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "HELLOT_SITE_ERROR")
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(respBody))

	data := &HelloTNewsResult{}

	json.Unmarshal([]byte(respBody), data)


	return handler.APIResultHandler(c, true, http.StatusOK,
		data.ResultData)

}

