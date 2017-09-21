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
	"fmt"
	"path/filepath"
	"strings"
	"strconv"
	//"github.com/disintegration/imaging"
	"os"
	"time"
	//"image/png"
	//"github.com/nfnt/resize"
	//"image/jpeg"
	"github.com/disintegration/imaging"
)

type ContentAPI struct {
	requestPost ContentRequestPost
	requestPut  ContentRequestPut
}

type ContentRequestPost struct {
	CategoryId		string		`validate:"required" json:"category_id"`
	ContentName 		string 		`validate:"required" json:"content_name"`
	Tags    		[]ContentTagRequest 	`json:"tags"`
	OriginFileName 		string 		`validate:"required" json:"origin_filename"`
	GeneratedFileName 	string 		`validate:"required" json:"generated_filename"`
	Type 			string 		`validate:"required" json:"file_type"`
	Ext 			string 		`validate:"required" json:"file_ext"`
	Size 			int64		`validate:"required" json:"file_size"`
	FilePath 		string 		`validate:"required" json:"file_path"`
	Duration 		string 		`validate:"required" json:"duration"`
	MpegDash		string 		`json:"mpeg_dash"`
	RTMP			string  	`json:"rtmp"`
	HLS			string		`json:"hls"`
	HDS             	string		`json:"hds"`
	IOS         		string 		`json:"ios"`
	Android         	string  	`json:"android"`

}

type ContentTagRequest struct {
	Idx				int 		`json:"idx"`
	Name			string 		`json:"name"`
}

type ContentRequestPut struct {
	CategoryId		string		`validate:"required" json:"category_id"`
	ContentName 		string 		`validate:"required" json:"content_name"`
	Tags    		[]ContentTagRequest 	`json:"tags"`
	ThumbChange		string		`validate:"required" json:"thumb_change"`
	ThumbType		string		`json:"thumb_type"`
	ThumbTime		string		`json:"thumb_time"`
	ThumbNails		[]ThumbNailRequest		`json:"thumbnails"`


}

type ThumbNailRequest struct {
	ContentID 		string			`json:"content_id"`
	SavePath		string 			`json:"save_path"`
	WebPath			string 			`json:"web_path"`
	//Size 			int64		`json:"size"`
	Resolution  		string 		`json:"resolution"`
}

