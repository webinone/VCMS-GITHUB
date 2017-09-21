package front

import (
	"github.com/labstack/echo"
	rdbModel "VCMS/apps/models/rdb"
	"net/http"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"VCMS/apps/handler"
	"strconv"
)

type WebinarBannerAPI struct {
	//requestPost WebinarBannerMasterRequest
	//requestPut  ContentRequestPut
}


 //배너 리스트 조회
func (api WebinarBannerAPI) GetWebinarBanners (c echo.Context) error  {

	//claims 	 := apiModel.GetJWTClaims(c)

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	webinar_site_id		:= c.QueryParam("webinar_site_id")
	banner_type 		:= c.QueryParam("banner_type")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	logrus.Debug("##### banner_type : ", banner_type)

	webinar_banners  := []rdbModel.WebinarBanner{}

	tx := c.Get("Tx").(*gorm.DB)

	tx  = tx.Where("webinar_site_id = ? ", webinar_site_id)

	if banner_type != "" {
		tx = tx.Where("banner_type = ? ", banner_type)
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


// 배너 사이트 이동 Redirect !!!
func (api WebinarBannerAPI) GetWebinarBannerRedirect (c echo.Context) error  {

	//claims 			:= apiModel.GetFrontJWTClaims(c)
	idx, _		:= strconv.Atoi(c.QueryParam("idx"))

	//member_id		:= claims.MemberId
	webinar_banner  := rdbModel.WebinarBanner{}

	tx := c.Get("Tx").(*gorm.DB)
	tx.Where("idx = ? ", idx).Find(&webinar_banner)


	return c.Redirect(http.StatusMovedPermanently, webinar_banner.LinkUrl)
}


