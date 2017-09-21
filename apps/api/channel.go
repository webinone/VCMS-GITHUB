package api

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"VCMS/apps/handler"
	apiModel "VCMS/apps/models/api"
	rdbModel "VCMS/apps/models/rdb"
	appConfig "VCMS/apps/config"
	appLibs "VCMS/apps/libs"
	"github.com/satori/go.uuid"
	"strconv"
	//"path/filepath"
	//"os"
	//"fmt"
	"path/filepath"
	"fmt"
	"os"
)

type ChannelAPI struct {
	request ChannelRequest
}

type ChannelRequest struct {
	ChannelId 	string	`json:"channel_id"`
	Name 		string 	`validate:"required" json:"channel_name"`
	Bitrate		string	`json:"bitrate"`
}

// 채널 등록
func (api ChannelAPI) PostChannel (c echo.Context) error  {

	payload := &api.request
	c.Bind(payload)

	if err   := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	claims 	 := apiModel.GetJWTClaims(c)

	logrus.Debug("#### TenantId : ", claims.TenantId)

	channel_id := uuid.NewV4().String()

	mpegDashUrl 	:= "http://" + appConfig.Config.WOWZA.StreamUrl + "/" + claims.TenantId + "/" + channel_id + "/manifest.mpd"
	rtmpUrl 	:= "rtmp://" + appConfig.Config.WOWZA.StreamUrl + "/" + claims.TenantId
	hlsUrl 		:= "http://" + appConfig.Config.WOWZA.StreamUrl + "/" + claims.TenantId + "/" + channel_id + "/manifest.m3u8"
	hdsUrl 		:= "http://" + appConfig.Config.WOWZA.StreamUrl + "/" + claims.TenantId + "/" + channel_id + "/manifest.f4m"
	iosUrl 		:= "http://" + appConfig.Config.WOWZA.StreamUrl + "/" + claims.TenantId + "/" + channel_id + "/playlist.m3u8"
	androidUrl 	:= "rtsp://" + appConfig.Config.WOWZA.StreamUrl + "/" + claims.TenantId + "/" + channel_id

	channel := &rdbModel.Channel{
		TenantID	: claims.TenantId,
		ChannelId	: channel_id,
		Name		: payload.Name,
		Bitrate		: payload.Bitrate,
		UpdatedId	: claims.UserId,
		Stream			:rdbModel.Stream {
			ChannelId	: channel_id,
			MpegDash	: mpegDashUrl,
			RTMP		: rtmpUrl,
			HLS		: hlsUrl,
			HDS       	: hdsUrl,
			IOS         	: iosUrl,
			Android         : androidUrl,
		},
	}

	tx := c.Get("Tx").(*gorm.DB)

	if !tx.Where("name = ? ",
		channel.ChannelId ).Find(channel).RecordNotFound() {
		return echo.NewHTTPError(http.StatusConflict, "Already exists Channel ")
	}

	tx.Create(channel)

	// Stream 생성


	// WOWZA
	//---------------------------------------------------------------------------
	// 채널 파일이 존재하는지 확인한다.
	switchSmilFile := appConfig.Config.APP.ContentsRoot + string(filepath.Separator) + claims.TenantId + string(filepath.Separator) + "switch.smil"

	fmt.Println(switchSmilFile)

	if _, err := os.Stat(switchSmilFile); os.IsNotExist(err) {

		//if err != nil {
		//	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		//}
		// path/to/whatever does not exist
		appLibs.WowzaHttpClient{ APIUrl: appConfig.Config.WOWZA.ApiUrl }.CreateScheduleSwitchSmil(claims.TenantId, channel_id, payload.Bitrate)
	} else {

		restURI := appConfig.Config.WOWZA.ApiUrl + "/v2/servers/_defaultServer_/vhosts/_defaultVHost_/applications/" + claims.TenantId + "/smilfiles/switch"

		switchSmil := appLibs.SwitchSmil{}
		switchSmil.RestURI = restURI

		channels := []rdbModel.Channel{}

		tx.Where(" tenant_id = ? ", claims.TenantId).Find(&channels)

		for _, v := range channels {
			switchSmil.SmilStreams = append(switchSmil.SmilStreams,
				appLibs.SmilStream {
					Src: v.ChannelId,
					SystemBirate:v.Bitrate,
					Type:"video",
					RestURI: restURI,
				},
			)
		}

		fmt.Println(switchSmil)

		err := appLibs.WowzaHttpClient{ APIUrl: appConfig.Config.WOWZA.ApiUrl }.UpdateScheduleSwitchSmil(switchSmil)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}


		//appLibs.ScheduleSmilParser{ FilePath: switchSmilFile }.UpdateScheduleSwitchSmil(channel_id, payload.Bitrate)
	}

	//---------------------------------------------------------------------------

	return handler.APIResultHandler(c, true, http.StatusCreated, channel)
}