// 컨텐츠 등록
func (api ContentAPI) PostContent(c echo.Context) error {

	fmt.Println(">>>>>>>>>>>>>>>>>> PostContent !!!")

	//fmt.Println(c.Request().GetBody())

	payload := &api.requestPost
	c.Bind(payload)

	fmt.Println(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims 	 := apiModel.GetJWTClaims(c)

	// 썸네일 생성
	// 5초를 표준으로 가져간다.
	var thumbPath string

	videoPath := payload.FilePath + string(filepath.Separator) + payload.GeneratedFileName

	//"640x320",
	//"320x240"
	// 썸네일 2개만 만들자.
	// 썸네일 생성
	//------------------------------------------------------------------------------------------------

	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	thumbPath = appConfig.Config.APP.ThumbNailRoot + string(filepath.Separator) + "thumb_"+ payload.GeneratedFileName[:strings.Index(payload.GeneratedFileName, ".")] + "_"+timestamp+"_640x320" + ".png"
	err := appLibs.FFMpegExec{}.RunCreateThumbNail(videoPath, thumbPath, "00:00:05", "640x320")
	if err != nil {

	}

	thumb_640x320_path := thumbPath
	thumb_640x320_file := "thumb_"+payload.GeneratedFileName[:strings.Index(payload.GeneratedFileName, ".")] + "_" + timestamp + "_640x320" + ".png"

	thumbPath = appConfig.Config.APP.ThumbNailRoot + string(filepath.Separator) + "thumb_"+ payload.GeneratedFileName[:strings.Index(payload.GeneratedFileName, ".")] + "_" + timestamp +  "_320x240" + ".png"
	err = appLibs.FFMpegExec{}.RunCreateThumbNail(videoPath, thumbPath, "00:00:05", "320x240")
	if err != nil {

	}

	thumb_320x240_path := thumbPath
	thumb_320x240_file := "thumb_"+ payload.GeneratedFileName[:strings.Index(payload.GeneratedFileName, ".")] + "_" + timestamp + "_320x240" + ".png"
	//------------------------------------------------------------------------------------------------

	logrus.Debug("#### TenantId : ", claims.TenantId)

	content_id := uuid.NewV4().String()

	// 태그가 존재한다면..
	tags := []rdbModel.ContentTag{}
	if len(payload.Tags) > 0 {
		for _, v := range payload.Tags {

			s := append (tags, rdbModel.ContentTag{

				TenantId		: claims.TenantId,
				ContentId		: content_id,
				Name 			: v.Name,
				UpdatedId       : claims.UserId,
			})
			tags = s
		}
	}

	content := &rdbModel.Content{
		TenantID		:claims.TenantId,
		CategoryId		:payload.CategoryId,
		ContentId		:content_id,
		ContentName		:payload.ContentName,
		OriginFileName		:payload.OriginFileName,
		GeneratedFileName	:payload.GeneratedFileName,
		Type			:payload.Type,
		Ext			:payload.Ext,
		FilePath		:payload.FilePath,
		Duration		:payload.Duration,
		Size			:payload.Size,
		UpdatedId		:claims.UserId,
		ThumbType		:"1",
		ThumbTime		:"00:00:05",
		ThumbNails		:[]rdbModel.ThumbNail{
			{
				ContentID:content_id,
				Resolution:"640x320",
				SavePath:thumb_640x320_path,
				WebPath:"/thumbnail/"+thumb_640x320_file,
			},
			{
				ContentID:content_id,
				Resolution:"320x240",
				SavePath:thumb_320x240_path,
				WebPath:"/thumbnail/"+thumb_320x240_file,
			},
		},
		Stream			:rdbModel.Stream {
			ContentId	: content_id,
			MpegDash	: payload.MpegDash,
			RTMP		: payload.RTMP,
			HLS		: payload.HLS,
			HDS       	: payload.HDS,
			IOS         	: payload.IOS,
			Android         : payload.Android,
		},
		ContentTags: tags,
	}

	tx := c.Get("Tx").(*gorm.DB)

	if !tx.Where("content_id = ?",
		content.ContentId).Find(content).RecordNotFound() {
		return echo.NewHTTPError(http.StatusConflict, "Already exists Content ")
	}

	tx.Create(content)

	return handler.APIResultHandler(c, true, http.StatusCreated, content)

}


// 컨텐츠 업로드
func (api ContentAPI) UploadContent(c echo.Context) error {

	logrus.Debug(">>>>> UPLOAD CONTENT !!")

	uploadInfo := &appLibs.UploadInfo{}
	uploadInfo.UploadContent(c)

	return handler.APIResultHandler(c, true, http.StatusOK, uploadInfo)
}

// 컨텐츠 리스트 조회
func (api ContentAPI) GetContents (c echo.Context) error  {

	claims 	 := apiModel.GetJWTClaims(c)

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	start_date		:= c.QueryParam("start_date")
	end_date		:= c.QueryParam("end_date")

	tenant_id		:= claims.TenantId
	category_id		:= c.QueryParam("category_id");

	content_name 		:= c.QueryParam("content_name")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### start_date : ", start_date)
	logrus.Debug("##### end_date : ", end_date)
	logrus.Debug("##### tenant_id : ", tenant_id)
	logrus.Debug("##### category_id : ", category_id)
	logrus.Debug("##### content_name : ", content_name)

	contents     	:= []rdbModel.Content{}
	category  	:= &rdbModel.Category{}

	tx := c.Get("Tx").(*gorm.DB)

	tx.Where("tenant_id = ? and category_id = ? ", tenant_id, category_id).Find(category)

	tx = tx.Where("tenant_id = ? ", tenant_id)

	logrus.Debug("category.ParentIdx : ", category.ParentIdx)

	if category.ParentIdx != 0 {
		if category_id != "" {
			tx = tx.Where("category_id = ? ", category_id)
		}
	}

	if content_name != "" {
		tx = tx.Where("content_name LIKE ? ", "%" + content_name + "%")
	}

	if start_date != "" {
		//tx = tx.Where("created_at BETWEEN ? AND ?", start_date, end_date)
		tx = tx.Where("created_at LIKE ?", start_date + "%")
	}

	var count = 0
	tx.Find(&contents).Count(&count)

	tx.Preload("ContentTags").Preload("Stream").Preload("ThumbNails").Order("idx desc").Offset(offset).Limit(limit).Find(&contents)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{"total_count": count,
			"rows": contents})

}

