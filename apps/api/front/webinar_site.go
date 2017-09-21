package front

import (
	"github.com/labstack/echo"
	//"strconv"
	"github.com/jinzhu/gorm"
	"net/http"
	rdbModel "VCMS/apps/models/rdb"

	"VCMS/apps/handler"
	"time"
	"github.com/metakeule/fmtdate"
	"fmt"
	"strconv"
	"github.com/Sirupsen/logrus"
)

type WebinarSiteAPI struct {
	//requestPost WebinarSiteRequestPost
	//requestPut  ContentRequestPut
}

type WebinarSiteIfResult struct {
	WebinarSiteId 	string        `json:"webinar_site_id"`
	Title    		string            	`json:"title"`
}

// Webinar 한건 조회
func (api WebinarSiteAPI) GetWebinarSite (c echo.Context) error  {

	//idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)
	webinar_site_id  := c.Param("webinar_site_id")
	webinar_site     := &rdbModel.WebinarSite{}

	//var count = 0
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Preload("PreviewContent").
		Preload("PreviewContent.Stream").
		Preload("PreviewContent.ThumbNails").
		Preload("Schedule").
		Preload("Schedule.ScheduleOrders").
		Preload("Schedule.ScheduleOrders.Content").
		Preload("Schedule.ScheduleOrders.Content.ThumbNails").
		Preload("Schedule.ScheduleOrders.Content.Stream").
		Preload("Channel").
		Preload("Channel.Stream").
		Preload("PostContent").
		Preload("PostContent.Stream").
		Preload("PostContent.ThumbNails").
		Preload("WebinarSiteThumbNails").
		Preload("WebinarSiteFiles").
		Preload("WebinarSiteAdmins").
		Preload("WebinarSiteBackImages").
	//Preload("WebinarBanners").
		Where("webinar_site_id = ? ", webinar_site_id ).Find(webinar_site).RecordNotFound() {

		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR NOT FOUND")
	}

	offSet, _ := time.ParseDuration("+09.00h")
	now := time.Now().UTC().Add(offSet)
	//
	//fmt.Println("Today : ", fmtdate.Format("YYYY-MM-DD hh:mm:ss", now))
	//
	nowDateTime := fmtdate.Format("YYYY-MM-DD hh:mm:ss", now)

	return handler.APIResultHandler(c, true, http.StatusOK, map[string]interface{}{
		"now_datetime" 	: nowDateTime,
		"webinar_site"	: webinar_site,
	})
}

// 웨비나 사이트 들... 함대리랑 인터페이스...이건 아무나 다 되게 한다.
func (api WebinarSiteAPI) GetWebinarIfSites (c echo.Context) error  {

	presenter_company_no 	:= c.Param("presenter_company_no")
	webinar_sites     		:= []rdbModel.WebinarSite{}

	tx := c.Get("Tx").(*gorm.DB)

	var count = 0
	tx.Find(&webinar_sites).Count(&count)

	tx.Order("idx desc").Find(&webinar_sites)

	fmt.Println(len(webinar_sites))

	webinar_if_sites  := []WebinarSiteIfResult{}

	// 관련 웨비나 관련 파라미터 관련 웨비나는 사업자 번호로 구분한다.
	if presenter_company_no != "" {
		tx = tx.Where("presenter_company_no = ? ", presenter_company_no)
	}

	for index, entry := range webinar_sites {

		fmt.Println(index)

		webinar_if_site := WebinarSiteIfResult{}
		webinar_if_site.WebinarSiteId = entry.WebinarSiteId
		webinar_if_site.Title = entry.Title

		webinar_if_sites = append(webinar_if_sites, webinar_if_site)

		//webinar_if_sites[index] = webinar_if_site
	}

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_if_sites,
		})


}

// 관련 웨비나 리스트
func (api WebinarSiteAPI) GetWebinarSitesByCompanyNo (c echo.Context) error  {

	offset, _ 			 := strconv.Atoi(c.QueryParam("offset"))
	limit, _ 			 := strconv.Atoi(c.QueryParam("limit"))
	presenter_company_no := c.QueryParam("presenter_company_no")
	webinar_site_id 	 := c.QueryParam("webinar_site_id")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### presenter_company_no : ", presenter_company_no)

	webinar_sites     := []rdbModel.WebinarSite{}

	tx := c.Get("Tx").(*gorm.DB)

	//if channel_id == "" {
	//	return echo.NewHTTPError(http.StatusBadRequest, "Channel ID NOT EXIST")
	//}
	tx = tx.Preload("WebinarSiteThumbNails").Where("webinar_site_id != ? and presenter_company_no = ? ", webinar_site_id, presenter_company_no)

	var count = 0
	tx.Find(&webinar_sites).Count(&count)

	tx.Offset(offset).Limit(limit).Find(&webinar_sites)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_sites,
		})

}