// 채널 한건 조회
func (api ChannelAPI) GetChannel (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	channel := &rdbModel.Channel{}

	//var count = 0
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("idx = ? ", idx ).Find(channel).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "Channel NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, channel)

}

// 채널 리스트 조회
func (api ChannelAPI) GetChannels (c echo.Context) error  {

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	created_at		:= c.QueryParam("created_at")
	//end_date		:= c.QueryParam("end_date")

	channel_name 		:= c.QueryParam("channel_name")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### created_at : ", created_at)
	//logrus.Debug("##### end_date : ", end_date)
	logrus.Debug("##### channel_name : ", channel_name)

	channels := []rdbModel.Channel{}
	tx := c.Get("Tx").(*gorm.DB)

	claims 	 := apiModel.GetJWTClaims(c)

	tx = tx.Where("tenant_id = ? ", claims.TenantId)

	if channel_name != "" {
		tx = tx.Where("name LIKE ? ", "%" + channel_name + "%")
	}

	if created_at != "" {
		tx = tx.Where("created_at LIKE ?", created_at + "%")
	}

	var count = 0
	tx.Find(&channels).Count(&count)

	tx.Order("idx desc").Offset(offset).Limit(limit).Find(&channels)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{"total_count": count,
			"rows": channels})

}

// 채널 수정
func (api ChannelAPI) PutChannel (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	payload := &api.request
	c.Bind(payload)

	if err   := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims 	 := apiModel.GetJWTClaims(c)


	tx := c.Get("Tx").(*gorm.DB)

	channel := &rdbModel.Channel{}

	if tx.Model(channel).Where("idx = ? ", idx).
		Updates(
			rdbModel.Channel{
				Name: payload.Name,
				Bitrate:payload.Bitrate,
				UpdatedId:claims.UserId,
	}).RowsAffected == 0 {

		return echo.NewHTTPError(http.StatusNotFound, "Channel NOT FOUND")
	}

	// WOWZA
	//----------------------------------------------------------------------------------------------------
	restURI := appConfig.Config.WOWZA.ApiUrl + "/v2/servers/_defaultServer_/vhosts/_defaultVHost_/applications/" + claims.TenantId + "/smilfiles/switch"

	switchSmil := appLibs.SwitchSmil{}
	switchSmil.RestURI = restURI

	channels := []rdbModel.Channel{}

	tx.Where(" tenant_id = ? ", claims.TenantId).Find(&channels)

	for _, v := range channels {
		switchSmil.SmilStreams = append(switchSmil.SmilStreams,
			appLibs.SmilStream {
				Src: v.ChannelId,
				SystemBirate:v.Bitrate,
				Type:"video",
				RestURI: restURI,
			},
		)
	}

	appLibs.WowzaHttpClient{ APIUrl: appConfig.Config.WOWZA.ApiUrl }.UpdateScheduleSwitchSmil(switchSmil)
	//----------------------------------------------------------------------------------------------------

	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")

}


// 채널 삭제
func (api ChannelAPI) DeleteChannel (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	channel := &rdbModel.Channel{}

	tx := c.Get("Tx").(*gorm.DB)

	claims 	 := apiModel.GetJWTClaims(c)

	if tx.Where("idx = ? ", idx ).Find(channel).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "Channel NOT FOUND")
	}

	if tx.Delete(channel, "idx = ?", idx).RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Channel NOT FOUND")
	}

	stream := &rdbModel.Stream{}
	if tx.Delete(stream, "channel_id = ?", channel.ChannelId).RowsAffected == 0 {
		//return echo.NewHTTPError(http.StatusNotFound, "Channel NOT FOUND")
	}

	// WOWZA
	//----------------------------------------------------------------------------------------------------
	restURI := appConfig.Config.WOWZA.ApiUrl + "/v2/servers/_defaultServer_/vhosts/_defaultVHost_/applications/" + claims.TenantId + "/smilfiles/switch"

	switchSmil := appLibs.SwitchSmil{}
	switchSmil.RestURI = restURI

	channels := []rdbModel.Channel{}

	tx.Where(" tenant_id = ? ", claims.TenantId).Find(&channels)

	for _, v := range channels {
		switchSmil.SmilStreams = append(switchSmil.SmilStreams,
			appLibs.SmilStream {
				Src: v.ChannelId,
				SystemBirate:v.Bitrate,
				Type:"video",
				RestURI: restURI,
			},
		)
	}

	appLibs.WowzaHttpClient{ APIUrl: appConfig.Config.WOWZA.ApiUrl }.UpdateScheduleSwitchSmil(switchSmil)
	//----------------------------------------------------------------------------------------------------

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}