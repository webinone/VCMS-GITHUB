package api

import (
	//"VCMS/apps/handler"
	//"net/http"
	//"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	//appLibs "VCMS/apps/libs"
	rdbModel "VCMS/apps/models/rdb"
	apiModel "VCMS/apps/models/api"
	appConfig "VCMS/apps/config"
	//"github.com/jinzhu/gorm"
	//"github.com/satori/go.uuid"


	"fmt"
	"net/http"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"VCMS/apps/handler"
	"github.com/satori/go.uuid"
	"strings"
	"os"
	"path/filepath"
	"io"
	//"strconv"
	"strconv"
)

type WebinarBannerAPI struct {
	requestPost WebinarBannerMasterRequest
	//requestPut  ContentRequestPut
}


type WebinarBannerMasterRequest struct {

	WebinarSiteId		string 				`validate:"required" json:"webinar_site_id"`
	WebinarBanners		[]WebinarBannerSubRequest	`validate:"required" json:"webinar_banners"`

}

type WebinarBannerSubRequest struct {

	WebinarSiteId		string 		`validate:"required" json:"webinar_site_id"`
	BannerType		string		`validate:"required" json:"banner_type"` // 1: 경품배너 2:업체배너 3: 진행페이지 배너
	BannerTitle		string		`json:"banner_title"`
	BannerDesc		string		`json:"banner_desc"`
	LinkUrl			string		`validate:"required" json:"link_url"`
	Order                   int		`json:"order"`
	SavePath		string 		`validate:"required" json:"save_path"`
	WebPath			string 		`validate:"required" json:"web_path"`

}
//
//type ContentRequestPut struct {
//	CategoryId		string		`validate:"required" json:"category_id"`
//	ContentName 		string 		`validate:"required" json:"content_name"`
//	ThumbChange		string		`validate:"required" json:"thumb_change"`
//	ThumbType		string		`json:"thumb_type"`
//	ThumbTime		string		`json:"thumb_time"`
//
//}

type MaxValue struct {
	Value		int
}

// 웨비나 사이트 등록
func (api WebinarBannerAPI) PostWebinarBanner(c echo.Context) error {

	fmt.Println("########### PostWebinarBanner !!!")

	payload := &api.requestPost
	c.Bind(payload)

	fmt.Println(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims 	 := apiModel.GetJWTClaims(c)

	logrus.Debug("#### TenantId : ", claims.TenantId)


	tx := c.Get("Tx").(*gorm.DB)


	//var banner_order = 0
	max_value := &MaxValue{}
	tx.Raw(`
		SELECT (IFNULL(MAX(banner_order),0) + 1) AS value FROM TB_WEBINAR_BANNER
		WHERE tenant_id = ?
		AND webinar_site_id = ?
		AND banner_type = ?
	`, claims.TenantId, payload.WebinarSiteId, payload.WebinarBanners[0].BannerType).Scan(max_value)

	logrus.Debug("********** banner_order : ", max_value.Value)


	if len(payload.WebinarBanners) > 0 {

		for _, v := range payload.WebinarBanners {

			webinar_banner := &rdbModel.WebinarBanner{
				TenantId	: claims.TenantId,
				WebinarSiteId	: payload.WebinarSiteId,
				BannerTitle	: payload.WebinarBanners[0].BannerTitle,
				BannerDesc	: payload.WebinarBanners[0].BannerDesc,
				BannerId	: uuid.NewV4().String(),
				BannerType	: v.BannerType,
				BannerOrder     : max_value.Value,
				SavePath	: v.SavePath,
				WebPath		: v.WebPath,
				LinkUrl		: v.LinkUrl,
			}

			tx.Create(webinar_banner)
		}
	}

	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")

}

// 배너 업로드
func (api WebinarBannerAPI) UploadBanner(c echo.Context) error {

	logrus.Debug(">>>>> UPLOAD BANNER !!")

	banner_type := c.FormValue("banner_type")

	logrus.Debug("###### banner_type : ", banner_type)

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	banners := []rdbModel.WebinarBanner{}

	var banner rdbModel.WebinarBanner

	for i, file := range files {

		logrus.Debug("##### i : ", i)
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dstDir 		:= appConfig.Config.APP.BannerRoot
		originFileName 	:= file.Filename
		originFileExt 	:= originFileName[strings.LastIndex(originFileName, ".")+1:]

		generateName := "banner_" + uuid.NewV4().String() + "." + originFileExt

		// Destination
		dst, err := os.Create(dstDir + string(filepath.Separator) + generateName)

		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		generatedFile, _ := dst.Stat()

		logrus.Debug(generatedFile.Name())

		logrus.Debug(dstDir + string(filepath.Separator) + generateName)

		banner.BannerId	= uuid.NewV4().String()
		banner.BannerType = banner_type
		//banner.BannerOrder = i
		banner.SavePath = dstDir + string(filepath.Separator) + generateName
		banner.WebPath = "/banner/" + generateName

		logrus.Debug("################### BannerId : ", banner.BannerId)

		banners = append(banners, banner)
	}

	return handler.APIResultHandler(c, true, http.StatusOK, banners)
}

 //배너 리스트 조회
func (api WebinarBannerAPI) GetWebinarBanners (c echo.Context) error  {

	claims 	 := apiModel.GetJWTClaims(c)

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	tenant_id		:= claims.TenantId

	webinar_site_id		:= c.QueryParam("webinar_site_id");
	banner_type 		:= c.QueryParam("banner_type")

	banner_title		:= c.QueryParam("banner_title")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	logrus.Debug("##### banner_type : ", banner_type)
	logrus.Debug("##### banner_title : ", banner_title)

	webinar_banners  := []rdbModel.WebinarBanner{}

	tx := c.Get("Tx").(*gorm.DB)

	tx  = tx.Where("tenant_id = ? and webinar_site_id = ? ", tenant_id, webinar_site_id)

	if banner_type != "" {
		tx = tx.Where("banner_type = ? ",  banner_type)
	}


	if banner_title != "" {
		tx = tx.Where("banner_title LIKE ? ", "%" + banner_title + "%")
	}

	var count = 0
	tx.Find(&webinar_banners).Count(&count)

	tx.Order("banner_order asc").Offset(offset).Limit(limit).Find(&webinar_banners)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_banners,
		})

}

