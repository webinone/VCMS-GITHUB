package api

import (
	"github.com/labstack/echo"
	"strconv"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"VCMS/apps/handler"
	apiModel "VCMS/apps/models/api"
	rdbModel "VCMS/apps/models/rdb"
	appConfig "VCMS/apps/config"
	"net/http"
	"fmt"
	"os"
	"github.com/satori/go.uuid"
	"github.com/xuri/excelize"
	"github.com/metakeule/fmtdate"
)

type WebinarPollMemberAPI struct {
	requestPost WebinarPollMemberPostRequest
	requestPut  WebinarPollMemberPutRequest
}
//
type WebinarPollMemberPostRequest struct {
//	WebinarSiteId		string 				`validate:"required" json:"webinar_site_id"`
//	Title			string				`validate:"required" json:"title"`
//	Desc			string				`validate:"required" json:"desc"`
//	StartDate		string				`validate:"required" json:"start_date"`
//	EndDate			string				`validate:"required" json:"end_date"`
}

type WebinarPollMemberPutRequest struct {
	WebinarSiteId		string 			`validate:"required" json:"webinar_site_id"`
	WebinarPollId		string			`validate:"required" json:"webinar_poll_id"`
	FrontUserId		string				`validate:"required" json:"front_user_id"`
	WinYN		string 					`validate:"required" json:"win_yn"`


}
////
////type ContentRequestPut struct {
////	CategoryId		string		`validate:"required" json:"category_id"`
////	ContentName 		string 		`validate:"required" json:"content_name"`
////	ThumbChange		string		`validate:"required" json:"thumb_change"`
////	ThumbType		string		`json:"thumb_type"`
////	ThumbTime		string		`json:"thumb_time"`
////
////}
//
//
//
//// 웨비나 사이트 등록
//func (api WebinarPollAPI) PostWebinarPoll(c echo.Context) error {
//
//	logrus.Debug("########### PostWebinarPoll !!!")
//
//	payload := &api.requestPost
//	c.Bind(payload)
//
//	logrus.Debug(payload)
//
//	if err := c.Validate(payload); err != nil {
//		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
//	}
//
//	claims 	 := apiModel.GetJWTClaims(c)
//
//	logrus.Debug("#### TenantId : ", claims.TenantId)
//
//	tx := c.Get("Tx").(*gorm.DB)
//
//	webinar_poll_id := uuid.NewV4().String()
//
//	// TODO : 설문기간이 겹치는 설문은 등록할 수 없다.
//
//	// 설문등록
//	webinar_poll := &rdbModel.WebinarPoll{
//		TenantId: claims.TenantId,
//		WebinarSiteId: payload.WebinarSiteId,
//		WebinarPollId:webinar_poll_id,
//		Title:payload.Title,
//		Desc:payload.Desc,
//		StartDate:payload.StartDate,
//		EndDate:payload.EndDate,
//		UpdatedId:claims.UserId,
//	}
//
//	tx.Create(webinar_poll)
//
//	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")
//
//}
//
//설문 참여자 리스트 조회
func (api WebinarPollMemberAPI) GetWebinarPollMembers (c echo.Context) error  {

	claims 	 		:= apiModel.GetJWTClaims(c)

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	tenant_id		:= claims.TenantId

	webinar_site_id		:= c.QueryParam("webinar_site_id")
	webinar_poll_id		:= c.QueryParam("webinar_poll_id")
	win_yn 			:= c.QueryParam("win_yn")		// 당첨 여부

	member_info		:= c.QueryParam("member_info")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	logrus.Debug("##### webinar_poll_id : ", webinar_poll_id)
	logrus.Debug("##### win_yn : ", win_yn)

	webinar_poll_members  := []rdbModel.WebinarPollMember{}

	tx := c.Get("Tx").(*gorm.DB)

	if webinar_site_id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "webinar_site_id required parameter")
	}

	if webinar_poll_id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "webinar_poll_id required parameter")
	}

	tx  = tx.Where("tenant_id = ? and webinar_site_id = ? and webinar_poll_id = ? ", tenant_id, webinar_site_id, webinar_poll_id)

	if member_info != "" {
		tx = tx.Where("front_user_id LIKE ? OR  front_user_name LIKE ? ", "%" + member_info + "%", "%" + member_info + "%")
	}

	if win_yn != "" {
		tx  = tx.Where("win_yn = ? ", win_yn)

	}

	var count = 0
	tx.Find(&webinar_poll_members).Count(&count)

	tx.Preload("WebinarPollMemberResults").
		Preload("WebinarPollMemberResults.WebinarPollQuestionMaster").
		Preload("WebinarPollMemberResults.WebinarPollQuestionMaster.WebinarPollQuestionDetails").
		Order("idx desc").Offset(offset).Limit(limit).Find(&webinar_poll_members)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_poll_members,
		})

}

