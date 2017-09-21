package front

import (
	//"VCMS/apps/handler"
	//"net/http"
	//"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	//appLibs "VCMS/apps/libs"
	rdbModel "VCMS/apps/models/rdb"
	apiModel "VCMS/apps/models/api"
	//"github.com/jinzhu/gorm"
	//"github.com/satori/go.uuid"

	"net/http"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"VCMS/apps/handler"
	"github.com/satori/go.uuid"

	"strconv"
	"time"
	"fmt"
	"github.com/metakeule/fmtdate"
)

type WebinarPollAPI struct {
	requestPost WebinarPollPostRequest
	//requestPut  ContentRequestPut
}

type WebinarPollPostRequest struct {
	WebinarSiteId		string 				`validate:"required" json:"webinar_site_id"`
	Title			string				`validate:"required" json:"title"`
	Desc			string				`json:"desc"`
	StartDate		string				`validate:"required" json:"start_date"`
	EndDate			string				`validate:"required" json:"end_date"`
}

//type WebinarBannerSubRequest struct {
//
//	WebinarSiteId		string 		`validate:"required" json:"webinar_site_id"`
//	BannerType		string		`validate:"required" json:"banner_type"` // 1: 경품배너 2:업체배너 3: 진행페이지 배너
//	BannerTitle		string		`json:"banner_title"`
//	BannerDesc		string		`json:"banner_desc"`
//	LinkUrl			string		`validate:"required" json:"link_url"`
//	Order                   int		`json:"order"`
//	SavePath		string 		`validate:"required" json:"save_path"`
//	WebPath			string 		`validate:"required" json:"web_path"`
//
//}
//
//type ContentRequestPut struct {
//	CategoryId		string		`validate:"required" json:"category_id"`
//	ContentName 		string 		`validate:"required" json:"content_name"`
//	ThumbChange		string		`validate:"required" json:"thumb_change"`
//	ThumbType		string		`json:"thumb_type"`
//	ThumbTime		string		`json:"thumb_time"`
//
//}



// 웨비나 사이트 등록
func (api WebinarPollAPI) PostWebinarPoll(c echo.Context) error {

	logrus.Debug("########### PostWebinarPoll !!!")

	payload := &api.requestPost
	c.Bind(payload)

	logrus.Debug(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims 	 := apiModel.GetJWTClaims(c)

	logrus.Debug("#### TenantId : ", claims.TenantId)

	tx := c.Get("Tx").(*gorm.DB)

	webinar_poll_id := uuid.NewV4().String()

	// TODO : 설문기간이 겹치는 설문은 등록할 수 없다.
	// 시작일이 등록된 설문 중에 사이에 존재하면 등록할수 없다.

	webinar_polls := []rdbModel.WebinarPoll{}
	tx.Where(" start_date <= ? and end_date >= ? ", payload.StartDate, payload.StartDate ).Find(&webinar_polls)

	if (len(webinar_polls) > 0) {
		return echo.NewHTTPError(http.StatusBadRequest, "WEBINAR_POLL_EXIST")
	}

	// 설문등록
	webinar_poll := &rdbModel.WebinarPoll{
		TenantId: claims.TenantId,
		WebinarSiteId: payload.WebinarSiteId,
		WebinarPollId:webinar_poll_id,
		Title:payload.Title,
		Desc:payload.Desc,
		StartDate:payload.StartDate,
		EndDate:payload.EndDate,
		UpdatedId:claims.UserId,
	}

	tx.Create(webinar_poll)

	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")

}

//설문 리스트 조회
func (api WebinarPollAPI) GetWebinarPolls (c echo.Context) error  {

	claims 	 := apiModel.GetJWTClaims(c)

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	tenant_id		:= claims.TenantId

	webinar_site_id		:= c.QueryParam("webinar_site_id");
	status 			:= c.QueryParam("status")		// 종료 여부

	title			:= c.QueryParam("title")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	logrus.Debug("##### status : ", status)
	logrus.Debug("##### title : ", title)

	webinar_polls  := []rdbModel.WebinarPoll{}

	tx := c.Get("Tx").(*gorm.DB)

	tx  = tx.Where("tenant_id = ? and webinar_site_id = ? ", tenant_id, webinar_site_id)

	if status != "" {
		offSet, _ := time.ParseDuration("+09.00h")
		now := time.Now().UTC().Add(offSet)
		fmt.Println("Today : ", fmtdate.Format("YYYY-MM-DD hh:mm:ss", now))
		today := fmtdate.Format("YYYY-MM-DD", now)

		if status == "1" {
			// 미진행
			tx = tx.Where("start_date > ? ",  today)
		} else if (status == "2") {
			// 진행중
			tx = tx.Where("start_date <= ? and end_date >= ? ",  today, today)
		} else if (status == "3") {
			// 종료
			tx = tx.Where("end_date < ? ",  today)
		}
	}


	if title != "" {
		tx = tx.Where("title LIKE ? ", "%" + title + "%")
	}

	var count = 0
	tx.Find(&webinar_polls).Count(&count)


	tx.Preload("WebinarPollMembers").Preload("WebinarPollQuestionMasters").
		Order("idx desc").Offset(offset).Limit(limit).Find(&webinar_polls)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_polls,
		})

}