// 컨텐츠 한건 조회
func (api ContentAPI) GetContent (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	content := &rdbModel.Content{}

	//var count = 0
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Preload("ContentTags").Preload("Stream").Preload("ThumbNails").Where("idx = ? ", idx ).Find(content).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "Content NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, content)
}

// 컨텐츠 삭제
func (api ContentAPI) DeleteContent (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	content   := &rdbModel.Content{}
	stream    := &rdbModel.Stream{}
	thumbnail := &rdbModel.ThumbNail{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("idx = ? ", idx ).Find(content).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "Content NOT FOUND")
	}

	// 스케쥴에 존재한다면 삭제 불가
	schedule_orders := []rdbModel.ScheduleOrder{}
	if !tx.Where("content_id = ? ", content.ContentId ).Find(&schedule_orders).RecordNotFound() {
		return echo.NewHTTPError(http.StatusBadRequest, "SCHEDULE_EXISTS")
	}

	if tx.Delete(stream, "content_id = ?", content.ContentId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "STREAM NOT FOUND")
	}

	if tx.Delete(thumbnail, "content_id = ?", content.ContentId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "CONTENT NOT FOUND")
	}

	if tx.Delete(content, "idx = ?", idx).RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "CONTENT NOT FOUND")
	}

	// 기존 파일 모두 삭제
	pattern := appConfig.Config.APP.ThumbNailRoot + string(filepath.Separator) + "*" + content.ContentId + "*"
	matches, err := filepath.Glob(pattern)

	if err != nil {
		fmt.Println(err)
	}

	for _, value := range matches {
		fmt.Println(value)
		if _, err := os.Stat(value); err == nil {
			os.Remove(value)
		}
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}