// 설문 당첨 여부 수정
func (api WebinarPollMemberAPI) PutWebinarPollMemberWinYN (c echo.Context) error  {

	payload := &api.requestPut
	c.Bind(payload)

	if err   := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx := c.Get("Tx").(*gorm.DB)

	webinar_poll_member  	:= &rdbModel.WebinarPollMember{}

	if tx.Where("webinar_site_id = ? and webinar_poll_id = ? and front_user_id = ? ",
		payload.WebinarSiteId,
		payload.WebinarPollId,
		payload.FrontUserId,
	).Find(webinar_poll_member).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_POLL_MEMBER NOT FOUND")
	}

	webinar_poll_member.WinYN = payload.WinYN

	tx.Save(webinar_poll_member)

	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")

}

type WebinarPollMemberViewResult struct {
	Idx        						int64       	`json:"idx"`
	TenantID   						string      	`json:"tenant_id"`
	WebinarSiteId					string 			`json:"webinar_site_id"`
	WebinarPollId					string 			`json:"webinar_poll_id"`
	WebinarPollQuestionMasterId		string 			`json:"webinar_poll_question_master_id"`
	WebinarPollQuestionDetailId		string 			`json:"webinar_poll_question_detail_id"`
	ResultCount						string			`json:"result_count"`
	ResultPercent					string			`json:"result_percent"`
}

// 설문 통계
func (api WebinarPollMemberAPI) GetWebinarPollMemberStatistics (c echo.Context) error  {


	claims 	 		:= apiModel.GetJWTClaims(c)

	tenant_id		:= claims.TenantId

	webinar_site_id		:= c.QueryParam("webinar_site_id")
	webinar_poll_id		:= c.QueryParam("webinar_poll_id")

	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	logrus.Debug("##### webinar_poll_id : ", webinar_poll_id)

	webinar_poll_question_masters := []rdbModel.WebinarPollQuestionMaster{}
	webinar_poll_results  := []WebinarPollMemberViewResult{}

	tx := c.Get("Tx").(*gorm.DB)

	// 먼저 질문 마스터를 구한다.
	tx.Where("tenant_id = ? ", tenant_id).
		Where("webinar_site_id = ? ", webinar_site_id).
		Where("webinar_poll_id = ? ", webinar_poll_id).
		Preload("WebinarPollQuestionDetails").
		Order("idx asc").Find(&webinar_poll_question_masters)


	tx.Raw(`
		SELECT  a.idx,
				a.tenant_id,
				a.webinar_site_id,
				a.webinar_poll_id,
				a.webinar_poll_question_master_id,
				a.webinar_poll_question_detail_id,
				a.title,
				(
				  SELECT COUNT(*) FROM TB_WEBINAR_POLL_MEMBER_RESULT
				  WHERE webinar_poll_question_master_id = a.webinar_poll_question_master_id
				  AND answer = a.webinar_poll_question_detail_id
				) AS result_count,
				FLOOR(
				  (
				  SELECT COUNT(*) FROM TB_WEBINAR_POLL_MEMBER_RESULT
				  WHERE webinar_poll_question_master_id = a.webinar_poll_question_master_id
				  AND answer = a.webinar_poll_question_detail_id
				  )
				  /
				  (
				  SELECT COUNT(*) FROM TB_WEBINAR_POLL_MEMBER_RESULT
				  WHERE webinar_poll_question_master_id = a.webinar_poll_question_master_id
				  ) * 100
				) AS result_percent
			  FROM TB_WEBINAR_POLL_QUESTION_DETAIL AS a
		WHERE tenant_id = ?
		AND webinar_site_id = ?
		AND webinar_poll_id = ?
		AND deleted_at is null
		`, tenant_id, webinar_site_id, webinar_poll_id).Scan(&webinar_poll_results)


	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"questions": webinar_poll_question_masters,
			"results": webinar_poll_results,
		})

}