func (api WebinarPollAPI) GetWebinarPoll (c echo.Context) error  {

	claims 	 := apiModel.GetFrontJWTClaims(c)

	webinar_site_id		:= c.QueryParam("webinar_site_id");

	webinar_poll  := &rdbModel.WebinarPoll{}

	offSet, _ := time.ParseDuration("+09.00h")
	now := time.Now().UTC().Add(offSet)
	nowDateString := fmtdate.Format("YYYY-MM-DD", now)

	//var count = 0
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Preload("WebinarSite").
		Preload("WebinarPollQuestionMasters").
		Preload("WebinarPollQuestionMasters.WebinarPollQuestionDetails").
		Where("webinar_site_id = ? and start_date <= ? and end_date >= ? ", webinar_site_id, nowDateString, nowDateString ).
		Find(webinar_poll).RecordNotFound() {

		//return echo.NewHTTPError(http.StatusNotFound, "NOT_FOUND")
		//webinar_poll.WebinarPollId = "00";
	}

	webinar_poll_member := &rdbModel.WebinarPollMember{}

	join_yn := "N"

	if tx.Where("webinar_poll_id =? and front_user_id = ? ", webinar_poll.WebinarPollId, claims.MemberId ).Find(webinar_poll_member).RecordNotFound() {
		join_yn = "N"
	} else {
		join_yn = "Y"
	}

	return handler.APIResultHandler(c, true, http.StatusOK, map[string]interface{}{
		"poll_join_yn": join_yn,
		"webinar_poll": webinar_poll,
	})
}

// 설문 삭제
func (api WebinarPollAPI) DeleteWebinarPoll (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	// TODO : 참여자가 있는 설문은 삭제 할 수 없다.

	webinar_poll  			:= &rdbModel.WebinarPoll{}
	webinar_poll_member_result  	:= &rdbModel.WebinarPollMemberResult{}
	webinar_poll_member 		:= &rdbModel.WebinarPollMember{}
	webinar_poll_question_master	:= &rdbModel.WebinarPollQuestionMaster{}
	webinar_poll_question_detail	:= &rdbModel.WebinarPollQuestionDetail{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("idx = ? ", idx ).Find(webinar_poll).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR POLL NOT FOUND")
	}

	webinar_poll_id := webinar_poll.WebinarPollId

	if tx.Delete(webinar_poll, "webinar_poll_id = ?", webinar_poll_id).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	if tx.Delete(webinar_poll_member_result, "webinar_poll_id = ?", webinar_poll_id).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	if tx.Delete(webinar_poll_member, "webinar_poll_id = ?", webinar_poll_id).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "CONTENT NOT FOUND")
	}

	if tx.Delete(webinar_poll_question_detail, "webinar_poll_id = ?", webinar_poll_id).RowsAffected == 0 {

	}

	if tx.Delete(webinar_poll_question_master, "webinar_poll_id = ?", webinar_poll_id).RowsAffected == 0 {

	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}
//

// 설문 수정
func (api WebinarPollAPI) PutWebinarPoll (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	// TODO : 참여자가 있거나 지난 일정은 수정 할 수 없다.

	payload := &api.requestPost
	c.Bind(payload)

	if err   := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx := c.Get("Tx").(*gorm.DB)

	// TODO : 설문기간이 겹치는 설문은 등록할 수 없다.
	// 시작일이 등록된 설문 중에 사이에 존재하면 등록할수 없다.

	webinar_polls := []rdbModel.WebinarPoll{}
	tx.Where(" start_date <= ? and end_date >= ? ", payload.StartDate, payload.StartDate ).Find(webinar_polls)

	if (len(webinar_polls) > 0) {
		return echo.NewHTTPError(http.StatusBadRequest, "WEBINAR_POLL_EXIST")
	}

	webinar_poll  	:= &rdbModel.WebinarPoll{}

	if tx.Where("idx = ? ", idx ).Find(webinar_poll).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_POLL NOT FOUND")
	}

	webinar_poll.Title = payload.Title
	webinar_poll.Desc = payload.Desc
	webinar_poll.StartDate = payload.StartDate
	webinar_poll.EndDate = payload.EndDate

	tx.Save(webinar_poll)

	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")

}

