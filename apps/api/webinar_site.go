package api

import (
	"VCMS/apps/handler"
	"net/http"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	appLibs "VCMS/apps/libs"
	rdbModel "VCMS/apps/models/rdb"
	apiModel "VCMS/apps/models/api"
	appConfig "VCMS/apps/config"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"strconv"
	"os"
	"io"
	"strings"
	"path/filepath"
	"fmt"
	"time"
	"github.com/metakeule/fmtdate"
)

type WebinarSiteAPI struct {
	requestPost WebinarSiteRequestPost
	//requestPut  ContentRequestPut
}


type WebinarSiteRequestPost struct {
	Title 				string 					`validate:"required" json:"title"`
	SubTitle 			string					`json:"sub_title"`
	Tags    			[]WebinarSiteTagRequest `json:"tags"`
	StartDateTime		string					`validate:"required" json:"start_datetime"`
	EnterBeforeTime		string					`validate:"required" json:"enter_before_time"`
	EnterDateTime		string					`validate:"required" json:"enter_datetime"`

	PresenterName		string					`validate:"required" json:"presenter_name"`
	PresenterCompany	string					`json:"presenter_company"`
	PresenterCompanyNo	string					`validate:"required" json:"presenter_company_no"`
	PresenterDep		string					`json:"presenter_dep"`
	PresenterPosition	string					`json:"presenter_position"`
	PresenterEmail		string					`json:"presenter_email"`
	PresenterPhone		string					`json:"presenter_phone"`
	PresenterFax		string					`json:"presenter_fax"`

	PreviewType			string 						`validate:"required" json:"preview_type"`  // 1: 내부 , 2: 외부
	PreviewContentId	string					`json:"preview_content_id"`
	PreviewExtraUrl		string					`json:"preview_extra_url"`

	ChannelId			string						`json:"channel_id"`
	ScheduleOrders    	[]ScheduleOrderRequest 	`validate:"required" json:"schedule_orders"`

	PostType			string 						`validate:"required" json:"post_type"`  // 1: 내부 , 2: 외부
	PostContentId		string					`json:"post_content_id"`
	PostExtraUrl		string					`json:"post_extra_url"`

	ThumbNails			[]WebinarSiteFilesRequest 	`validate:"required" json:"thumbnails"`
	SiteFiles			[]WebinarSiteFilesRequest 	`json:"site_files"`
	BackImages			[]WebinarSiteFilesRequest 	`json:"back_images"`
	AdminUsers			[]WebinarSiteAdminUserRequest	`json:"admin_users"`

	YoutubeLiveUrl		string						`json:"youtube_live_url"`

	HostName			string						`validate:"required" json:"host_name"`
	HostEmail			string						`json:"host_email"`
	HostPhone			string						`json:"host_phone"`
	HostFax				string						`json:"host_fax"`

}

type WebinarSiteTagRequest struct {
	Idx				int 		`json:"idx"`
	Name			string 		`json:"name"`
}

type WebinarSiteFilesRequest struct {
	OriginalName		string		`validate:"required" json:"original_name"`
	SavePath		string 			`validate:"required" json:"save_path"`
	WebPath			string 			`validate:"required" json:"web_path"`

}

type WebinarSiteAdminUserRequest struct {
	MemberId		string			`validate:"required" json:"member_id"`
	MemberName		string			`validate:"required" json:"member_name"`
}