// 설문결과 엑셀 다운로드
func (api WebinarPollMemberAPI) DownloadExcelFile(c echo.Context) error {

	logrus.Debug("******** DownloadExcelFile ")

	webinar_site_id					:= c.QueryParam("webinar_site_id")
	webinar_poll_id     			:= c.QueryParam("webinar_poll_id")
	//webinar_poll_question_master_id := c.QueryParam("webinar_poll_question_master_id")

	webinar_poll_question_masters   := []rdbModel.WebinarPollQuestionMaster{}
	webinar_poll_members   			:= []rdbModel.WebinarPollMember{}

	tx := c.Get("Tx").(*gorm.DB)

	tx  = tx.Where("webinar_site_id = ? ",  webinar_site_id)
	tx  = tx.Where("webinar_poll_id = ? ",  webinar_poll_id)
	//tx  = tx.Where("webinar_poll_question_master_id = ? ",  webinar_poll_question_master_id)

	tx.Order("idx asc").Find(&webinar_poll_question_masters)

	tx.Preload("WebinarPollMemberResults").
		Preload("WebinarPollMemberResults.WebinarPollQuestionDetail").
		Order("idx desc").Find(&webinar_poll_members)

	xlsx_name := uuid.NewV4().String() + ".xlsx"

	logrus.Debug("xlsx_name : ", xlsx_name)

	xlsx := excelize.NewFile()

	xlsx.SetCellValue("Sheet1", "A1", "사용자ID")

	xlsx.SetCellValue("Sheet1", "B1", "사용자명")

	xlsx.SetCellValue("Sheet1", "C1", "당첨여부")

	xlsx.SetCellValue("Sheet1", "D1", "설문참여일시")

	index := 4
	colum_location := ""
	colum_title := ""

	for i, v := range webinar_poll_question_masters {

		colum_location = toCharStrArr(index) + "1"
		colum_title = v.Title

		fmt.Println(" i : " , i)

		xlsx.SetCellValue("Sheet1", colum_location, colum_title)

		index++

	}

	// HEADER 끝.....
	for i, v := range webinar_poll_members {

		//colum_location = toCharStrArr(i)

		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i + 2), v.FrontUserId)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i + 2), v.FrontUserName)
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(i + 2), v.WinYN)
		xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(i + 2), fmtdate.Format("YYYY-MM-DD hh:mm:ss", v.CreatedAt))

		//for z,
		result_index := 4
		for _, z := range v.WebinarPollMemberResults {

			colum_location = toCharStrArr(result_index)
			colum_title = z.WebinarPollQuestionDetail.Title

			if colum_title == "" {
				colum_title = z.Answer
			}

			xlsx.SetCellValue("Sheet1", colum_location + strconv.Itoa(i + 2), colum_title)

			result_index++
		}

	}

	// Save xlsx file by the given path.
	err := xlsx.SaveAs(appConfig.Config.APP.ExcelRoot + "/" + xlsx_name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	return c.Attachment(appConfig.Config.APP.ExcelRoot + "/" + xlsx_name, xlsx_name)

}

var arr = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
					  "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func toCharStrArr(i int) string {
	return arr[i]
}




