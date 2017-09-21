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
	"github.com/satori/go.uuid"
	"strings"
	"os"
	"path/filepath"
	"io"
	"strconv"
	"github.com/jinzhu/gorm"
)

type WebinarNoticeAPI struct {
	requestPost WebinarNoticeRequest
	//requestPut  ContentRequestPut
}


type WebinarNoticeRequest struct {
	WebinarSiteId		string 				`validate:"required" json:"webinar_site_id"`
	Title			string				`validate:"required" json:"title"`
	Content			string				`validate:"required" json:"content"`
	UseYN			string				`validate:"required" json:"use_yn"`
	WebinarNoticeFiles	[]WebinarNoticeFilesRequest	`validate:"required" json:"files"`

}


type WebinarNoticeFilesRequest struct {
	OriginalName		string		`validate:"required" json:"original_name"`
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

//type MaxValue struct {
//	Value		int
//}
//
// 웨비나 사이트 등록
func (api WebinarNoticeAPI) PostWebinarNotice(c echo.Context) error {

	logrus.Debug("########### PostWebinarNotice !!!")

	payload := &api.requestPost
	c.Bind(payload)

	logrus.Debug(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims 	 := apiModel.GetJWTClaims(c)

	logrus.Debug("#### TenantId : ", claims.TenantId)

	tx := c.Get("Tx").(*gorm.DB)

	webinar_notice_id := uuid.NewV4().String()

	// file이 있는경우 등록한다.
	if len(payload.WebinarNoticeFiles) > 0 {

		for _, v := range payload.WebinarNoticeFiles {

			webinar_notice_file := &rdbModel.WebinarNoticeFile {
				TenantId		: claims.TenantId,
				WebinarSiteId		: payload.WebinarSiteId,
				WebinarNoticeId		: webinar_notice_id,
				WebinarNoticeFileId 	: uuid.NewV4().String(),
				OriginalName		: v.OriginalName,
				SavePath		: v.SavePath,
				WebPath			: v.WebPath,
			}

			tx.Create(webinar_notice_file)
		}
	}

	// 공지사항 등록
	webinar_notice := &rdbModel.WebinarNotice{
		TenantId: claims.TenantId,
		WebinarSiteId: payload.WebinarSiteId,
		WebinarNoticeId:webinar_notice_id,
		Title:payload.Title,
		Content:payload.Content,
		UseYN:payload.UseYN,
		UpdatedId:claims.UserId,
	}

	tx.Create(webinar_notice)

	return handler.APIResultHandler(c, true, http.StatusCreated, "INSERT OK")

}

// 자료 업로드
func (api WebinarNoticeAPI) UploadNoticeFile(c echo.Context) error {

	logrus.Debug(">>>>> UploadNoticeFile !!")

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	notice_files := []rdbModel.WebinarNoticeFile{}

	var notice_file rdbModel.WebinarNoticeFile

	for i, file := range files {

		logrus.Debug("##### i : ", i)
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dstDir 		:= appConfig.Config.APP.DownloadRoot
		originFileName 	:= file.Filename
		originFileExt 	:= originFileName[strings.LastIndex(originFileName, ".")+1:]

		generateName := "notice_" + uuid.NewV4().String() + "." + originFileExt

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
		notice_file.OriginalName = originFileName
		notice_file.FileSize = generatedFile.Size()
		notice_file.SavePath = dstDir + string(filepath.Separator) + generateName
		notice_file.WebPath = "/download/" + generateName


		notice_files = append(notice_files, notice_file)
	}

	return handler.APIResultHandler(c, true, http.StatusOK, notice_files)
}

// 자료 다운로드
func (api WebinarNoticeAPI) DownloadNoticeFile(c echo.Context) error {

	logrus.Debug("******** DownloadNoticeFile ")

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	logrus.Debug("******** idx : ", idx)

	notice_file := &rdbModel.WebinarNoticeFile{}

	tx := c.Get("Tx").(*gorm.DB)

	tx.Where("idx = ? ", idx).Find(notice_file)

	logrus.Debug(">>>>>>>>> SavePath : ", notice_file.SavePath)

	//f, err := os.Open(notice_file.SavePath)
	//if err != nil {
	//	return err
	//}

	//data, _ := ioutil.ReadFile(notice_file.SavePath)

	//return c.Blob(http.StatusOK, "text/plain", data)

	//return c.Stream(http.StatusOK, "image/png", f)

	original_name := strings.Replace(notice_file.OriginalName,","," ",-1)

		//s = strings.Replace(s,"\n","<br>",-1)

	return c.Attachment(notice_file.SavePath, original_name)
}

 // 공지사항 리스트 조회
func (api WebinarNoticeAPI) GetWebinarNotices (c echo.Context) error  {

	claims 	 := apiModel.GetJWTClaims(c)

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	webinar_site_id		:= c.QueryParam("webinar_site_id");

	//start_date		:= c.QueryParam("start_date")
	//end_date		:= c.QueryParam("end_date")

	title			:= c.QueryParam("title")
	use_yn			:= c.QueryParam("use_yn")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### webinar_site_id : ", webinar_site_id)
	//logrus.Debug("##### start_date : ", start_date)
	//logrus.Debug("##### end_date : ", end_date)
	logrus.Debug("##### title : ", title)
	logrus.Debug("##### use_yn : ", use_yn)

	tx := c.Get("Tx").(*gorm.DB)

	webinar_notices    := []rdbModel.WebinarNotice{}

	tx  = tx.Where("tenant_id = ? and webinar_site_id = ? ", claims.TenantId, webinar_site_id)

	//if start_date != "" {
	//	tx = tx.Where("created_at BETWEEN ? AND ?", start_date, end_date)
	//}

	if title != "" {
		tx = tx.Where("title LIKE ? ", "%" + title + "%")
	}

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

// 공지사항 한건 조회
func (api WebinarNoticeAPI) GetWebinarNotice (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	webinar_notice  := &rdbModel.WebinarNotice{}

	//var count = 0
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("idx = ? ", idx ).Preload("WebinarNoticeFiles").
		Find(webinar_notice).RecordNotFound() {

		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR NOTICE  NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, webinar_notice)
}


 // 공지사항 삭제
func (api WebinarNoticeAPI) DeleteWebinarNotice (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	notice_file 	:= &rdbModel.WebinarNoticeFile{}
	notice 	        := &rdbModel.WebinarNotice{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("idx = ? ", idx ).Find(notice).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_NOTICE NOT FOUND")
	}

	if tx.Delete(notice_file, "webinar_site_id = ? and webinar_notice_id = ? ", notice.WebinarSiteId, notice.WebinarNoticeId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	if tx.Delete(notice, "webinar_site_id = ? and idx = ? ", notice.WebinarSiteId, idx).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "CONTENT NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}


// 공지사항 수정
func (api WebinarNoticeAPI) PutWebinarNotice (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	claims 	 := apiModel.GetJWTClaims(c)

	payload := &api.requestPost
	c.Bind(payload)

	if err   := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx := c.Get("Tx").(*gorm.DB)

	notice_file 	:= &rdbModel.WebinarNoticeFile{}
	notice 	        := &rdbModel.WebinarNotice{}

	if tx.Where("idx = ? ", idx ).Find(notice).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "WEBINAR_NOTICE NOT FOUND")
	}

	// file이 있는경우 등록한다.
	if len(payload.WebinarNoticeFiles) > 0 {

		// 1. 자료를 모두 삭제한다.
		if tx.Delete(notice_file, "webinar_site_id = ? and webinar_notice_id = ? ", notice.WebinarSiteId, notice.WebinarNoticeId).RowsAffected == 0 {
			//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
		}

		for _, v := range payload.WebinarNoticeFiles {

			webinar_notice_file := &rdbModel.WebinarNoticeFile {
				TenantId		: claims.TenantId,
				WebinarSiteId		: payload.WebinarSiteId,
				WebinarNoticeId		: notice.WebinarNoticeId,
				WebinarNoticeFileId 	: uuid.NewV4().String(),
				OriginalName		: v.OriginalName,
				SavePath		: v.SavePath,
				WebPath			: v.WebPath,
				UpdatedId		: claims.UserId,
			}

			tx.Create(webinar_notice_file)
		}
	}

	notice.Title = payload.Title
	notice.Content = payload.Content
	notice.UseYN = payload.UseYN

	tx.Save(notice)

	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")

}

