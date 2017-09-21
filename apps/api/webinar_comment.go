package api

import (
	"github.com/labstack/echo"
	//appLibs "VCMS/apps/libs"
	"VCMS/apps/handler"
	apiModel "VCMS/apps/models/api"
	rdbModel "VCMS/apps/models/rdb"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"

	//"strconv"
	"strconv"
)

type WebinarCommentAPI struct {
	//requestPost WebinarBannerMasterRequest
	//requestPut  ContentRequestPut
}

//댓글 리스트 조회
func (api WebinarCommentAPI) GetWebinarComments(c echo.Context) error {

	claims := apiModel.GetJWTClaims(c)

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	tenant_id := claims.TenantId

	webinar_site_id := c.QueryParam("webinar_site_id")
	member_info := c.QueryParam("member_info")
	created_at := c.QueryParam("created_at")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	logrus.Debug("##### created_at : ", created_at)

	logrus.Debug("##### member_info : ", member_info)

	webinar_comments := []rdbModel.WebinarComment{}

	tx := c.Get("Tx").(*gorm.DB)

	tx = tx.Where("tenant_id = ? and webinar_site_id = ? ", tenant_id, webinar_site_id)

	if member_info != "" {
		tx = tx.Where("front_user_id LIKE ? OR  front_user_name LIKE ? ", "%"+member_info+"%", "%"+member_info+"%")
	}

	if created_at != "" {
		tx = tx.Where("created_at LIKE ? ", created_at+"%")
	}

	var count = 0
	tx.Find(&webinar_comments).Count(&count)

	tx.Order("idx desc ").Offset(offset).Limit(limit).Find(&webinar_comments)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows":        webinar_comments,
		})

}

// 댓글 삭제
func (api WebinarCommentAPI) DeleteWebinarComment(c echo.Context) error {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	tx := c.Get("Tx").(*gorm.DB)

	webinar_comment := &rdbModel.WebinarComment{}

	if tx.Where("idx = ? ", idx).Find(webinar_comment).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR COMMENT NOT FOUND")
	}

	if tx.Delete(webinar_comment, "webinar_site_id = ? and idx = ? ", webinar_comment.WebinarSiteId, idx).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "CONTENT NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}