// 컨텐츠 수정
func (api ContentAPI) PutContent (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	claims 	 := apiModel.GetJWTClaims(c)


	payload := &api.requestPut
	c.Bind(payload)

	if err   := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx := c.Get("Tx").(*gorm.DB)

	content := &rdbModel.Content{}

	if tx.Where("idx = ? ", idx ).Find(content).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "Content NOT FOUND")
	}

	if payload.ThumbTime != content.ThumbTime {
		payload.ThumbChange = "Y"
	}

	logrus.Debug(" ######### payload.ThumbChange : ", payload.ThumbChange)
	logrus.Debug(" ######### payload.ThumbType : ", payload.ThumbType)
	logrus.Debug(" ######### payload.ThumbTime : ", payload.ThumbTime)
	logrus.Debug(" ######### content.ThumbTime : ", content.ThumbTime)

	// content 업로드
	if payload.ThumbChange == "Y" {
		// 변경이 되었다면..

		thumbnails := [] rdbModel.ThumbNail{}

		tx.Where("content_id = ?", content.ContentId).Find(&thumbnails)

		for _, value := range thumbnails {
			// 일단 물리적으로 삭제한다.
			if _, err := os.Stat(value.SavePath); err == nil {
				// path/to/whatever exists
				os.Remove(value.SavePath)
			}
		}

		// 직접 입력이라면... 1:시간입력 -> 2: 직접 업로드
		if payload.ThumbType == "1" {

			// 썸네일 생성
			//------------------------------------------------------------------------------------------------
			videoPath := content.FilePath + content.GeneratedFileName

			timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

			thumbPath := appConfig.Config.APP.ThumbNailRoot + string(filepath.Separator) + "thumb_"+ content.GeneratedFileName[:strings.Index(content.GeneratedFileName, ".")] + "_" + timestamp +  "_640x320" + ".png"
			err := appLibs.FFMpegExec{}.RunCreateThumbNail(videoPath, thumbPath, payload.ThumbTime, "640x320")
			if err != nil {

			}

			thumb_640x320_path := thumbPath
			thumb_640x320_file := "thumb_"+content.GeneratedFileName[:strings.Index(content.GeneratedFileName, ".")] + "_" + timestamp +  "_640x320" + ".png"

			thumbPath = appConfig.Config.APP.ThumbNailRoot + string(filepath.Separator) + "thumb_"+ content.GeneratedFileName[:strings.Index(content.GeneratedFileName, ".")] + "_" + timestamp +  "_320x240" + ".png"
			err = appLibs.FFMpegExec{}.RunCreateThumbNail(videoPath, thumbPath, payload.ThumbTime, "320x240")
			if err != nil {

			}

			thumb_320x240_path := thumbPath
			thumb_320x240_file := "thumb_"+ content.GeneratedFileName[:strings.Index(content.GeneratedFileName, ".")] + "_" + timestamp +  "_320x240" + ".png"
			//------------------------------------------------------------------------------------------------

			thumbnail := &rdbModel.ThumbNail{
				ContentID: content.ContentId,
			}

			if tx.Delete(thumbnail, "content_id = ?", thumbnail.ContentID).RowsAffected == 0 {
				//return echo.NewHTTPError(http.StatusNotFound, "Thumbnail NOT FOUND")
			}


			thumbnail.Resolution = "640x320"
			thumbnail.SavePath = thumb_640x320_path
			thumbnail.WebPath  = "/thumbnail/"+thumb_640x320_file

			tx.Create(thumbnail)

			thumbnail2 := &rdbModel.ThumbNail{
				ContentID: content.ContentId,
			}

			thumbnail2.Resolution = "320x240"
			thumbnail2.SavePath = thumb_320x240_path
			thumbnail2.WebPath  = "/thumbnail/"+thumb_320x240_file

			tx.Create(thumbnail2)

		} else {

			// 직접 업로드

			//ThumbNailRequest

			// 태그가 존재한다면..
			//tags := []rdbModel.ContentTag{}
			if len(payload.ThumbNails) > 0 {
				for _, v := range payload.ThumbNails {

					thumbnail := &rdbModel.ThumbNail{
						ContentID: v.ContentID,
						Resolution:v.Resolution,
						SavePath:v.SavePath,
						WebPath:v.WebPath,
					}

					tx.Save(thumbnail)
				}
			}

		}

	}

	tag := &rdbModel.ContentTag{}

	// 4. 태그 삭제
	if tx.Delete(tag, "content_id = ?", content.ContentId).RowsAffected == 0 {

	}

	tags := []rdbModel.ContentTag{}
	if len(payload.Tags) > 0 {
		for _, v := range payload.Tags {

			s := append (tags, rdbModel.ContentTag{

				TenantId		: claims.TenantId,
				ContentId		: content.ContentId,
				Name 			: v.Name,
				UpdatedId       : claims.UserId,
			})
			tags = s
		}
	}

	if tx.Model(content).Where("idx = ? ", idx).
		Updates(rdbModel.Content{
			CategoryId	:payload.CategoryId,
			ContentName	:payload.ContentName,
			ContentTags :tags,
			ThumbType	:payload.ThumbType,
			ThumbTime	:payload.ThumbTime,
	}).RowsAffected == 0 {

		return echo.NewHTTPError(http.StatusNotFound, "Cotent NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")

}

