package front

import (
	//"VCMS/apps/handler"
	//"net/http"
	"github.com/labstack/echo"
	//appLibs "VCMS/apps/libs"
	rdbModel "VCMS/apps/models/rdb"
	//apiModel "VCMS/apps/models/api"
	//appConfig "VCMS/apps/config"
	"VCMS/apps/handler"
	apiModel "VCMS/apps/models/api"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type WebinarFrontQnAAPI struct {
	requestPost WebinarFrontQnAPostRequest
	requestPut  WebinarFrontQnAPutRequest
	requestAdminPost WebinarAdminQnAPostRequest
	requestAdminPut	WebinarAdminQnAPutRequest
}

type WebinarFrontQnAPostRequest struct {
	WebinarSiteId string `validate:"required" json:"webinar_site_id"`
	//QuestionType		string		`validate:"required" json:"question_type"`
	//FrontUserId		string		`validate:"required" json:"front_user_id"`
	QuestionContent   string `validate:"required" json:"question_content"`
	QuestionVideoTime string `validate:"required" json:"question_video_time"`
	//ReplyYN			string		`validate:"required" json:"reply_yn"`
}

type WebinarFrontQnAPutRequest struct {
	//FrontUserId		string		`validate:"required" json:"front_user_id"`
	QuestionContent string `validate:"required" json:"question_content"`
	//ReplyYN			string		`validate:"required" json:"reply_yn"`
}

type WebinarAdminQnAPostRequest struct {
	WebinarSiteId 		string 		`validate:"required" json:"webinar_site_id"`
	WebinarQnaId		string 		`validate:"required" json:"webinar_qna_id"`
	ReplyContent 		string		`validate:"required" json:"reply_content"`
	EmailSendYN			string		`validate:"required" json:"email_send_yn"` // Y:예, N:아니오
}

type WebinarAdminQnAPutRequest struct {
	FrontUserId		string		`validate:"required" json:"front_user_id"`
	QuestionContent		string		`validate:"required" json:"question_content"`
	//ReplyYN			string		`validate:"required" json:"reply_yn"`
}

// 웨비나 Q&A 질문 등록
func (api WebinarFrontQnAAPI) PostWebinarQnA(c echo.Context) error {

	logrus.Debug("########### PostWebinarQnA !!!")

	claims := apiModel.GetFrontJWTClaims(c)

	payload := &api.requestPost
	c.Bind(payload)

	logrus.Debug(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//claims 	 := apiModel.GetJWTClaims(c)
	//
	//logrus.Debug("#### TenantId : ", claims.TenantId)

	// TODO : 라이브 영상 시작 후 질문 시간을 계산해서 넣는다.

	// 현재시간에서 라이브 시작시간을 뺀다.
	// 그것을 초로 환산하고...그게 질문 시간이 된다.
	// 만약 혹시나 동영상 총 재생시간 보다 크다면...총 동영상 시간으로 셋팅한다.

	webinar_site := &rdbModel.WebinarSite{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("webinar_site_id = ? ", payload.WebinarSiteId).Find(webinar_site).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_SITE NOT FOUND")
	}

	// webinar_site_id로 tenant_id 구하기
	tenant_id := webinar_site.TenantId

	webinar_qna_id := uuid.NewV4().String()

	// 질문등록
	webinar_front_qna := &rdbModel.WebinarFrontQnA{
		TenantId:          tenant_id,
		WebinarSiteId:     payload.WebinarSiteId,
		WebinarQnaId:      webinar_qna_id,
		QuestionType:      "1",
		FrontUserId:       claims.MemberId,
		FrontUserName:     claims.MemberName,
		QuestionContent:   payload.QuestionContent,
		QuestionVideoTime: payload.QuestionVideoTime,
		//ReplyYN		: payload.ReplyYN,
	}

	tx.Create(webinar_front_qna)

	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")

}

// QnA 삭제
func (api WebinarFrontQnAAPI) DeleteWebinarQnA(c echo.Context) error {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	tx := c.Get("Tx").(*gorm.DB)

	webinar_front_qna := &rdbModel.WebinarFrontQnA{}
	webinar_admin_qna := &rdbModel.WebinarAdminQnA{}

	if tx.Where("idx = ? ", idx).Find(webinar_front_qna).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_FRONT_QNA NOT FOUND")
	}

	if tx.Delete(webinar_admin_qna, "webinar_site_id = ? and webinar_qna_id = ? ",
		webinar_front_qna.WebinarSiteId, webinar_front_qna.WebinarQnaId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	if tx.Delete(webinar_front_qna, "webinar_site_id = ? and idx = ? ", webinar_front_qna.WebinarSiteId, idx).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "CONTENT NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}

// Q&A 수정
func (api WebinarFrontQnAAPI) PutWebinarQnA(c echo.Context) error {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	//claims 	 := apiModel.GetFrontJWTClaims(c)

	payload := &api.requestPut
	c.Bind(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx := c.Get("Tx").(*gorm.DB)

	webinar_front_qna := &rdbModel.WebinarFrontQnA{}

	if tx.Where("idx = ? ", idx).Find(webinar_front_qna).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_FRONT_QNA NOT FOUND")
	}

	if tx.Model(webinar_front_qna).Where("idx = ? ", idx).
		Updates(rdbModel.WebinarFrontQnA{
			QuestionContent: payload.QuestionContent,
		}).RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_FRONT_QNA NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")

}

// QnA 리스트 조회
func (api WebinarFrontQnAAPI) GetWebinarQnAs(c echo.Context) error {

	claims := apiModel.GetFrontJWTClaims(c)

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	webinar_site_id := c.QueryParam("webinar_site_id")
	question_type 	:= c.QueryParam("question_type")

	member_id 		:= c.QueryParam("member_id")
	reply_yn		:= c.QueryParam("reply_yn")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	logrus.Debug("##### member_id : ", claims.MemberId)
	logrus.Debug("##### reply_yn : ", reply_yn)

	tx := c.Get("Tx").(*gorm.DB)

	webinar_front_qnas := []rdbModel.WebinarFrontQnA{}

	tx = tx.Where("webinar_site_id = ? ", webinar_site_id)

	if question_type != "" {
		tx = tx.Where("question_type = ? ", question_type)
	}

	if reply_yn != "" {
		tx = tx.Where("reply_yn = ? ", reply_yn)
	}

	if member_id != "" {
		tx = tx.Where("front_user_id = ? ", member_id)
	}

	var count = 0
	tx.Find(&webinar_front_qnas).Count(&count)

	tx.Preload("WebinarAdminQnA").
		//Order("idx desc").Offset(offset).Limit(limit).Find(&webinar_front_qnas)
		Order("created_at desc").Find(&webinar_front_qnas)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows":        webinar_front_qnas,
		})
}


// 웨비나 Q&A 답변 등록
func (api WebinarFrontQnAAPI) PostWebinarAdminQnA(c echo.Context) error {

	logrus.Debug("########### PostWebinarQnA !!!")

	payload := &api.requestAdminPost
	c.Bind(payload)

	logrus.Debug(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims 	 := apiModel.GetFrontJWTClaims(c)

	logrus.Debug("#### MemberId : ", claims.MemberId)

	webinar_front_qna := &rdbModel.WebinarFrontQnA{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("webinar_qna_id = ? ", payload.WebinarQnaId ).Find(webinar_front_qna).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_Q&A NOT FOUND")
	}

	// 답변이 이미 존재하지 않으면... UPDATE 존재하지 않으면 INSERT 한다.
	if webinar_front_qna.ReplyYN == "N" {
		// 답변 등록
		webinar_admin_qna := &rdbModel.WebinarAdminQnA{
			TenantId	: webinar_front_qna.TenantId,
			WebinarSiteId	: webinar_front_qna.WebinarSiteId,
			WebinarQnaId	: payload.WebinarQnaId,
			ReplyContent	: payload.ReplyContent,
			EmailSendYN	: payload.EmailSendYN,
			UpdatedId: claims.MemberId,
			//ReplyYN		: payload.ReplyYN,
		}
		tx.Create(webinar_admin_qna)
	} else {
		webinar_admin_qna := &rdbModel.WebinarAdminQnA{}
		// 존재하면 UPDATE
		if tx.Model(webinar_admin_qna).Where("webinar_qna_id = ? ", payload.WebinarQnaId).
			Updates(rdbModel.WebinarAdminQnA{
			ReplyContent	: payload.ReplyContent,
			EmailSendYN		: payload.EmailSendYN,
			UpdatedId		: claims.MemberId,
		}).RowsAffected == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_FRONT_QNA NOT FOUND")
		}
	}

	// 일단 답변 여부 "Y"로 update 한다.
	webinar_front_qna.ReplyYN = "Y";

	tx.Save(webinar_front_qna)

	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")
}

// 답변 삭제
func (api WebinarFrontQnAAPI) DeleteWebinarAdminQnA(c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	tx := c.Get("Tx").(*gorm.DB)

	webinar_admin_qna	:= &rdbModel.WebinarAdminQnA{}
	webinar_front_qna	:= &rdbModel.WebinarFrontQnA{}

	if tx.Where("idx = ? ", idx ).Find(webinar_admin_qna).RecordNotFound() {

	}

	webinar_qna_id := webinar_admin_qna.WebinarQnaId

	if tx.Model(webinar_front_qna).Where("webinar_qna_id = ? ", webinar_qna_id).
		Updates(rdbModel.WebinarFrontQnA{
		ReplyYN:"N",
	}).RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_FRONT_QNA NOT FOUND")
	}

	if tx.Delete(webinar_admin_qna, "idx = ? ", idx).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}