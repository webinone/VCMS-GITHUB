package api

import (
	rdbModel "VCMS/apps/models/rdb"
	apiModel "VCMS/apps/models/api"
	appLibs "VCMS/apps/libs"
	appConfig "VCMS/apps/config"
	"github.com/labstack/echo"
	"net/http"
	"VCMS/apps/handler"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"strconv"
	"github.com/Sirupsen/logrus"
	"path/filepath"
	"fmt"
)

type ScheduleAPI struct {
	request ScheduleRequest
}

type ScheduleRequest struct {
	ScheduleId   		string 				`json:"schedule_id"`
	ChannelId		string				`validate:"required" json:"channel_id"`
	ScheduleName 		string 				`validate:"required" json:"schedule_name"`
	StartDateTime 		string 				`validate:"required" json:"start_datetime"`
	ScheduleOrders    	[]ScheduleOrderRequest 		`validate:"required" json:"schedule_orders"`
}

type ScheduleOrderRequest struct {
	ScheduleId   		string 				`json:"schedule_id"`
	ScheduleOrderId		string 				`json:"schedule_order_id"`
	ContentId    		string 				`validate:"required" json:"content_id"`
	GeneratedFileName	string				`json:"generated_filename"`
	StartSec 		string 				`validate:"required" json:"start_sec"`
	EndSec 			string 				`validate:"required" json:"end_sec"`
}

// 스케쥴 등록
func (api ScheduleAPI) PostSchedule(c echo.Context) error {

	payload := &api.request
	c.Bind(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	schedule_id := uuid.NewV4().String()

	claims 	 := apiModel.GetJWTClaims(c)

	schedule_orders := []rdbModel.ScheduleOrder{}

	videos 		:= []appLibs.Video{}

	total_time	:= 0

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
				Src    : "mp4:assets/"+v.ContentId+".mp4",
				Start  : v.StartSec,
				Length : v.EndSec,
			})
			videos = video

			number_end_sec, _ 	:= strconv.Atoi(v.EndSec)
			number_start_sec, _ 	:= strconv.Atoi(v.StartSec)

			total_time += number_end_sec - number_start_sec
		}

	}

	schedule := &rdbModel.Schedule {
		TenantID	:	claims.TenantId,
		ChannelId	:	payload.ChannelId,
		ScheduleId	:	schedule_id,
		Name	  	:	payload.ScheduleName,
		StartDateTime 	:	payload.StartDateTime,
		TotalTime	:	strconv.Itoa(total_time),
		ScheduleOrders  :	schedule_orders,
		UpdatedId       : 	claims.UserId,
		Stream			:rdbModel.Stream {
			ChannelId	: payload.ChannelId,
			MpegDash	: "http://"+ appConfig.Config.WOWZA.StreamUrl + "/"+ claims.TenantId  +"/" + payload.ChannelId + "/manifest.mpd",
			RTMP		: "rtmp://"+ appConfig.Config.WOWZA.StreamUrl + "/"+ claims.TenantId  +"/" + payload.ChannelId,
			HDS       	: "rtmp://"+ appConfig.Config.WOWZA.StreamUrl + "/"+ claims.TenantId  +"/" + payload.ChannelId + "/manifest.f4m",
			IOS         	: "http://"+ appConfig.Config.WOWZA.StreamUrl + "/"+ claims.TenantId  +"/" + payload.ChannelId + "/playlist.m3u8",
			Android         : "rtsp://"+ appConfig.Config.WOWZA.StreamUrl + "/"+ claims.TenantId  +"/" + payload.ChannelId,
		},
	}

	// TODO : 스케쥴 생성시에 WOWZA 수정한다
	//-------------------------------------------------------------------------------
	smilFilePath := appConfig.Config.APP.ContentsRoot + string(filepath.Separator) + claims.TenantId + string(filepath.Separator) + "streamschedule.smil"

	err := appLibs.ScheduleSmilParser{ FilePath: smilFilePath}.
		UpdateScheduleSmil( payload.ChannelId, "false",  schedule_id, payload.StartDateTime, videos)

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

	// 기존에 등록되어 있는 Tenant인지 체크 한다.
	tx := c.Get("Tx").(*gorm.DB)

	tx.Create(schedule)

	return handler.APIResultHandler(c, true, http.StatusCreated, schedule)
}

