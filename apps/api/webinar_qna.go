package api

import (
	"github.com/labstack/echo"
	appLibs "VCMS/apps/libs"
	rdbModel "VCMS/apps/models/rdb"
	apiModel "VCMS/apps/models/api"
	appConfig "VCMS/apps/config"
	"net/http"
	"github.com/Sirupsen/logrus"
	"VCMS/apps/handler"
	"github.com/jinzhu/gorm"
	"strconv"
	"github.com/metakeule/fmtdate"
	"fmt"
	"github.com/xuri/excelize"
	"os"
	"github.com/satori/go.uuid"
)

type WebinarAdminQnAAPI struct {
	requestPost WebinarAdminQnAPostRequest
	requestPut  WebinarAdminQnAPutRequest
}


type WebinarAdminQnAPostRequest struct {
	WebinarQnaId		string 		`validate:"required" json:"webinar_qna_id"`
	ReplyContent 		string		`validate:"required" json:"reply_content"`
	EmailSendYN		string		`validate:"required" json:"email_send_yn"` // Y:예, N:아니오
}

type WebinarAdminQnAPutRequest struct {
	FrontUserId		string		`validate:"required" json:"front_user_id"`
	QuestionContent		string		`validate:"required" json:"question_content"`
	//ReplyYN			string		`validate:"required" json:"reply_yn"`
}

// 웨비나 사이트 등록
func (api WebinarAdminQnAAPI) PostWebinarQnA(c echo.Context) error {

	logrus.Debug("########### PostWebinarQnA !!!")

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
			//ReplyYN		: payload.ReplyYN,
		}
		tx.Create(webinar_admin_qna)
	} else {
		webinar_admin_qna := &rdbModel.WebinarAdminQnA{}
		// 존재하면 UPDATE
		if tx.Model(webinar_admin_qna).Where("webinar_qna_id = ? ", payload.WebinarQnaId).
			Updates(rdbModel.WebinarAdminQnA{
			ReplyContent	: payload.ReplyContent,
			EmailSendYN	: payload.EmailSendYN,
		}).RowsAffected == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_FRONT_QNA NOT FOUND")
		}
	}

	// 일단 답변 여부 "Y"로 update 한다.
	webinar_front_qna.ReplyYN = "Y";

	tx.Save(webinar_front_qna)

	fmt.Println(" !!!! payload.EmailSendYN : ", payload.EmailSendYN)

	// EMAIL 발송
	if payload.EmailSendYN == "Y" {
		// 발송여부가 Y 인 경우만
		// 사용자 메일 주소 구하기
		member_default := &rdbModel.MemberDefault{}
		tx.Where("member_id = ? ", webinar_front_qna.FrontUserId ).Find(member_default)

		member_id   	:= member_default.MemberId
		member_name 	:= member_default.MemberName
		member_email 	:= member_default.EmailAddr

		// 메일이 있는 경우만...
		if member_email != "" {
			webinar_site     := &rdbModel.WebinarSite{}
			if tx.Where("webinar_site_id = ? ", webinar_front_qna.WebinarSiteId ).Find(webinar_site).RecordNotFound() {

				return echo.NewHTTPError(http.StatusNotFound, "WEBINAR NOT FOUND")
			}

			send_mail_client := &appLibs.SendMailClient {
				SmtpURL		: appConfig.Config.EMAIL.SmtpUrl,
				SmtpPort	: appConfig.Config.EMAIL.SmtpPort,
				User		: appConfig.Config.EMAIL.User,
				Password	: appConfig.Config.EMAIL.Password,
				From		: appConfig.Config.EMAIL.User,
				To			: member_email,
				Subject		: "[웨비나 질문/답변] " + webinar_site.Title,
			}

			send_mail_client.SendWebinarQnAMail(
				member_id,
				member_name,
				webinar_front_qna.QuestionContent,
				fmtdate.Format("YYYY-MM-DD hh:mm:ss", webinar_front_qna.CreatedAt),
				payload.ReplyContent,
				"2017-10-11 12:30:11",
				"http://webinar.hellot.net/site/" + webinar_site.WebinarSiteId,
				webinar_site.Title,
				"2017-12-11 12:30:11",
			)

		} else {

			// 메일이 없는 경우에는 보내지 않는다.
			logrus.Error("메일이 없어서 보내지 않는다.")
		}
	}


	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")

}