//// Webinar 한건 조회
//func (api WebinarSiteAPI) GetWebinarSite (c echo.Context) error  {
//
//	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)
//
//	webinar_site     := &rdbModel.WebinarSite{}
//
//	//var count = 0
//	tx := c.Get("Tx").(*gorm.DB)
//
//	if tx.Preload("PreviewContent").Preload("PreviewContent.Stream").
//		Preload("Schedule").Preload("Channel").
//		Preload("Channel.Stream").Preload("PostContent").
//		Preload("PostContent.Stream").Preload("WebinarBanners").
//		Where("idx = ? ", idx ).Find(webinar_site).RecordNotFound() {
//
//		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR NOT FOUND")
//	}
//
//	return handler.APIResultHandler(c, true, http.StatusOK, webinar_site)
//}
//
//
//// Webinar 삭제
//func (api WebinarSiteAPI) DeleteWebinarSite (c echo.Context) error  {
//
//	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)
//	// TODO : 이미 진행 된 건은 삭제 할 수 없다.
//	// 1. 배너 삭제 (물리적 삭제는 하지 말자)
//	// 2. 스케쥴 Order 삭제
//	// 3. 스케쥴 삭제
//	// 4. 웨비사 사이트 삭제
//
//	webinar_banner  := &rdbModel.WebinarBanner{}
//	schedule_order  := &rdbModel.ScheduleOrder{}
//	schedule 	:= &rdbModel.Schedule{}
//	webinar_site	:= &rdbModel.WebinarSite{}
//
//	tx := c.Get("Tx").(*gorm.DB)
//
//	if tx.Where("idx = ? ", idx ).Find(webinar_site).RecordNotFound() {
//		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_SITE NOT FOUND")
//	}
//
//
//	if tx.Delete(webinar_banner, "webinar_site_id = ?", webinar_site.WebinarSiteId).RowsAffected == 0 {
//		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
//	}
//
//	if tx.Delete(schedule_order, "schedule_id = ?", webinar_site.ScheduleId).RowsAffected == 0 {
//		//return echo.NewHTTPError(http.StatusNotFound, "CONTENT NOT FOUND")
//	}
//
//	if tx.Delete(schedule, "schedule_id = ?", webinar_site.ScheduleId).RowsAffected == 0 {
//
//	}
//
//	if tx.Delete(webinar_site, "idx = ?", idx).RowsAffected == 0 {
//
//	}
//
//	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
//}
//

// 배너 수정
func (api WebinarBannerAPI) PutWebinarBanner (c echo.Context) error  {

	banner_type := c.Param("banner_type")
	claims 	 := apiModel.GetJWTClaims(c)

	payload := &api.requestPost
	c.Bind(payload)

	if err   := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx := c.Get("Tx").(*gorm.DB)


	webinar_banner  := &rdbModel.WebinarBanner{}

	//
	// 1. 배너를 모두 삭제 한다.
	if tx.Delete(webinar_banner, "webinar_site_id = ? and banner_type = ?", payload.WebinarSiteId, banner_type).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	// 2. 배너를 모두 다시 insert 한다.
	if len(payload.WebinarBanners) > 0 {

		for i, v := range payload.WebinarBanners {

			webinar_banner := &rdbModel.WebinarBanner{
				TenantId	: claims.TenantId,
				WebinarSiteId	: payload.WebinarSiteId,
				BannerTitle	: v.BannerTitle,
				BannerDesc	: v.BannerDesc,
				BannerId	: uuid.NewV4().String(),
				BannerType	: v.BannerType,
				BannerOrder     : i+1,
				SavePath	: v.SavePath,
				WebPath		: v.WebPath,
				LinkUrl		: v.LinkUrl,
			}

			tx.Create(webinar_banner)
		}
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")

}

