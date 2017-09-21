package front

import (
	//"VCMS/apps/handler"
	//"net/http"
	"github.com/labstack/echo"
	//appLibs "VCMS/apps/libs"
	rdbModel "VCMS/apps/models/rdb"
	apiModel "VCMS/apps/models/api"
	//appConfig "VCMS/apps/config"
	"net/http"
	"github.com/Sirupsen/logrus"
	"VCMS/apps/handler"
	"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
	"strconv"
)


type WebinarCommentAPI struct {
	requestPost WebinarCommentPostRequest
}

type WebinarCommentPostRequest struct {
	WebinarSiteId		string 		`validate:"required" json:"webinar_site_id"`
	Comment			string		`validate:"required" json:"comment"`
}

// 웨비나 프론트 참여 JOIN
func (api WebinarCommentAPI) PostWebinarComment(c echo.Context) error {

	logrus.Debug("########### PostWebinarComment !!!")

	claims := apiModel.GetFrontJWTClaims(c)

	payload := &api.requestPost
	c.Bind(payload)

	logrus.Debug(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	webinar_comment     := &rdbModel.WebinarComment{}
	webinar_site     := &rdbModel.WebinarSite{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("webinar_site_id = ? ", payload.WebinarSiteId ).Find(webinar_site).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_SITE NOT FOUND")
	}

	// webinar_site_id로 tenant_id 구하기
	tenant_id := webinar_site.TenantId

	webinar_comment.TenantId = tenant_id
	webinar_comment.FrontUserId = claims.MemberId
	webinar_comment.FrontUserName = claims.MemberName
	webinar_comment.WebinarSiteId = payload.WebinarSiteId
	webinar_comment.WebinarCommentId = uuid.NewV4().String()
	webinar_comment.Comment = payload.Comment

	tx.Create(webinar_comment)

	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")
}


func (api WebinarCommentAPI) PutWebinarComment(c echo.Context) error {

	logrus.Debug("########### PutWebinarComment !!!")

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	//claims := apiModel.GetFrontJWTClaims(c)

	payload := &api.requestPost
	c.Bind(payload)

	logrus.Debug(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	webinar_comment     := &rdbModel.WebinarComment{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("idx = ? ", idx ).Find(webinar_comment).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR COMMENT NOT FOUND")
	}


	webinar_comment.Comment = payload.Comment

	tx.Save(webinar_comment)

	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")

}


//댓글 리스트 조회
func (api WebinarCommentAPI) GetWebinarComments (c echo.Context) error  {

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	webinar_site_id	:= c.QueryParam("webinar_site_id");

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)


	webinar_comments  := []rdbModel.WebinarComment{}

	tx := c.Get("Tx").(*gorm.DB)

	tx  = tx.Where("webinar_site_id = ? ", webinar_site_id)

	var count = 0
	tx.Find(&webinar_comments).Count(&count)

	tx.Order("idx desc ").Offset(offset).Limit(limit).Find(&webinar_comments)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_comments,
		})

}

// 댓글 삭제
func (api WebinarCommentAPI) DeleteWebinarComment(c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	tx := c.Get("Tx").(*gorm.DB)

	webinar_comment  := &rdbModel.WebinarComment{}

	if tx.Where("idx = ? ", idx ).Find(webinar_comment).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR COMMENT NOT FOUND")
	}


	if tx.Delete(webinar_comment, "webinar_site_id = ? and idx = ? ", webinar_comment.WebinarSiteId, idx).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "CONTENT NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}