type WebinarSiteAcademyResult struct {
	WebinarSiteId 		string        	`json:"webinar_site_id"`
	Url 				string			`json:"url"`
	Title    			string        	`json:"title"`
	Desc            	string			`json:"desc"`
	Status              string			`json:"status"`
	StartDateTime		string			`json:"start_date_time"`
	PresenterName   	string			`json:"presenter_name"`
	PresenterCompany	string			`json:"presenter_company"`
	PresenterDep		string			`json:"presenter_dep"`
	PresenterEmail		string			`json:"presenter_email"`
	PresenterPhone		string			`json:"presenter_phone"`
	ThumbnailPath		string			`json:"thumb_path"`

}

// 아카데미 연동 웨비나 사이트
func (api WebinarSiteAPI) GetWebinarSites (c echo.Context) error  {

	offset, _ 			 := strconv.Atoi(c.QueryParam("offset"))
	limit, _ 			 := strconv.Atoi(c.QueryParam("limit"))

	// 행사명
	title				 := c.QueryParam("title")

	//sort_name		:= c.QueryParam("sort_name")

	// 진행 상태
	status				:= c.QueryParam("status")

	sort_name			:= "start_date_time"
	sort_dir			:= c.QueryParam("sort_dir")

	if sort_dir == "" {
		sort_dir  = "DESC"
	}

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### title : ", title)
	logrus.Debug("##### sort_name : ", sort_name)
	logrus.Debug("##### sort_dir : ", sort_dir)


	webinar_sites     := []rdbModel.WebinarSite{}

	tx := c.Get("Tx").(*gorm.DB)

	if title != "" {
		tx = tx.Where("title LIKE ? ", "%" + title + "%")
	}

	offSet, _ := time.ParseDuration("+09.00h")
	now := time.Now().UTC().Add(offSet)
	fmt.Println("Today : ", fmtdate.Format("YYYY-MM-DD hh:mm:ss", now))
	today := fmtdate.Format("YYYY-MM-DD HH:mm:ss", now)

	if status != "" {
		if status == "1" {
			// 미진행
			tx = tx.Where("start_date_time > ? ",  today)
		} else if (status == "2") {
			// 진행중
			tx = tx.Where("start_date_time <= ? and end_date_time >= ? ",  today, today)
		} else if (status == "3") {
			// 종료
			tx = tx.Where("end_date_time < ? ",  today)
		}
	}

	var count = 0
	tx.Find(&webinar_sites).Count(&count)

	tx.Preload("WebinarSiteThumbNails").
		Order(sort_name +" "+sort_dir).Offset(offset).Limit(limit).Find(&webinar_sites)


	webinar_academy_sites  := []WebinarSiteAcademyResult{}

	for index, entry := range webinar_sites {

		fmt.Println(index)

		webinar_academy_site := WebinarSiteAcademyResult{}
		webinar_academy_site.WebinarSiteId 	= entry.WebinarSiteId
		webinar_academy_site.Url 			= "http://webinar.hellot.net/site/" + entry.WebinarSiteId
		webinar_academy_site.Title 			= entry.Title
		webinar_academy_site.Desc 			= entry.SubTitle
		webinar_academy_site.StartDateTime 	= entry.StartDateTime
		webinar_academy_site.PresenterCompany = entry.PresenterCompany
		webinar_academy_site.PresenterDep = entry.PresenterDep
		webinar_academy_site.PresenterName = entry.PresenterName
		webinar_academy_site.PresenterEmail = entry.PresenterEmail
		webinar_academy_site.PresenterPhone = entry.PresenterPhone

		if entry.StartDateTime > today {
			webinar_academy_site.Status = "1"
		}

		if entry.StartDateTime <= today && entry.EnterDateTime >= today {
			webinar_academy_site.Status = "2"
		}

		if entry.EnterDateTime < today {
			webinar_academy_site.Status = "3"
		}

		if len(entry.WebinarSiteThumbNails) > 0 {
			webinar_academy_site.ThumbnailPath = entry.WebinarSiteThumbNails[0].WebPath
		}

		webinar_academy_sites = append(webinar_academy_sites, webinar_academy_site)
	}

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_academy_sites,
		})


}
