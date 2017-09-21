package api

import (
	//"VCMS/apps/handler"
	//"net/http"
	"github.com/labstack/echo"
	//appLibs "VCMS/apps/libs"
	rdbModel "VCMS/apps/models/rdb"
	apiModel "VCMS/apps/models/api"
	appConfig "VCMS/apps/config"
	"net/http"
	"github.com/Sirupsen/logrus"
	"VCMS/apps/handler"
	"github.com/jinzhu/gorm"
	"strconv"
	//"github.com/satori/go.uuid"
	//"github.com/xuri/excelize"
	//"github.com/metakeule/fmtdate"
	//"fmt"
	//"os"
	//"github.com/xuri/excelize"
	//"fmt"
	//"os"
	//"github.com/metakeule/fmtdate"
	"github.com/satori/go.uuid"
	"github.com/xuri/excelize"
	"github.com/metakeule/fmtdate"
	"fmt"
	"os"
)


type WebinarAdminJoinAPI struct {
	requestPost WebinarAdminJoinPostRequest
}


type WebinarAdminJoinPostRequest struct {
	WebinarSiteId		string 		`validate:"required" json:"webinar_site_id"`
	FrontUserId		string		`validate:"required" json:"memer_id"`
}

//type WebinarFrontQnAPutRequest struct {
//	FrontUserId		string		`validate:"required" json:"front_user_id"`
//	QuestionContent		string		`validate:"required" json:"question_content"`
//	//ReplyYN			string		`validate:"required" json:"reply_yn"`
//}

// 공지사항 리스트 조회
func (api WebinarAdminJoinAPI) GetWebinarJoins (c echo.Context) error  {

	claims 	 := apiModel.GetJWTClaims(c)

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	webinar_site_id		:= c.QueryParam("webinar_site_id");
	created_at		:= c.QueryParam("created_at")

	member_info		:= c.QueryParam("member_info")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	logrus.Debug("##### created_at : ", created_at)

	logrus.Debug("##### member_info : ", member_info)

	tx := c.Get("Tx").(*gorm.DB)

	webinar_joins    := []rdbModel.WebinarJoin{}

	tx  = tx.Where("tenant_id = ? and webinar_site_id = ? ", claims.TenantId, webinar_site_id)

	if created_at !=  "" {
		tx = tx.Where("created_at LIKE ? ", created_at + "%")
	}

	if member_info != "" {
		tx = tx.Where("front_user_id LIKE ? OR  front_user_name LIKE ? ", "%" + member_info + "%", "%" + member_info + "%")
	}

	var count = 0
	tx.Find(&webinar_joins).Count(&count)

	tx.Order("idx desc").Offset(offset).Limit(limit).Find(&webinar_joins)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_joins,
		})

}

// 자료 다운로드
func (api WebinarAdminJoinAPI) DownloadExcelFile(c echo.Context) error {

	logrus.Debug("******** DownloadNoticeFile ")

	//claims 	 := apiModel.GetJWTClaims(c)

	webinar_site_id		:= c.QueryParam("webinar_site_id");

	created_at		:= c.QueryParam("created_at")

	member_info		:= c.QueryParam("member_info")

	//logrus.Debug("##### offset : ", offset)
	//logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	//logrus.Debug("##### start_date : ", start_date)
	//logrus.Debug("##### end_date : ", end_date)
	logrus.Debug("##### created_at : ", created_at)


	tx := c.Get("Tx").(*gorm.DB)

	webinar_joins    := []rdbModel.WebinarJoin{}

	tx  = tx.Where("webinar_site_id = ? ",  webinar_site_id)

	if created_at !=  "" {
		tx = tx.Where("created_at LIKE ? ", created_at + "%")
	}


	if member_info != "" {
		tx = tx.Where("front_user_id LIKE ? OR  front_user_name LIKE ? ", "%" + member_info + "%", "%" + member_info + "%")
	}


	tx.Order("idx desc").Find(&webinar_joins)

	xlsx_name := uuid.NewV4().String() + ".xlsx"

	logrus.Debug("xlsx_name : ", xlsx_name)

	xlsx := excelize.NewFile()

	xlsx.SetCellValue("Sheet1", "A1", "사용자ID")
	// 회원명
	xlsx.SetCellValue("Sheet1", "B1", "사용자명")
	// 질문일시
	xlsx.SetCellValue("Sheet1", "C1", "휴대폰")
	// 동영상 시간
	xlsx.SetCellValue("Sheet1", "D1", "이메일")
	// 질문내용
	xlsx.SetCellValue("Sheet1", "E1", "회사명")

	xlsx.SetCellValue("Sheet1", "F1", "부서")

	xlsx.SetCellValue("Sheet1", "G1", "직급 또는 직책")

	xlsx.SetCellValue("Sheet1", "H1", "일반전화번호")

	xlsx.SetCellValue("Sheet1", "I1", "주소")

	xlsx.SetCellValue("Sheet1", "J1", "상세주소")

	xlsx.SetCellValue("Sheet1","K1", "참여일시")
	//
	for i, v := range webinar_joins {

		// 회원ID
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i + 2), v.FrontUserId)
		// 회원명
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i + 2), v.FrontUserName)
		// 휴대폰
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(i + 2), v.MobilePhoneNum)
		// 이메일
		xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(i + 2), v.Email)
		// 회사명
		xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(i + 2), v.CompanyName)
		// 부서명
		xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(i + 2), v.Department)
		// 직급또는직책
		xlsx.SetCellValue("Sheet1", "G"+strconv.Itoa(i + 2), v.Position)
		// 일반전화번호
		xlsx.SetCellValue("Sheet1", "H"+strconv.Itoa(i + 2), v.PhoneNum)
		// 주소
		xlsx.SetCellValue("Sheet1", "I"+strconv.Itoa(i + 2), v.Address1)
		// 상세주소
		xlsx.SetCellValue("Sheet1", "J"+strconv.Itoa(i + 2), v.Address2)
		// 참여일시
		xlsx.SetCellValue("Sheet1", "K"+strconv.Itoa(i + 2), fmtdate.Format("YYYY-MM-DD hh:mm:ss", v.CreatedAt))

	}
	//
	// Save xlsx file by the given path.
	err := xlsx.SaveAs(appConfig.Config.APP.ExcelRoot + "/" + xlsx_name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	return c.Attachment(appConfig.Config.APP.ExcelRoot + "/" + xlsx_name, xlsx_name)

}

