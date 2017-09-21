package front

import (
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"net/http"
	"github.com/labstack/echo"
	"strconv"
	rdbModel "VCMS/apps/models/rdb"
	"VCMS/apps/handler"
)

type WebinarNoticeAPI struct {
	//requestPost WebinarNoticeRequest
	//requestPut  ContentRequestPut
}


// 공지사항 리스트 조회
func (api WebinarNoticeAPI) GetWebinarNotices (c echo.Context) error  {

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	webinar_site_id		:= c.QueryParam("webinar_site_id");

	//title			:= c.QueryParam("title")
	use_yn			:= c.QueryParam("use_yn")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	//logrus.Debug("##### start_date : ", start_date)
	//logrus.Debug("##### end_date : ", end_date)
	//logrus.Debug("##### title : ", title)
	logrus.Debug("##### use_yn : ", use_yn)

	tx := c.Get("Tx").(*gorm.DB)

	webinar_notices    := []rdbModel.WebinarNotice{}

	tx  = tx.Where("webinar_site_id = ? ", webinar_site_id)

	////if start_date != "" {
	////	tx = tx.Where("created_at BETWEEN ? AND ?", start_date, end_date)
	////}
	//
	//if title != "" {
	//	tx = tx.Where("title LIKE ? ", "%" + title + "%")
	//}

	if use_yn != "" {
		tx = tx.Where("use_yn = ? ", use_yn)
	}

	var count = 0
	tx.Find(&webinar_notices).Count(&count)

	tx.Preload("WebinarNoticeFiles").Order("idx desc").Offset(offset).Limit(limit).Find(&webinar_notices)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_notices,
		})

}