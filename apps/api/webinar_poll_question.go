package api

import (
	//"VCMS/apps/handler"
	//"net/http"
	//"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	//appLibs "VCMS/apps/libs"
	rdbModel "VCMS/apps/models/rdb"
	apiModel "VCMS/apps/models/api"
	"net/http"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"VCMS/apps/handler"
	"github.com/satori/go.uuid"
	"strconv"
)

type WebinarPollQuestionAPI struct {
	requestPost WebinarPollQuestionPostRequest
	//requestPut  ContentRequestPut
}

type WebinarPollQuestionPostRequest struct {
	WebinarSiteId		string 							`validate:"required" json:"webinar_site_id"`
	WebinarPollId		string							`validate:"required" json:"webinar_poll_id"`
	WebinarPollQuestionMasters	[]WebinarPollQuestionMasterRequest		`validate:"required" json:"poll_questions"`
}

type WebinarPollQuestionMasterRequest struct {

	Title					string					`validate:"required" json:"title"`
	QuestionType				string					`validate:"required" json:"question_type"`
	QuestionCount				string					`validate:"required" json:"question_count"`
	WebinarPollQuestionDetails 	[]WebinarPollQuestionDetailRequest		`validate:"required" json:"webinar_question_detail"`

}

type WebinarPollQuestionDetailRequest struct {
	Order			int64				`validate:"required" json:"order"`
	Title			string				`validate:"required" json:"title"`
}