// 스케쥴 수정
func (api ScheduleAPI) PutSchedule(c echo.Context) error {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	payload := &api.request
	c.Bind(payload)
	claims 	 := apiModel.GetJWTClaims(c)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx := c.Get("Tx").(*gorm.DB)

	schedule 	:= &rdbModel.Schedule{}

	if tx.Where("idx = ? ",
		idx ).Find(schedule).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "Tenant NOT FOUND")
	}

	schedule_order 	:= &rdbModel.ScheduleOrder{
		ScheduleId: schedule.ScheduleId,
	}

	tx.Delete(schedule_order, "schedule_id = ?", schedule.ScheduleId)

	schedule_orders := []rdbModel.ScheduleOrder{}
	videos 		:= []appLibs.Video{}

	total_time 	:= 0

	if len(payload.ScheduleOrders) > 0 {

		for i, v := range payload.ScheduleOrders {

			s := append (schedule_orders, rdbModel.ScheduleOrder{

				TenantID	: claims.TenantId,
				ScheduleId 	: schedule.ScheduleId,
				ScheduleOrderId : uuid.NewV4().String(),
				Order		: i+1 ,
				ContentId 	: v.ContentId,
				StartSec	: v.StartSec,
				EndSec		: v.EndSec,
				UpdatedId       : claims.UserId,
			})
			schedule_orders = s

			video := append (videos, appLibs.Video{
				Src    : "mp4:assets/"+v.ContentId+".mp4",
				Start  : v.StartSec,
				Length : v.EndSec,
			})
			videos = video

			number_end_sec, _ 	:= strconv.Atoi(v.EndSec)
			number_start_sec, _ 	:= strconv.Atoi(v.StartSec)

			total_time += number_end_sec - number_start_sec
		}

	}

	// TODO :  WOWZA 수정한다
	//-------------------------------------------------------------------------------
	smilFilePath := appConfig.Config.APP.ContentsRoot + string(filepath.Separator) + claims.TenantId + string(filepath.Separator) + "streamschedule.smil"

	err := appLibs.ScheduleSmilParser{ FilePath: smilFilePath}.
		UpdateScheduleSmil( payload.ChannelId, "false",  schedule.ScheduleId, payload.StartDateTime, videos)

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

	// 스케쥴 Order 삭제후 재 입력한다.
	if tx.Model(schedule).Where("idx = ? ", idx).
		Updates(
			rdbModel.Schedule{
				ChannelId:payload.ChannelId,
				Name: payload.ScheduleName,
				StartDateTime	:	payload.StartDateTime,
				TotalTime	:	strconv.Itoa(total_time),
				ScheduleOrders	:	schedule_orders,
				UpdatedId: claims.UserId,
			}).RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Schedule NOT FOUND")
	}

	// 스케쥴 입력


	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")
}

// 스케쥴 삭제
func (api ScheduleAPI) DeleteSchedule(c echo.Context) error {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)
	claims 	 := apiModel.GetJWTClaims(c)

	schedule 	:= &rdbModel.Schedule {}
	schedule_order 	:= &rdbModel.ScheduleOrder {}
	stream 		:= &rdbModel.Stream {}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("idx = ? ",
		idx ).Find(schedule).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "Schedule NOT FOUND")
	}

	tx.Delete(stream, "channel_id = ?", schedule.ChannelId)
	tx.Delete(schedule_order, "schedule_id = ?", schedule.ScheduleId)
	tx.Delete(schedule, "idx = ?", idx)

	// TODO : WOWZA 연동 삭제
	//-------------------------------------------------------------------------------
	smilFilePath := appConfig.Config.APP.ContentsRoot + string(filepath.Separator) + claims.TenantId + string(filepath.Separator) + "streamschedule.smil"

	err := appLibs.ScheduleSmilParser{ FilePath: smilFilePath}.
		DeleteScheduleSmil( schedule.ScheduleId )

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

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}

// 스케쥴 한건 조회
func (api ScheduleAPI) GetSchedule(c echo.Context) error {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	schedule := &rdbModel.Schedule{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Preload("ScheduleOrders").Preload("ScheduleOrders.Content").Preload("ScheduleOrders.Content.Stream").Where("idx = ? ",
		idx ).Find(schedule).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "Tenant NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, schedule)
}

// 스케쥴 여러건 조회
func (api ScheduleAPI) GetSchedules(c echo.Context) error {

	channel_id 		:= c.QueryParam("channel_id")

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	start_date		:= c.QueryParam("start_date")
	end_date		:= c.QueryParam("end_date")

	schedule_name 		:= c.QueryParam("schedule_name")

	logrus.Debug("##### channel_id : ", channel_id)
	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### start_date : ", start_date)
	logrus.Debug("##### end_date : ", end_date)
	logrus.Debug("##### schedule_name : ", schedule_name)

	schedules := &[]rdbModel.Schedule{}
	tx := c.Get("Tx").(*gorm.DB)

	tx = tx.Preload("Channel").Preload("ScheduleOrders").Preload("ScheduleOrders.Content").Preload("ScheduleOrders.Content.Stream")

	if schedule_name != "" {
		tx = tx.Where("name LIKE ?", "%" + schedule_name + "%")
	}

	if start_date != "" {
		tx = tx.Where("created_at BETWEEN ? AND ?", start_date, end_date)
	}

	if channel_id != "" {
		tx = tx.Where("channel_id = ?", channel_id)
	}

	var count = 0
	tx.Find(schedules).Count(&count)

	tx.Order("idx desc").Offset(offset).Limit(limit).Find(schedules)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{"total_count": count,
			"rows": schedules})
}
