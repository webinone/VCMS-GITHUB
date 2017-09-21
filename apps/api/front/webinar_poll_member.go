package front


import (
	"github.com/labstack/echo"
	"github.com/Sirupsen/logrus"
	"net/http"
	rdbModel "VCMS/apps/models/rdb"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"VCMS/apps/handler"
	apiModel "VCMS/apps/models/api"
	"fmt"
)

type WebinarPollMemberAPI struct {
	requestPost WebinarPollMemberPostRequest
	//requestPut  ContentRequestPut
}

type WebinarPollMemberPostRequest struct {
	WebinarSiteId			string 				`validate:"required" json:"webinar_site_id"`
	WebinarPollId			string				`validate:"required" json:"webinar_poll_id"`

	//FrontUserId			string				`validate:"required" json:"front_user_id"`
	WinYN				string				`validate:"required" json:"win_yn"`
	WebinarPollMemberResults 	[]WebinarPollMemberResultRequest `validate:"required" json:"poll_member_results"`

}

type WebinarPollMemberResultRequest struct {
	WebinarPollQuestionMasterId	string 			`validate:"required" json:"webinar_poll_question_master_id"`
	Answer			string				`validate:"required" json:"answer"`

}

// 웨비나 설문 참여 등록
func (api WebinarPollMemberAPI) PostWebinarPollMember(c echo.Context) error {

	logrus.Debug("########### PostWebinarPollMember !!!")



	claims 	 := apiModel.GetFrontJWTClaims(c)

	fmt.Println("############## Claims START ####################")
	fmt.Println("member_id : " + claims.MemberId)
	fmt.Println("member_name : " + claims.MemberName)
	fmt.Println("############## Claims END ####################")

	// TODO : 이미 참가 했으면 참가 못하도록 막아야 함.

	payload := &api.requestPost
	c.Bind(payload)

	logrus.Debug(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	webinar_site     := &rdbModel.WebinarSite{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("webinar_site_id = ? ", payload.WebinarSiteId ).Find(webinar_site).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_SITE NOT FOUND")
	}

	// webinar_site_id로 tenant_id 구하기
	tenant_id := webinar_site.TenantId

	webinar_poll_member_results := []rdbModel.WebinarPollMemberResult{}

	// 2. 배너를 모두 다시 insert 한다.
	if len(payload.WebinarPollMemberResults) > 0 {

		for _, v := range payload.WebinarPollMemberResults {

			s := append(webinar_poll_member_results, rdbModel.WebinarPollMemberResult{
				TenantId:tenant_id,
				WebinarSiteId:payload.WebinarSiteId,
				WebinarPollId:payload.WebinarPollId,
				WebinarPollQuestionMasterId:v.WebinarPollQuestionMasterId,
				FrontUserId:claims.MemberId,
				Answer: v.Answer,
			})
			webinar_poll_member_results = s
		}
	}

	webinar_poll_member := &rdbModel.WebinarPollMember{

		TenantId:tenant_id,
		WebinarSiteId: payload.WebinarSiteId,
		WebinarPollId: payload.WebinarPollId,
		WebinarPollMemberId: uuid.NewV4().String(),
		FrontUserId: claims.MemberId,
		FrontUserName: claims.MemberName,
		WinYN: "N",
		WebinarPollMemberResults : webinar_poll_member_results,
	}

	tx.Create(webinar_poll_member)



	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")

}