// 컨텐츠 썸네일 직접 업로드 수정
// 컨텐츠 업로드
func (api ContentAPI) UploadThumbnail(c echo.Context) error {

	logrus.Debug(">>>>> UPLOAD Thumbnail !!")

	content_id := c.FormValue("content_id")

	uploadInfo := &appLibs.UploadInfo{}
	err := uploadInfo.UploadThumbnail(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tx := c.Get("Tx").(*gorm.DB)

	thumbnails := [] rdbModel.ThumbNail{}

	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	fmt.Println("timestamp : ", timestamp)

	tx.Where("content_id = ?", content_id).Find(&thumbnails)

	for _, value := range thumbnails {
		// 일단 물리적으로 삭제한다.
		if _, err := os.Stat(value.SavePath); err == nil {
			// path/to/whatever exists
			os.Remove(value.SavePath)
		}
	}

	// 썸네일 모두 삭제하고 해당 썸네일로 INSERT 한다.
	thumbnail := rdbModel.ThumbNail{
		ContentID: content_id,
	}


	if tx.Delete(thumbnail, "content_id = ?", thumbnail.ContentID).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "Thumbnail NOT FOUND")
	}

	// 640 320 generate
	// -------------------------------------------------------------------------

	// -------------------------------------------------------------------------
	orgin_filepath := uploadInfo.FilePath  + string(filepath.Separator) + uploadInfo.GeneratedName

	srcImage, err := imaging.Open(orgin_filepath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	fmt.Println("origin_filepath : ", orgin_filepath)
	//
	thumb_640x320_path := uploadInfo.FilePath
	thumb_640x320_file := uploadInfo.GeneratedName[:strings.Index(uploadInfo.GeneratedName, ".")] + "_" + timestamp +  "_640x320." + uploadInfo.FileExt

	fmt.Println("640 gen !!")
	fmt.Println("thumb_640x320_path : ", thumb_640x320_path)
	fmt.Println("thumb_640x320_file : ", thumb_640x320_file)

	dstImage640 := imaging.Resize(srcImage, 640, 320, imaging.Lanczos)

	fmt.Println("uploadInfo.FilePath + string(filepath.Separator) + thumb_640x320_file : ", uploadInfo.FilePath + string(filepath.Separator) + thumb_640x320_file)

	err = imaging.Save(dstImage640, uploadInfo.FilePath + string(filepath.Separator) + thumb_640x320_file)

	if err != nil {
		fmt.Println(err)
		return err
	}

	//--------------------------------------------------------------------------------------------

	thumbnail.Resolution = "640x320"
	thumbnail.SavePath = thumb_640x320_path + string(filepath.Separator) + thumb_640x320_file
	thumbnail.WebPath  = "/thumbnail/"+thumb_640x320_file


	//--------------------------------------------------------------------------------------------


	thumbnail2 := rdbModel.ThumbNail{
		ContentID: content_id,
	}
	//
	thumb_320x240_path := uploadInfo.FilePath
	thumb_320x240_file := uploadInfo.GeneratedName[:strings.Index(uploadInfo.GeneratedName, ".")] + "_" + timestamp +  "_320x240." + uploadInfo.FileExt
	//--------------------------------------------------------------------------------------------



	dstImage320 := imaging.Resize(srcImage, 320, 240, imaging.Lanczos)

	err = imaging.Save(dstImage320, uploadInfo.FilePath + string(filepath.Separator) + thumb_320x240_file)

	if err != nil {
		fmt.Println(err)
		return err
	}

	//--------------------------------------------------------------------------------------------
	thumbnail2.Resolution = "320x240"
	thumbnail2.SavePath = thumb_320x240_path + string(filepath.Separator) + thumb_320x240_file
	thumbnail2.WebPath  = "/thumbnail/"+thumb_320x240_file


	//--------------------------------------------------------------------------------------------

	// 원본삭제
	os.Remove(uploadInfo.FilePath  + string(filepath.Separator) + uploadInfo.GeneratedName)

	//tx.Create(thumbnail)
	//tx.Create(thumbnail2)
	thumbnail_results := []rdbModel.ThumbNail{}

	s1 := append (thumbnail_results, thumbnail)
	thumbnail_results = s1

	s2 := append (thumbnail_results, thumbnail2)
	thumbnail_results = s2

	fmt.Println(thumbnail_results)

	return handler.APIResultHandler(c, true, http.StatusOK,
		thumbnail_results)
}