// QnA 삭제
func (api WebinarAdminQnAAPI) DeleteWebinarQnA(c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	tx := c.Get("Tx").(*gorm.DB)

	webinar_front_qna	:= &rdbModel.WebinarFrontQnA{}
	webinar_admin_qna	:= &rdbModel.WebinarAdminQnA{}

	if tx.Where("idx = ? ", idx ).Find(webinar_front_qna).RecordNotFound() {
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

// QnA 리스트 조회
func (api WebinarAdminQnAAPI) GetWebinarQnAs (c echo.Context) error  {

	claims 	 := apiModel.GetJWTClaims(c)

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	webinar_site_id		:= c.QueryParam("webinar_site_id");

	question_type		:= c.QueryParam("question_type")
	question_content	:= c.QueryParam("question_content")
	reply_yn		:= c.QueryParam("reply_yn")
	created_at		:= c.QueryParam("created_at")

	member_info		:= c.QueryParam("member_info")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	//logrus.Debug("##### start_date : ", start_date)
	//logrus.Debug("##### end_date : ", end_date)
	logrus.Debug("##### question_content : ", question_content)
	logrus.Debug("##### reply_yn : ", reply_yn)
	logrus.Debug("##### created_at : ", created_at)


	tx := c.Get("Tx").(*gorm.DB)

	webinar_front_qnas    := []rdbModel.WebinarFrontQnA{}

	tx  = tx.Where("tenant_id = ? and webinar_site_id = ? ", claims.TenantId, webinar_site_id)

	if question_type != "" {
		tx = tx.Where("question_type = ? ", question_type)
	}

	if reply_yn != "" {
		tx = tx.Where("reply_yn = ? ", reply_yn)
	}

	if created_at !=  "" {
		tx = tx.Where("created_at LIKE ? ", created_at + "%")
	}

	if question_content != "" {
		tx = tx.Where("question_content LIKE ? ", "%" + question_content + "%")
	}

	if member_info != "" {
		tx = tx.Where("front_user_id LIKE ? OR  front_user_name LIKE ? ", "%" + member_info + "%", "%" + member_info + "%")
	}


	var count = 0
	tx.Find(&webinar_front_qnas).Count(&count)

	tx.Preload("WebinarAdminQnA").
		//Preload("MemberDefault").
		//Preload("MemberDefault.MemberSub").
		Preload("WebinarJoin").Order("idx desc").Offset(offset).Limit(limit).Find(&webinar_front_qnas)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_front_qnas,
		})

}


// GetWebinarQnA

func (api WebinarAdminQnAAPI) GetWebinarQnA (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	webinar_front_qna  := &rdbModel.WebinarFrontQnA{}

	//var count = 0
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("idx = ? ", idx ).Preload("WebinarAdminQnA").
		Preload("MemberDefault").Preload("MemberDefault.MemberSub").
		Find(webinar_front_qna).RecordNotFound() {

		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR NOTICE  NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, webinar_front_qna)

}


// 자료 다운로드
func (api WebinarAdminQnAAPI) DownloadExcelFile(c echo.Context) error {

	logrus.Debug("******** DownloadNoticeFile ")


	//claims 	 := apiModel.GetJWTClaims(c)

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	webinar_site_id		:= c.QueryParam("webinar_site_id");

	question_type		:= c.QueryParam("question_type")
	question_content	:= c.QueryParam("question_content")
	reply_yn		:= c.QueryParam("reply_yn")
	created_at		:= c.QueryParam("created_at")

	member_info		:= c.QueryParam("member_info")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	//logrus.Debug("##### start_date : ", start_date)
	//logrus.Debug("##### end_date : ", end_date)
	logrus.Debug("##### question_content : ", question_content)
	logrus.Debug("##### reply_yn : ", reply_yn)
	logrus.Debug("##### created_at : ", created_at)


	tx := c.Get("Tx").(*gorm.DB)

	webinar_front_qnas    := []rdbModel.WebinarFrontQnA{}

	tx  = tx.Where("webinar_site_id = ? ",  webinar_site_id)

	if question_type != "" {
		tx = tx.Where("question_type = ? ", question_type)
	}

	if reply_yn != "" {
		tx = tx.Where("reply_yn = ? ", reply_yn)
	}

	if created_at !=  "" {
		tx = tx.Where("created_at LIKE ? ", created_at + "%")
	}

	if question_content != "" {
		tx = tx.Where("question_content LIKE ? ", "%" + question_content + "%")
	}

	if member_info != "" {
		tx = tx.Where("front_user_id LIKE ? OR  front_user_name LIKE ? ", "%" + member_info + "%", "%" + member_info + "%")
	}


	var count = 0
	tx.Find(&webinar_front_qnas).Count(&count)

	tx.Preload("WebinarAdminQnA").
		Preload("WebinarJoin").Order("idx desc").Offset(offset).Limit(limit).Find(&webinar_front_qnas)

	xlsx_name := uuid.NewV4().String() + ".xlsx"

	xlsx := excelize.NewFile()

	xlsx.SetCellValue("Sheet1", "A1", "질문자ID")
	// 회원명
	xlsx.SetCellValue("Sheet1", "B1", "질문자명")
	// 질문일시
	xlsx.SetCellValue("Sheet1", "C1", "질문일시")
	// 동영상 시간
	xlsx.SetCellValue("Sheet1", "D1", "질문동영상시간")
	// 질문내용
	xlsx.SetCellValue("Sheet1", "E1", "질문내용")

	xlsx.SetCellValue("Sheet1", "F1", "답변내용")

	for i, v := range webinar_front_qnas {

		// 회원ID
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i + 2), v.FrontUserId)
		// 회원명
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i + 2), v.FrontUserName)
		// 질문일시
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(i + 2), fmtdate.Format("YYYY-MM-DD hh:mm:ss", v.CreatedAt))
		// 동영상 시간
		xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(i + 2), v.QuestionVideoTime)
		// 질문내용
		xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(i + 2), v.QuestionContent)
		// 답변내용
		xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(i + 2), v.WebinarAdminQnA.ReplyContent)

	}

	// Save xlsx file by the given path.
	err := xlsx.SaveAs(appConfig.Config.APP.ExcelRoot + "/" + xlsx_name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	return c.Attachment(appConfig.Config.APP.ExcelRoot + "/" + xlsx_name, xlsx_name)
}