// 웨비나 사이트 등록
func (api WebinarSiteAPI) PostWebinarSite(c echo.Context) error {

	fmt.Println("########### PostWebinarSite !!!")

	payload := &api.requestPost
	c.Bind(payload)

	fmt.Println(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims 	 := apiModel.GetJWTClaims(c)

	// TODO : 해당 채널에 동시간대 스케쥴이 존재한다면 등록을 막는다.


	logrus.Debug("#### TenantId : ", claims.TenantId)

	webinar_site_id := uuid.NewV4().String()
	schedule_id 	:= uuid.NewV4().String()

	logrus.Debug("############ schedule_id : ", schedule_id)

	schedule_orders := []rdbModel.ScheduleOrder{}
	tags 		:= []rdbModel.WebinarSiteTag{}

	thumbnails 	:= []rdbModel.WebinarSiteThumbNail{}
	site_files 	:= []rdbModel.WebinarSiteFile{}
	backimages 	:= []rdbModel.WebinarSiteBackImage{}
	admin_users	:= []rdbModel.WebinarSiteAdmin{}

	total_time	:= 0

	videos 		:= []appLibs.Video{}

	if len(payload.ScheduleOrders) > 0 {

		for i, v := range payload.ScheduleOrders {

			s := append (schedule_orders, rdbModel.ScheduleOrder{

				TenantID	: claims.TenantId,
				ScheduleId 	: schedule_id,
				ScheduleOrderId : uuid.NewV4().String(),
				Order		: i+1 ,
				ContentId 	: v.ContentId,
				StartSec	: v.StartSec,
				EndSec		: v.EndSec,
				UpdatedId       : claims.UserId,
			})
			schedule_orders = s

			video := append (videos, appLibs.Video{
				Src    : "mp4:assets/"+v.GeneratedFileName,
				Start  : v.StartSec,
				Length : v.EndSec,
			})
			videos = video

			number_end_sec, _ 	:= strconv.Atoi(v.EndSec)
			number_start_sec, _ 	:= strconv.Atoi(v.StartSec)

			total_time += number_end_sec - number_start_sec
		}

	}

	// 태그가 존재한다면..
	if len(payload.Tags) > 0 {
		for _, v := range payload.Tags {

			s := append (tags, rdbModel.WebinarSiteTag{

				TenantId: claims.TenantId,
				WebinarSiteId	: webinar_site_id,
				Name 			: v.Name,
				UpdatedId       : claims.UserId,
			})
			tags = s

		}
	}

	// 썸네일이 존재한다면..
	if len(payload.ThumbNails) > 0 {

		for _, v := range payload.ThumbNails {

			s := append (thumbnails, rdbModel.WebinarSiteThumbNail {

				TenantId: claims.TenantId,
				WebinarSiteId: webinar_site_id,
				WebinarSiteThumbId 	: uuid.NewV4().String(),
				OriginalName		: v.OriginalName,
				SavePath		: v.SavePath,
				WebPath			: v.WebPath,
				UpdatedId       : claims.UserId,
			})
			thumbnails = s
		}
	}

	// 세미나 자료가 존재한다면..
	if len(payload.SiteFiles) > 0 {

		for _, v := range payload.SiteFiles {

			s := append (site_files, rdbModel.WebinarSiteFile {

				TenantId: claims.TenantId,
				WebinarSiteId: webinar_site_id,
				WebinarSiteFileId 	: uuid.NewV4().String(),
				OriginalName		: v.OriginalName,
				SavePath		: v.SavePath,
				WebPath			: v.WebPath,
				UpdatedId       : claims.UserId,
			})
			site_files = s
		}
	}

	// TOP 이미지가 존재한다면...
	if len(payload.BackImages) > 0 {

		for _, v := range payload.BackImages {

			s := append (backimages, rdbModel.WebinarSiteBackImage {

				TenantId: claims.TenantId,
				WebinarSiteId			: webinar_site_id,
				WebinarSiteBackImageId 	: uuid.NewV4().String(),
				OriginalName			: v.OriginalName,
				SavePath				: v.SavePath,
				WebPath				: v.WebPath,
			})
			backimages = s
		}
	}

	// 관리자 유저들이 존재한다면...
	if len(payload.AdminUsers) > 0 {
		for _, v := range payload.AdminUsers {
			s := append (admin_users, rdbModel.WebinarSiteAdmin {
				TenantId: claims.TenantId,
				WebinarSiteId			: webinar_site_id,
				MemberId: v.MemberId,
				MemberName:v.MemberName,
			})
			admin_users = s
		}
	}

	offSet, _ := time.ParseDuration("+09.00h")
	start_date_time, _ := fmtdate.Parse("YYYY-MM-DD hh:mm:ss", payload.StartDateTime)
	start_date_time.UTC().Add(offSet)

	end_date_time := start_date_time.Add(time.Duration(total_time) * time.Second)

	fmt.Println("end_date_time : ", fmtdate.Format("YYYY-MM-DD hh:mm:ss", end_date_time))

	webinar_site := &rdbModel.WebinarSite{
		TenantId			: claims.TenantId,
		WebinarSiteId		: webinar_site_id,
		Title				: payload.Title,
		SubTitle			: payload.SubTitle,
		WebinarSiteTags 	: tags,
		StartDateTime		: payload.StartDateTime,
		TotalTime			: strconv.Itoa(total_time),
		EndDateTime			: fmtdate.Format("YYYY-MM-DD hh:mm:ss", end_date_time),
		EnterBeforeTime 	: payload.EnterBeforeTime,
		EnterDateTime 		: payload.EnterDateTime,
		PresenterName		: payload.PresenterName,
		PresenterCompany	: payload.PresenterCompany,
		PresenterCompanyNo	: payload.PresenterCompanyNo,
		PresenterDep		: payload.PresenterDep,
		PresenterPosition	: payload.PresenterPosition,
		PresenterEmail		: payload.PresenterEmail,
		PresenterPhone		: payload.PresenterPhone,
		PresenterFax		: payload.PresenterFax,

		PreviewType		: payload.PreviewType,
		PreviewContentId	: payload.PreviewContentId,
		PreviewExtraUrl		: payload.PreviewExtraUrl,

		ChannelId		: payload.ChannelId,

		ScheduleId		: schedule_id,
		Schedule		: rdbModel.Schedule {
			TenantID	:	claims.TenantId,
			ChannelId	:	payload.ChannelId,
			ScheduleId	:	schedule_id,
			Name	  	:	payload.Title,
			StartDateTime 	:	payload.StartDateTime,
			TotalTime	:	strconv.Itoa(total_time),
			ScheduleOrders  :	schedule_orders,
			UpdatedId       : 	claims.UserId,
		},

		WebinarSiteThumbNails	: thumbnails,
		WebinarSiteFiles	:	site_files,
		WebinarSiteBackImages:  backimages,
		WebinarSiteAdmins: admin_users,

		PostType		: payload.PostType,
		PostContentId		: payload.PostContentId,
		PostExtraUrl		: payload.PostExtraUrl,

		YoutubeLiveUrl: payload.YoutubeLiveUrl,

		HostName		: payload.HostName,
		HostEmail		: payload.HostEmail,
		HostPhone		: payload.HostPhone,
		HostFax			: payload.HostFax,

	}

	tx := c.Get("Tx").(*gorm.DB)

	if !tx.Where("webinar_site_id = ?",
		webinar_site_id).Find(webinar_site).RecordNotFound() {
		return echo.NewHTTPError(http.StatusConflict, "Already exists WebinarSite ")
	}

	// TODO : 스케쥴 생성시에 WOWZA 수정한다 (운영 적용후 테스트)
	//-------------------------------------------------------------------------------
	smilFilePath := appConfig.Config.APP.ContentsRoot + string(filepath.Separator) + claims.TenantId + string(filepath.Separator) + "streamschedule.smil"

	err := appLibs.ScheduleSmilParser{ FilePath: smilFilePath}.
		UpdateScheduleSmil( payload.ChannelId, "false",  schedule_id, payload.EnterDateTime, videos)

	if err != nil {
		fmt.Println(">>>>>>>>>>>>>>> UpdateScheduleSmil ")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//-------------------------------------------------------------------------------

	// 수정 후에는 반드시 Reload를 호출해야 한다.
	//--------------------------------------------------------------------------------
	_, err = appLibs.WowzaHttpClient{APIUrl:appConfig.Config.WOWZA.CustomApiUrl}.ReloadScheduleSmil(claims.TenantId);
	if err != nil {
		fmt.Println(">>>>>>>>>>>>>>> ReloadScheduleSmil ")
		return echo.NewHTTPError(http.StatusInternalServerError, "Wowza Server Error : " + err.Error())
	}
	//-------------------------------------------------------------------------------


	tx.Create(webinar_site)

	return handler.APIResultHandler(c, true, http.StatusCreated, webinar_site)

}

// 배너 업로드
func (api WebinarSiteAPI) UploadBanner(c echo.Context) error {

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

		banner.BannerType = banner_type
		banner.BannerOrder = i
		banner.SavePath = dstDir + string(filepath.Separator) + generateName
		banner.WebPath = "/banner/" + generateName

		banners = append(banners, banner)
	}

	return handler.APIResultHandler(c, true, http.StatusOK, banners)
}

// 자료 업로드
func (api WebinarSiteAPI) UploadThumbNail(c echo.Context) error {

	logrus.Debug(">>>>> UploadThumbNail !!")

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	thumbnails := []rdbModel.WebinarSiteThumbNail{}

	var thumbnail rdbModel.WebinarSiteThumbNail

	for i, file := range files {

		logrus.Debug("##### i : ", i)
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dstDir 			:= appConfig.Config.APP.ThumbNailRoot
		originFileName 	:= file.Filename
		originFileExt 	:= originFileName[strings.LastIndex(originFileName, ".")+1:]

		generateName 	:= "webinar_site_thmb_" + uuid.NewV4().String() + "." + originFileExt

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

		//banner.BannerOrder = i
		thumbnail.OriginalName = originFileName
		thumbnail.FileSize = generatedFile.Size()
		thumbnail.SavePath = dstDir + string(filepath.Separator) + generateName
		thumbnail.WebPath = "/thumbnail/" + generateName

		thumbnails = append(thumbnails, thumbnail)
	}
	return handler.APIResultHandler(c, true, http.StatusOK, thumbnails)
}

// 자료 업로드
func (api WebinarSiteAPI) UploadSiteFile(c echo.Context) error {

	logrus.Debug(">>>>> UploadSiteFile !!")

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	site_files := []rdbModel.WebinarSiteFile{}

	var site_file rdbModel.WebinarSiteFile

	for i, file := range files {

		logrus.Debug("##### i : ", i)
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dstDir 			:= appConfig.Config.APP.DownloadRoot
		originFileName 	:= file.Filename
		originFileExt 	:= originFileName[strings.LastIndex(originFileName, ".")+1:]

		generateName 	:= "webinar_site_" + uuid.NewV4().String() + "." + originFileExt

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

		//banner.BannerOrder = i
		site_file.OriginalName = originFileName
		site_file.FileSize = generatedFile.Size()
		site_file.SavePath = dstDir + string(filepath.Separator) + generateName
		site_file.WebPath = "/download/" + generateName

		site_files = append(site_files, site_file)
	}
	return handler.APIResultHandler(c, true, http.StatusOK, site_files)
}

// 메인 TOP 이미지
func (api WebinarSiteAPI) UploadBackImage(c echo.Context) error {

	logrus.Debug(">>>>> UploadBackImage !!")

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	backimages := []rdbModel.WebinarSiteBackImage{}

	var backimage rdbModel.WebinarSiteBackImage

	for i, file := range files {

		logrus.Debug("##### i : ", i)
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dstDir 			:= appConfig.Config.APP.ThumbNailRoot
		originFileName 	:= file.Filename
		originFileExt 	:= originFileName[strings.LastIndex(originFileName, ".")+1:]

		generateName 	:= "webinar_site_backimage_" + uuid.NewV4().String() + "." + originFileExt

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

		//banner.BannerOrder = i
		backimage.OriginalName = originFileName
		backimage.FileSize = generatedFile.Size()
		backimage.SavePath = dstDir + string(filepath.Separator) + generateName
		backimage.WebPath = "/thumbnail/" + generateName

		backimages = append(backimages, backimage)
	}
	return handler.APIResultHandler(c, true, http.StatusOK, backimages)
}

// 컨텐츠 리스트 조회
func (api WebinarSiteAPI) GetWebinarSites (c echo.Context) error  {

	claims 	 := apiModel.GetJWTClaims(c)

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	start_date		:= c.QueryParam("start_date")
	end_date		:= c.QueryParam("end_date")

	start_datetime		:= c.QueryParam("start_datetime")

	tenant_id		:= claims.TenantId

	channel_id		:= c.QueryParam("channel_id");
	title 			:= c.QueryParam("title")

	status          := c.QueryParam("status")

	sort_name		:= c.QueryParam("sort_name")
	sort_dir		:= c.QueryParam("sort_dir")

	if sort_name == "" {
		sort_name = "idx"
		sort_dir  = "DESC"
	}

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### start_date : ", start_date)
	logrus.Debug("##### end_date : ", end_date)
	logrus.Debug("##### tenant_id : ", tenant_id)
	logrus.Debug("##### channel_id : ", channel_id)
	logrus.Debug("##### title : ", title)

	webinar_sites     := []rdbModel.WebinarSite{}

	tx := c.Get("Tx").(*gorm.DB)

	//if channel_id == "" {
	//	return echo.NewHTTPError(http.StatusBadRequest, "Channel ID NOT EXIST")
	//}

	tx = tx.Where("tenant_id = ? ", claims.TenantId)


	if channel_id != "" {
		tx = tx.Where("channel_id = ? ", channel_id)
	}

	if title != "" {
		tx = tx.Where("title LIKE ? ", "%" + title + "%")
	}

	if start_date != "" {
		tx = tx.Where("created_at BETWEEN ? AND ?", start_date, end_date)
	}

	if start_datetime != "" {
		tx = tx.Where("start_date_time LIKE ?", start_datetime + "%")
	}

	if status != "" {

		offSet, _ := time.ParseDuration("+09.00h")
		now := time.Now().UTC().Add(offSet)
		fmt.Println("Today : ", fmtdate.Format("YYYY-MM-DD hh:mm:ss", now))
		today := fmtdate.Format("YYYY-MM-DD HH:mm:ss", now)

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

	tx.Preload("PreviewContent").
		Preload("PreviewContent.Stream").
		Preload("Schedule").
		Preload("Schedule.ScheduleOrders").
		Preload("Channel").
		Preload("Channel.Stream").
		Preload("PostContent").
		Preload("PostContent.Stream").
		Preload("WebinarBanners").Order(sort_name +" "+sort_dir).Offset(offset).Limit(limit).Find(&webinar_sites)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{
			"total_count": count,
			"rows": webinar_sites,
		})

}



// Webinar 한건 조회
func (api WebinarSiteAPI) GetWebinarSite (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	webinar_site     := &rdbModel.WebinarSite{}

	//var count = 0
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Preload("WebinarSiteTags").
		Preload("PreviewContent").
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
		Where("idx = ? ", idx ).Find(webinar_site).RecordNotFound() {

		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, webinar_site)
}


// Webinar 삭제
func (api WebinarSiteAPI) DeleteWebinarSite (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)
	// TODO : 이미 진행 된 건은 삭제 할 수 없다.
	// 1. 배너 삭제 (물리적 삭제는 하지 말자)
	// 2. 스케쥴 Order 삭제
	// 3. 스케쥴 삭제
	// 4. 웨비사 사이트 삭제

	webinar_banner  := &rdbModel.WebinarBanner{}
	schedule_order  := &rdbModel.ScheduleOrder{}
	schedule 	:= &rdbModel.Schedule{}
	webinar_site	:= &rdbModel.WebinarSite{}
	thumbnail       := &rdbModel.WebinarSiteThumbNail{}
	site_file		:= &rdbModel.WebinarSiteFile{}

	backimage		:= &rdbModel.WebinarSiteBackImage{}
	admin_user		:= &rdbModel.WebinarSiteAdmin{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("idx = ? ", idx ).Find(webinar_site).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_SITE NOT FOUND")
	}


	if tx.Delete(webinar_banner, "webinar_site_id = ?", webinar_site.WebinarSiteId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	if tx.Delete(schedule_order, "schedule_id = ?", webinar_site.ScheduleId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "CONTENT NOT FOUND")
	}

	if tx.Delete(schedule, "schedule_id = ?", webinar_site.ScheduleId).RowsAffected == 0 {

	}

	if tx.Delete(thumbnail, "webinar_site_id = ?", webinar_site.WebinarSiteId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	if tx.Delete(site_file, "webinar_site_id = ?", webinar_site.WebinarSiteId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	if tx.Delete(backimage, "webinar_site_id = ?", webinar_site.WebinarSiteId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	if tx.Delete(admin_user, "webinar_site_id = ?", webinar_site.WebinarSiteId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	if tx.Delete(webinar_site, "idx = ?", idx).RowsAffected == 0 {

	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}


// Webinar 수정
func (api WebinarSiteAPI) PutWebinarSite (c echo.Context) error  {

	idx, _ 	 := strconv.ParseInt(c.Param("idx"), 0, 64)
	claims 	 := apiModel.GetJWTClaims(c)

	payload := &api.requestPost

	fmt.Println("payload !!")
	fmt.Println(payload)

	c.Bind(payload)

	if err   := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx := c.Get("Tx").(*gorm.DB)

	webinar_site := &rdbModel.WebinarSite{}

	if tx.Where("idx = ? ", idx ).Find(webinar_site).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR SITE NOT FOUND")
	}

	schedule_id 	:= webinar_site.ScheduleId

	schedule_order  := &rdbModel.ScheduleOrder{}
	schedule 		:= &rdbModel.Schedule{}
	tag 			:= &rdbModel.WebinarSiteTag{}
	thumbnail       := &rdbModel.WebinarSiteThumbNail{}
	site_file		:= &rdbModel.WebinarSiteFile{}

	backimage		:= &rdbModel.WebinarSiteBackImage{}
	admin_user		:= &rdbModel.WebinarSiteAdmin{}

	// TODO : 채널이 변경되었으면 채널이 유효한지 체크한다. 해당 날짜에 가능한 것인지?????
	// 2. 스케쥴 order도 모두 삭제한다.
	if tx.Delete(schedule_order, "schedule_id = ?", schedule_id).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "CONTENT NOT FOUND")
	}

	// 3. 스케쥴도 삭제한다.
	if tx.Delete(schedule, "schedule_id = ?", webinar_site.ScheduleId).RowsAffected == 0 {

	}

	// 4. 태그 삭제
	if tx.Delete(tag, "webinar_site_id = ?", webinar_site.WebinarSiteId).RowsAffected == 0 {

	}

	if len(payload.ThumbNails) > 0 {
		// 5. 썸네일 삭제
		if tx.Delete(thumbnail, "webinar_site_id = ?", webinar_site.WebinarSiteId).RowsAffected == 0 {

		}
	}

	if len(payload.SiteFiles) > 0 {
		// 6. 사이트 파일 삭제
		if tx.Delete(site_file, "webinar_site_id = ?", webinar_site.WebinarSiteId).RowsAffected == 0 {

		}
	}

	// backimage가 존재한다면..
	if len(payload.BackImages) > 0 {
		if tx.Delete(backimage, "webinar_site_id = ?", webinar_site.WebinarSiteId).RowsAffected == 0 {

		}
	}

	// 관리자 유저 모두 삭제
	if tx.Delete(admin_user, "webinar_site_id = ?", webinar_site.WebinarSiteId).RowsAffected == 0 {

	}

	schedule_orders := []rdbModel.ScheduleOrder{}
	tags := []rdbModel.WebinarSiteTag{}
	total_time	:= 0

	videos 		:= []appLibs.Video{}

	// 그리고 다시 순서대로 넣는다.
	// 스케쥴을 수정한다.

	if len(payload.ScheduleOrders) > 0 {

		for i, v := range payload.ScheduleOrders {

			s := append (schedule_orders, rdbModel.ScheduleOrder{

				TenantID	: claims.TenantId,
				ScheduleId 	: schedule_id,
				ScheduleOrderId : uuid.NewV4().String(),
				Order		: i+1 ,
				ContentId 	: v.ContentId,
				StartSec	: v.StartSec,
				EndSec		: v.EndSec,
				UpdatedId       : claims.UserId,
			})
			schedule_orders = s

			video := append (videos, appLibs.Video{
				Src    : "mp4:assets/"+v.GeneratedFileName,
				Start  : v.StartSec,
				Length : v.EndSec,
			})
			videos = video

			number_end_sec, _ 	:= strconv.Atoi(v.EndSec)
			number_start_sec, _ 	:= strconv.Atoi(v.StartSec)

			total_time += number_end_sec - number_start_sec
		}

	}

	// 태그가 존재한다면..
	if len(payload.Tags) > 0 {
		for _, v := range payload.Tags {

			s := append (tags, rdbModel.WebinarSiteTag{

				TenantId: claims.TenantId,
				WebinarSiteId	: webinar_site.WebinarSiteId,
				Name 			: v.Name,
				UpdatedId       : claims.UserId,
			})
			tags = s

		}
	}

	thumbnails 	:= []rdbModel.WebinarSiteThumbNail{}
	site_files 	:= []rdbModel.WebinarSiteFile{}


	// 썸네일이 존재한다면..
	if len(payload.ThumbNails) > 0 {

		for _, v := range payload.ThumbNails {

			s := append (thumbnails, rdbModel.WebinarSiteThumbNail {

				TenantId: claims.TenantId,
				WebinarSiteId: webinar_site.WebinarSiteId,
				WebinarSiteThumbId 	: uuid.NewV4().String(),
				OriginalName		: v.OriginalName,
				SavePath		: v.SavePath,
				WebPath			: v.WebPath,
				UpdatedId       : claims.UserId,
			})
			thumbnails = s
		}
	}

	// 세미나 자료가 존재한다면..
	if len(payload.SiteFiles) > 0 {

		for _, v := range payload.SiteFiles {

			s := append (site_files, rdbModel.WebinarSiteFile {

				TenantId: claims.TenantId,
				WebinarSiteId: webinar_site.WebinarSiteId,
				WebinarSiteFileId 	: uuid.NewV4().String(),
				OriginalName		: v.OriginalName,
				SavePath		: v.SavePath,
				WebPath			: v.WebPath,
				UpdatedId       : claims.UserId,
			})
			site_files = s
		}
	}

	backimages 	:= []rdbModel.WebinarSiteBackImage{}
	admin_users	:= []rdbModel.WebinarSiteAdmin{}

	// TOP 이미지가 존재한다면...
	if len(payload.BackImages) > 0 {

		for _, v := range payload.BackImages {

			s := append (backimages, rdbModel.WebinarSiteBackImage {

				TenantId: claims.TenantId,
				WebinarSiteId			:  webinar_site.WebinarSiteId,
				WebinarSiteBackImageId 	: uuid.NewV4().String(),
				OriginalName			: v.OriginalName,
				SavePath				: v.SavePath,
				WebPath				: v.WebPath,
			})
			backimages = s
		}
	}

	// 관리자 유저들이 존재한다면...
	if len(payload.AdminUsers) > 0 {
		for _, v := range payload.AdminUsers {
			s := append (admin_users, rdbModel.WebinarSiteAdmin {
				TenantId: claims.TenantId,
				WebinarSiteId			:  webinar_site.WebinarSiteId,
				MemberId: v.MemberId,
				MemberName:v.MemberName,
			})
			admin_users = s
		}
	}


	offSet, _ := time.ParseDuration("+09.00h")

	fmt.Println("payload.StartDateTime : ", payload.StartDateTime)

	start_date_time, _ := fmtdate.Parse("YYYY-MM-DD hh:mm:ss", payload.StartDateTime)
	start_date_time.UTC().Add(offSet)

	//fmt.Println("end_date_time : ", fmtdate.Format("YYYY-MM-DD hh:mm:ss", end_date_time))

	end_date_time := start_date_time.Add(time.Duration(total_time) * time.Second)

	fmt.Println("end_date_time : ", fmtdate.Format("YYYY-MM-DD hh:mm:ss", end_date_time))

	webinar_site.Title 				= payload.Title
	webinar_site.SubTitle 			= payload.SubTitle
	webinar_site.WebinarSiteTags 	= tags
	webinar_site.StartDateTime 		= payload.StartDateTime
	webinar_site.TotalTime 			= strconv.Itoa(total_time)
	webinar_site.EndDateTime		= fmtdate.Format("YYYY-MM-DD hh:mm:ss", end_date_time)

	webinar_site.EnterBeforeTime 	= payload.EnterBeforeTime
	webinar_site.EnterDateTime 		= payload.EnterDateTime

	webinar_site.PresenterName 		= payload.PresenterName
	webinar_site.PresenterCompany 	= payload.PresenterCompany
	webinar_site.PresenterCompanyNo = payload.PresenterCompanyNo
	webinar_site.PresenterDep 		= payload.PresenterDep
	webinar_site.PresenterPosition 	= payload.PresenterPosition
	webinar_site.PresenterEmail 	= payload.PresenterEmail
	webinar_site.PresenterPhone 	= payload.PresenterPhone
	webinar_site.PresenterFax 		= payload.PresenterFax
	webinar_site.PreviewType 		= payload.PreviewType
	webinar_site.PreviewContentId 	= payload.PreviewContentId
	webinar_site.PreviewExtraUrl 	= payload.PreviewExtraUrl

	webinar_site.ChannelId 		= payload.ChannelId

	webinar_site.ScheduleId 	= schedule_id
	webinar_site.Schedule		= rdbModel.Schedule {
		TenantID		:	claims.TenantId,
		ChannelId		:	payload.ChannelId,
		ScheduleId		:	schedule_id,
		Name	  		:	payload.Title,
		StartDateTime 	:	payload.StartDateTime,
		TotalTime		:	strconv.Itoa(total_time),
		ScheduleOrders  :	schedule_orders,
		UpdatedId       : 	claims.UserId,
	}

	webinar_site.PostType 		= payload.PostType
	webinar_site.PostContentId 	= payload.PostContentId
	webinar_site.PostExtraUrl 	= payload.PostExtraUrl

	webinar_site.HostName 		= payload.HostName
	webinar_site.HostEmail 		= payload.HostEmail
	webinar_site.HostPhone 		= payload.HostPhone
	webinar_site.HostFax 		= payload.HostFax

	webinar_site.YoutubeLiveUrl = payload.YoutubeLiveUrl

	if len(thumbnails) > 0 {
		webinar_site.WebinarSiteThumbNails = thumbnails
	}

	if len(site_files) > 0 {
		webinar_site.WebinarSiteFiles = site_files
	}

	if len(backimages) > 0 {
		webinar_site.WebinarSiteBackImages = backimages
	}

	if len(admin_users) > 0 {
		webinar_site.WebinarSiteAdmins = admin_users
	}

	/// TODO :  WOWZA 수정한다
	//-------------------------------------------------------------------------------
	smilFilePath := appConfig.Config.APP.ContentsRoot + string(filepath.Separator) + claims.TenantId + string(filepath.Separator) + "streamschedule.smil"

	err := appLibs.ScheduleSmilParser{ FilePath: smilFilePath}.
		UpdateScheduleSmil( payload.ChannelId, "false",  schedule_id, payload.EnterDateTime, videos)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//-------------------------------------------------------------------------------

	// 수정 후에는 반드시 Reload를 호출해야 한다.
	//--------------------------------------------------------------------------------
	_, err = appLibs.WowzaHttpClient{APIUrl:appConfig.Config.WOWZA.CustomApiUrl}.ReloadScheduleSmil(claims.TenantId);
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Wowza Server Error : " + err.Error())
	}
	//-------------------------------------------------------------------------------

	tx.Save(webinar_site)

	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")

}


// 자료 다운로드
func (api WebinarSiteAPI) DownloadSiteFile(c echo.Context) error {

	logrus.Debug("******** DownloadSiteFile ")

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	logrus.Debug("******** idx : ", idx)

	site_file := &rdbModel.WebinarSiteFile{}

	tx := c.Get("Tx").(*gorm.DB)

	tx.Where("idx = ? ", idx).Find(site_file)

	logrus.Debug(">>>>>>>>> SavePath : ", site_file.SavePath)


	original_name := strings.Replace(site_file.OriginalName,","," ",-1)

	//s = strings.Replace(s,"\n","<br>",-1)

	return c.Attachment(site_file.SavePath, original_name)
}
