package api

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"VCMS/apps/handler"
	rdbModel "VCMS/apps/models/rdb"
	"strconv"
	appLibs "VCMS/apps/libs"
	appConfig "VCMS/apps/config"
	"path/filepath"
	"os"
	"github.com/satori/go.uuid"
)

type TenantAPI struct {
	request TenantRequest
}

type TenantRequest struct {
	TenantId   string        `json:"tenant_id"`
	TenantName string        `validate:"required" json:"name"`
	TenantDesc string        `json:"tenant_desc"`
	TimeZone   string        `validate:"required" json:"time_zone"`
}

// 테넌트 등록
func (api TenantAPI) PostTenant (c echo.Context) error  {

	payload := &api.request
	c.Bind(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 암호전 패스워드
	logrus.Debug("TenantId  : ", payload.TenantId)
	logrus.Debug("TenantName  : ", payload.TenantName)
	logrus.Debug("TenantDesc  : ", payload.TenantDesc)
	logrus.Debug("TimeZone  : ", payload.TimeZone)

	tenant := &rdbModel.Tenant{
		TenantId  	:payload.TenantId,
		Name	  	:payload.TenantName,
		TenantDesc 	:payload.TenantDesc,
		TimeZone  	:payload.TimeZone,
	}

	// 기존에 등록되어 있는 Tenant인지 체크 한다.
	tx := c.Get("Tx").(*gorm.DB)
	if !tx.Where("tenant_id = ? ",
		tenant.TenantId ).Find(tenant).RecordNotFound() {
		return echo.NewHTTPError(http.StatusConflict, "Already exists Tenant ")
	}

	// TODO : 어플리케이션을 2개 생성한다.

	// TODO : Live 어플리케이션
	_, err := appLibs.WowzaHttpClient{APIUrl:appConfig.Config.WOWZA.ApiUrl}.CreateApplication(payload.TenantId, "Live", "live")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Wowza Server Error : " + err.Error())
	}

	// TODO : VOD 어플리케이션
	//_, err = appLibs.WowzaHttpClient{APIUrl:appConfig.Config.WOWZA.ApiUrl}.CreateApplication(payload.TenantId, "VOD", "default")
	//if err != nil {
	//	return echo.NewHTTPError(http.StatusInternalServerError, "Wowza Server Error : " + err.Error())
	//}


	// 디렉토리 권한 주기
	//contentDir 	:= appConfig.Config.APP.ContentsRoot + string(filepath.Separator) + payload.TenantId
	//// content 디렉토리 권한 주기
	//err = os.Chmod(contentDir, 0777)
	//if err != nil {
	//	return echo.NewHTTPError(http.StatusInternalServerError, "Wowza Server Error : " + err.Error())
	//}

	 //VOD 디렉토리 생성 content/vod/application명/
	vodDir 		:= appConfig.Config.APP.ContentsRoot + string(filepath.Separator) + "vod" + string(filepath.Separator) + payload.TenantId
	if _, err := os.Stat(vodDir); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			err := os.Mkdir(vodDir,0777);

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Server Error : " + err.Error())
			}
		}
	}

	// TODO : SMIL 파일도 이때 만들자. 스케쥴 SMIL
	_, err = appLibs.WowzaHttpClient{APIUrl:appConfig.Config.WOWZA.ApiUrl}.CreateScheduleSmil(payload.TenantId);
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Wowza Server Error : " + err.Error())
	}

	tx.Create(tenant)

	// TODO : 기본 root 카테고리를 입력한다.
	category := &rdbModel.Category{
		TenantId  	: payload.TenantId,
		CategoryId	: uuid.NewV4().String(),
		Name	  	: payload.TenantName,
		Desc		: "",
		ParentIdx	: 0,
		UseYN		: "Y",
	}

	tx.Create(category)

	return handler.APIResultHandler(c, true, http.StatusCreated, tenant)
}

// 테넌트 삭제
func (api TenantAPI) DeleteTenant (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	tenant := &rdbModel.Tenant{}
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("idx = ? ",
		idx ).Find(tenant).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "Tenant NOT FOUND")
	}

	// Live Application VOD 어플리케이션 모두 삭제
	//------------------------------------------------------------------------------------------------------------------
	_, err := appLibs.WowzaHttpClient{APIUrl:appConfig.Config.WOWZA.ApiUrl}.DeleteApplication(tenant.TenantId, "Live")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Wowza Server Error : " + err.Error())
	}

	// TODO : VOD 파일 삭제는 일단 보류
	//-------------------------------------------------------------------------------------------------------------------

	// 카테고리 삭제
	category := &rdbModel.Category{
		TenantId  	:tenant.TenantId,
	}

	tx.Delete(category, "tenant_id = ?", category.TenantId)

	tx.Delete(tenant, "idx = ?", idx)

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}


// 테넌트 수정
func (api TenantAPI) PutTenant (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	payload := &api.request
	c.Bind(payload)

	tx := c.Get("Tx").(*gorm.DB)

	tenant := &rdbModel.Tenant{}

	if tx.Model(tenant).Where("idx = ? ", idx).
		Updates(rdbModel.Tenant{Name: payload.TenantName,
					TenantDesc:payload.TenantDesc,
					TimeZone:payload.TimeZone}).RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Tenant NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")
}


// 테넌트 한건 조회
func (api TenantAPI) GetTenant (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	tenant := &rdbModel.Tenant{}

	tx := c.Get("Tx").(*gorm.DB)


	if tx.Where("idx = ? ",
		idx ).Find(tenant).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "Tenant NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, tenant)

}

// 테넌트 리스트 조회
func (api TenantAPI) GetTenants (c echo.Context) error  {

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	start_date		:= c.QueryParam("start_date")
	end_date		:= c.QueryParam("end_date")

	search_string 		:= c.QueryParam("search_string")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### start_date : ", start_date)
	logrus.Debug("##### end_date : ", end_date)
	logrus.Debug("##### search_string : ", search_string)

	tenants := &[]rdbModel.Tenant{}
	tx := c.Get("Tx").(*gorm.DB)

	if search_string != "" {
		tx = tx.Where("name LIKE ? OR  tenant_desc LIKE ? ", "%" + search_string + "%", "%" + search_string + "%")
	}

	if start_date != "" {
		tx = tx.Where("created_at BETWEEN ? AND ?", start_date, end_date)
	}

	var count = 0
	tx.Find(tenants).Count(&count)

	tx.Order("idx desc").Offset(offset).Limit(limit).Find(tenants)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{"total_count": count,
					"rows": tenants})

}