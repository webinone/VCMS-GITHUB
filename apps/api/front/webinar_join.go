package front

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
	"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
)


type WebinarFrontJoinAPI struct {
	requestPost WebinarFrontJoinPostRequest
}

type WebinarFrontJoinPostRequest struct {
	WebinarSiteId		string 		`validate:"required" json:"webinar_site_id"`
	MobilePhoneNum		string		`validate:"required" json:"mobile_phone_num"`
	Email			string		`validate:"required" json:"email"`
	CompanyName		string		`validate:"required" json:"company_name"`
	Department		string		`validate:"required" json:"department"`
	Position		string		`validate:"required" json:"position"`
	PhoneNum		string		`validate:"required" json:"phone_num"`
	ZipCode			string		`validate:"required" json:"zip_code"`
	Address1		string		`validate:"required" json:"address1"`
	Address2		string		`validate:"required" json:"address2"`
}

// 웨비나 프론트 참여 JOIN
func (api WebinarFrontJoinAPI) PostWebinarJoin(c echo.Context) error {

	logrus.Debug("########### PostWebinarFrontJoin !!!")

	claims := apiModel.GetFrontJWTClaims(c)

	payload := &api.requestPost
	c.Bind(payload)

	logrus.Debug(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//claims 	 := apiModel.GetJWTClaims(c)
	//logrus.Debug("#### TenantId : ", claims.TenantId)
	webinar_join     := &rdbModel.WebinarJoin{}
	webinar_site     := &rdbModel.WebinarSite{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("webinar_site_id = ? ", payload.WebinarSiteId ).Find(webinar_site).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_SITE NOT FOUND")
	}

	// webinar_site_id로 tenant_id 구하기
	tenant_id := webinar_site.TenantId

	if !tx.Where("tenant_id = ?  AND webinar_site_id = ? AND front_user_id = ? ", tenant_id, payload.WebinarSiteId, claims.MemberId ).
		Find(webinar_join).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_JOIN_EXISTS")
	}

	webinar_join.TenantId 		= tenant_id
	webinar_join.FrontUserId 	= claims.MemberId
	webinar_join.FrontUserName 	= claims.MemberName
	webinar_join.WebinarSiteId 	= payload.WebinarSiteId
	webinar_join.WebinarJoinId 	= uuid.NewV4().String()
	webinar_join.MobilePhoneNum = payload.MobilePhoneNum
	webinar_join.Email 	    	= payload.Email
	webinar_join.CompanyName    = payload.CompanyName
	webinar_join.Department     = payload.Department
	webinar_join.CompanyName    = payload.CompanyName
	webinar_join.Position       = payload.Position
	webinar_join.PhoneNum       = payload.PhoneNum
	webinar_join.ZipCode        = payload.ZipCode
	webinar_join.Address1       = payload.Address1
	webinar_join.Address2       = payload.Address2

	tx.Create(webinar_join)

	// WEBINAR JOIN

	webinar_join_member := &rdbModel.WebinarJoinMember{}
	webinar_join_member.TenantId 		= tenant_id
	webinar_join_member.FrontUserId 	= claims.MemberId
	webinar_join_member.FrontUserName 	= claims.MemberName
	webinar_join_member.MobilePhoneNum  = payload.MobilePhoneNum
	webinar_join_member.Email			= payload.Email
	webinar_join_member.CompanyName 	= payload.CompanyName
	webinar_join_member.Department 		= payload.Department
	webinar_join_member.Position 		= payload.Position
	webinar_join_member.PhoneNum       	= payload.PhoneNum
	webinar_join_member.ZipCode        	= payload.ZipCode
	webinar_join_member.Address1       	= payload.Address1
	webinar_join_member.Address2       	= payload.Address2

	tx.Save(webinar_join_member)

	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")

}


// 웨비나 프론트 참여 JOIN
func (api WebinarFrontJoinAPI) GetWebinarJoinCheck(c echo.Context) error {

	claims := apiModel.GetFrontJWTClaims(c)

	///api/front/webinar/join/check?webinar_site_id=f4dbaf32-ed6b-4499-8c6f-9e75b2317c5e

	webinar_site_id  := c.QueryParam("webinar_site_id")

	webinar_join     	:= &rdbModel.WebinarJoin{}
	webinar_join_member	:= &rdbModel.WebinarJoinMember{}
	member_default	 	:= &rdbModel.MemberDefault{}

	helloTtx, err := gorm.Open(appConfig.Config.RDB[1].Product, appConfig.Config.RDB[1].ConnectString)
	if err != nil {
		panic(err)
	}
	defer helloTtx.Close()
	helloTtx.LogMode(appConfig.Config.RDB[1].Debug)

	// 회원 기본 정보
	if helloTtx.Preload("MemberSub").
		Where("member_id = ? ", claims.MemberId ).
		Find(member_default).RecordNotFound() {
	}

	tx := c.Get("Tx").(*gorm.DB)

	// 회원 웨비나 참여 회원 정보
	if tx.Where("front_user_id = ? ", claims.MemberId ).
		Find(webinar_join_member).RecordNotFound() {
	}

	if tx.Where("webinar_site_id = ? AND front_user_id = ? ", webinar_site_id, claims.MemberId ).
		Find(webinar_join).RecordNotFound() {

		// 존재하지 않는다.
		// 웨비나 사용 가입 페이지로 이동해야 한다.
		return handler.APIResultHandler(c, true, http.StatusOK,
			map[string]interface{}{
				"webinar_site_id" : webinar_site_id,
				"message": "NOT_EXIST",
				"member": member_default,
				"join_member": webinar_join_member,
			})
	} else {
		return handler.APIResultHandler(c, true, http.StatusOK,
			map[string]interface{}{
				"webinar_site_id" : webinar_site_id,
				"message": "EXIST",
				"member": member_default,
				"join_member": webinar_join_member,
			})

	}
}