// 웨비나 사이트 등록
func (api WebinarPollQuestionAPI) PostWebinarPollQuestion(c echo.Context) error {

	logrus.Debug("########### PostWebinarPollQuestion !!!")

	// TODO : 참여자가 있는 경우 등록할 수 없다.

	payload := &api.requestPost
	c.Bind(payload)

	logrus.Debug(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims 	 := apiModel.GetJWTClaims(c)

	logrus.Debug("#### TenantId : ", claims.TenantId)

	tx := c.Get("Tx").(*gorm.DB)

	// TODO : 설문기간이 겹치는 설문은 등록할 수 없다.

	webinar_poll_question_master := &rdbModel.WebinarPollQuestionMaster{}
	webinar_poll_question_detail := &rdbModel.WebinarPollQuestionDetail{}

	// 1. 질문을 모두 삭제 한다.
	if tx.Delete(webinar_poll_question_detail, "webinar_site_id = ? and webinar_poll_id = ?",
		payload.WebinarSiteId, payload.WebinarPollId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	if tx.Delete(webinar_poll_question_master, "webinar_site_id = ? and webinar_poll_id = ?",
		payload.WebinarSiteId, payload.WebinarPollId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	// 2. 질문을 모두 다시 insert 한다.
	if len(payload.WebinarPollQuestionMasters) > 0 {

		for _, v := range payload.WebinarPollQuestionMasters {

			webinar_poll_question_master_id := uuid.NewV4().String()
			webinar_poll_question_details := []rdbModel.WebinarPollQuestionDetail{}

			for i, d := range v.WebinarPollQuestionDetails {

				s := append(webinar_poll_question_details, rdbModel.WebinarPollQuestionDetail{
					TenantId:claims.TenantId,
					WebinarSiteId:payload.WebinarSiteId,
					WebinarPollId:payload.WebinarPollId,
					WebinarPollQuestionMasterId:webinar_poll_question_master_id,
					WebinarPollQuestionDetailId:uuid.NewV4().String(),
					Title: d.Title,
					Order: i+1,
				})
				webinar_poll_question_details = s

			}

			webinar_poll_question_master := &rdbModel.WebinarPollQuestionMaster{
				TenantId	: claims.TenantId,
				WebinarSiteId	: payload.WebinarSiteId,
				WebinarPollId: payload.WebinarPollId,
				WebinarPollQuestionMasterId: webinar_poll_question_master_id,
				Title: v.Title,
				QuestionType:v.QuestionType,
				QuestionCount:v.QuestionCount,
				WebinarPollQuestionDetails:webinar_poll_question_details,
			}

			tx.Create(webinar_poll_question_master)
		}
	}

	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")

}

//설문 리스트 조회
func (api WebinarPollQuestionAPI) GetWebinarPollQuestions (c echo.Context) error  {

	claims 	 := apiModel.GetJWTClaims(c)

	tenant_id		:= claims.TenantId

	webinar_site_id		:= c.QueryParam("webinar_site_id");
	webinar_poll_id 	:= c.QueryParam("webinar_poll_id")		// 종료 여부


	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	logrus.Debug("##### tenant_id : ", tenant_id)
	logrus.Debug("##### webinar_poll_id : ", webinar_poll_id)

	webinar_poll_question_masters  := []rdbModel.WebinarPollQuestionMaster{}

	tx := c.Get("Tx").(*gorm.DB)

	tx  = tx.Where("tenant_id = ? and webinar_site_id = ? and webinar_poll_id = ?", tenant_id, webinar_site_id, webinar_poll_id)

	var count = 0
	tx.Find(&webinar_poll_question_masters).Count(&count)

	tx.Preload("WebinarPollQuestionDetails").Order("idx desc").Find(&webinar_poll_question_masters)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_poll_question_masters,
		})

}

// Webinar 설문 문항 한건 조회
func (api WebinarPollQuestionAPI) GetWebinarPollQuestion (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	webinar_poll_question_master  := &rdbModel.WebinarPollQuestionMaster{}

	//var count = 0
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Preload("WebinarSite").
		Preload("WebinarPollQuestionDetails").
		Where("idx = ? ", idx ).Find(webinar_poll_question_master).RecordNotFound() {

		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR POLL QUESTION NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, webinar_poll_question_master)
}


//// 설문 삭제
//func (api WebinarPollAPI) DeleteWebinarPoll (c echo.Context) error  {
//
//	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)
//
//	// TODO : 참여자가 있는 설문은 삭제 할 수 없다.
//
//	webinar_poll  			:= &rdbModel.WebinarPoll{}
//	webinar_poll_member_result  	:= &rdbModel.WebinarPollMemberResult{}
//	webinar_poll_member 		:= &rdbModel.WebinarPollMember{}
//	webinar_poll_question_master	:= &rdbModel.WebinarPollQuestionMaster{}
//	webinar_poll_question_detail	:= &rdbModel.WebinarPollQuestionDetail{}
//
//	tx := c.Get("Tx").(*gorm.DB)
//
//	if tx.Where("idx = ? ", idx ).Find(webinar_poll).RecordNotFound() {
//		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR POLL NOT FOUND")
//	}
//
//	webinar_poll_id := webinar_poll.WebinarPollId
//
//	if tx.Delete(webinar_poll, "webinar_poll_id = ?", webinar_poll_id).RowsAffected == 0 {
//		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
//	}
//
//	if tx.Delete(webinar_poll_member_result, "webinar_poll_id = ?", webinar_poll_id).RowsAffected == 0 {
//		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
//	}
//
//	if tx.Delete(webinar_poll_member, "webinar_poll_id = ?", webinar_poll_id).RowsAffected == 0 {
//		//return echo.NewHTTPError(http.StatusNotFound, "CONTENT NOT FOUND")
//	}
//
//	if tx.Delete(webinar_poll_question_detail, "webinar_poll_id = ?", webinar_poll_id).RowsAffected == 0 {
//
//	}
//
//	if tx.Delete(webinar_poll_question_master, "webinar_poll_id = ?", webinar_poll_id).RowsAffected == 0 {
//
//	}
//
//	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
//}
////
//
//// 설문 수정
//func (api WebinarPollAPI) PutWebinarPoll (c echo.Context) error  {
//
//	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)
//
//	// TODO : 참여자가 있거나 지난 일정은 수정 할 수 없다.
//
//	payload := &api.requestPost
//	c.Bind(payload)
//
//	if err   := c.Validate(payload); err != nil {
//		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
//	}
//
//	tx := c.Get("Tx").(*gorm.DB)
//
//	webinar_poll  	:= &rdbModel.WebinarPoll{}
//
//	if tx.Where("idx = ? ", idx ).Find(webinar_poll).RecordNotFound() {
//		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_POLL NOT FOUND")
//	}
//
//	webinar_poll.Title = payload.Title
//	webinar_poll.Desc = payload.Desc
//	webinar_poll.StartDate = payload.StartDate
//	webinar_poll.EndDate = payload.EndDate
//
//	tx.Save(webinar_poll)
//
//	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")
//
//}

