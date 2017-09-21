package api

import (
	"github.com/labstack/echo"
	"github.com/Sirupsen/logrus"
	apiModel "VCMS/apps/models/api"
	rdbModel "VCMS/apps/models/rdb"
	appLibs "VCMS/apps/libs"
	appConfig "VCMS/apps/config"
	"time"
	"strconv"
	"fmt"
	"github.com/jinzhu/gorm"
	"VCMS/apps/handler"
	"net/http"
	"path/filepath"
)

type DashBoardAPI struct {
}


type ChartLineData struct {
	Time 		string
	Total 		float64
	Vod 		float64
	Live 		float64
}

type ChartLabelValueData struct {
	Label 		string
	Value 		float64
}

type ChartValueData struct {
	Value 		float64
}

type ChartLabelIdValueData struct {
	Label 		string
	Id 		string
	Value 		float64
}

// WOWZA 호출
func (api DashBoardAPI) GetCurrentApplicationConnection (c echo.Context) error  {


	return nil

}

// WOWZA 호출
func (api DashBoardAPI) GetCurrentChannelConnection (c echo.Context) error  {



	return nil

}

func (api DashBoardAPI) GetDashBoardData (c echo.Context) error  {


	//getMonthConnections(c)

	return nil

}



func (api DashBoardAPI) GetMonthConnections(c echo.Context) error {

	claims 	 := apiModel.GetJWTClaims(c)

	logrus.Debug("#### TenantId : ", claims.TenantId)

	offSet, _ := time.ParseDuration("+09.00h")

	now := time.Now().UTC().Add(offSet)

	now_year 	:= now.Year()
	now_month 	:= int(now.Month())
	now_day     	:= now.Day()

	nowDate := strconv.Itoa(now_year) + appLibs.StrPad(strconv.Itoa(now_month), 2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(now_day),2, "0","LEFT")
	fmt.Println(nowDate)


	monthAgo := now.AddDate(0, -1, 0)

	monthAgo_year	:= monthAgo.Year()
	monthAgo_month	:= int(monthAgo.Month())
	monthAgo_day	:= monthAgo.Day()

	monthAgoDate 	:= strconv.Itoa(monthAgo_year) + appLibs.StrPad(strconv.Itoa(monthAgo_month), 2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(monthAgo_day),2, "0","LEFT")

	fmt.Println(monthAgoDate)

	tx := c.Get("Tx").(*gorm.DB)

	chartLineData := []ChartLineData{}

	tx.Raw(`
		SELECT  a.time,
			a.total,
			(
			  SELECT count(*) FROM TB_VCMS_PLAY_LOG
			  WHERE content_id != '-'
			  AND group_year_month_day = a.time
			) as vod,
			(
			  SELECT count(*) FROM TB_VCMS_PLAY_LOG
			  WHERE channel_id != '-'
			  AND group_year_month_day = a.time
			) as live
		FROM
		  (
			  SELECT group_year_month_day AS time,
			     COUNT(*)                         AS total
			   FROM TB_VCMS_PLAY_LOG
			   WHERE tenant_id = ?
				 AND group_year_month_day between ? and ?
			   GROUP BY group_year_month_day
			   ORDER BY group_year_month_day DESC
		  ) a
		`, claims.TenantId, monthAgoDate, nowDate).Scan(&chartLineData) // (*sql.Rows, error)

	var chart_data string

	chart_data += "["


	for i, val := range chartLineData {

		fmt.Println(i)
		fmt.Println(val.Time)
		fmt.Println(val.Total)
		fmt.Println(val.Vod)
		fmt.Println(val.Live)

		chart_data += "["

		chart_data += `"` + fmt.Sprintf("%s/%s/%s", val.Time[:4], val.Time[4:6], val.Time[6:8]) + `",`
		chart_data += fmt.Sprintf("%.0f", val.Total) + ","
		chart_data += fmt.Sprintf("%.0f", val.Vod) + ","
		chart_data += fmt.Sprintf("%.0f", val.Live) + ""

		chart_data += "]"


		if i != len(chartLineData)-1 {
			chart_data += ","
		}

	}

	chart_data += "]"

	fmt.Println(chart_data)

	return handler.APIResultHandler(c, true, http.StatusCreated, chart_data)

}

func (api DashBoardAPI) GetWeeksConnections(c echo.Context) error {


	claims 	 := apiModel.GetJWTClaims(c)

	logrus.Debug("#### TenantId : ", claims.TenantId)

	offSet, _ := time.ParseDuration("+09.00h")

	now := time.Now().UTC().Add(offSet)

	now_year 	:= now.Year()
	now_month 	:= int(now.Month())
	now_day     	:= now.Day()

	nowDate := strconv.Itoa(now_year) + appLibs.StrPad(strconv.Itoa(now_month), 2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(now_day),2, "0","LEFT")
	fmt.Println(nowDate)

	monthAgo := now.AddDate(0, 0, -7)

	monthAgo_year	:= monthAgo.Year()
	monthAgo_month	:= int(monthAgo.Month())
	monthAgo_day	:= monthAgo.Day()

	monthAgoDate 	:= strconv.Itoa(monthAgo_year) + appLibs.StrPad(strconv.Itoa(monthAgo_month), 2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(monthAgo_day),2, "0","LEFT")

	fmt.Println(monthAgoDate)

	tx := c.Get("Tx").(*gorm.DB)

	chartLineData := []ChartLineData{}

	tx.Raw(`
		SELECT  a.time,
			a.total,
			(
			  SELECT count(*) FROM TB_VCMS_PLAY_LOG
			  WHERE content_id != '-'
			  AND group_year_month_day = a.time
			) as vod,
			(
			  SELECT count(*) FROM TB_VCMS_PLAY_LOG
			  WHERE channel_id != '-'
			  AND group_year_month_day = a.time
			) as live
		FROM
		  (
			  SELECT group_year_month_day AS time,
			     COUNT(*)                         AS total
			   FROM TB_VCMS_PLAY_LOG
			   WHERE tenant_id = ?
				 AND group_year_month_day between ? and ?
			   GROUP BY group_year_month_day
			   ORDER BY group_year_month_day DESC
		  ) a
		`, claims.TenantId, monthAgoDate, nowDate).Scan(&chartLineData) // (*sql.Rows, error)

	var chart_data string

	chart_data += "["


	for i, val := range chartLineData {

		fmt.Println(i)
		fmt.Println(val.Time)
		fmt.Println(val.Total)
		fmt.Println(val.Vod)
		fmt.Println(val.Live)

		chart_data += "["

		chart_data += `"` + fmt.Sprintf("%s/%s/%s", val.Time[:4], val.Time[4:6], val.Time[6:8]) + `",`
		chart_data += fmt.Sprintf("%.0f", val.Total) + ","
		chart_data += fmt.Sprintf("%.0f", val.Vod) + ","
		chart_data += fmt.Sprintf("%.0f", val.Live) + ""

		chart_data += "]"


		if i != len(chartLineData)-1 {
			chart_data += ","
		}

	}

	chart_data += "]"

	fmt.Println(chart_data)

	return handler.APIResultHandler(c, true, http.StatusCreated, chart_data)
}


func (api DashBoardAPI) GetTodayConnections(c echo.Context) error {

	claims 	 := apiModel.GetJWTClaims(c)

	logrus.Debug("#### TenantId : ", claims.TenantId)

	offSet, _ := time.ParseDuration("+09.00h")

	now := time.Now().UTC().Add(offSet)

	now_year 	:= now.Year()
	now_month 	:= int(now.Month())
	now_day     	:= now.Day()

	nowDate := strconv.Itoa(now_year) + appLibs.StrPad(strconv.Itoa(now_month), 2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(now_day),2, "0","LEFT")
	fmt.Println(nowDate)

	tx := c.Get("Tx").(*gorm.DB)

	// 일단 채널 목록을 가져온다.
	channels := []rdbModel.Channel{}

	tx.Where("tanant_id = ? ", claims.TenantId).Find(&channels)

	sqlString := `
	 		SELECT  a.time,
			  a.total,
		    	(
			  SELECT count(*) FROM TB_VCMS_PLAY_LOG
			  WHERE content_id != '-'
			  AND group_year_month_day_hour_minute = a.time
			) as vod,
			(
			  SELECT count(*) FROM TB_VCMS_PLAY_LOG
			  WHERE channel_id != '-'
			  AND group_year_month_day_hour_minute = a.time
			) as live
			FROM
			  (SELECT
			     group_year_month_day_hour_minute AS time,
			     COUNT(*)                         AS total
			   FROM TB_VCMS_PLAY_LOG
			   WHERE tenant_id = ?
				 AND group_year_month_day = ?
			   GROUP BY group_year_month_day_hour_minute
			   ORDER BY group_year_month_day_hour_minute DESC
			  ) a
			  `

	chartLineData := []ChartLineData{}

	tx.Raw(sqlString, claims.TenantId, nowDate).Scan(&chartLineData) // (*sql.Rows, error)

	var chart_data string

	chart_data += "["


	for i, val := range chartLineData {


		chart_data += "["


		chart_data += `"` + fmt.Sprintf("%s:%s", val.Time[8:10], val.Time[10:]) + `",`
		chart_data += fmt.Sprintf("%.0f", val.Total) + ","
		chart_data += fmt.Sprintf("%.0f", val.Vod) + ","
		chart_data += fmt.Sprintf("%.0f", val.Live) + ""

		chart_data += "]"


		if i != len(chartLineData)-1 {
			chart_data += ","
		}

	}

	chart_data += "]"

	fmt.Println(chart_data)

	return handler.APIResultHandler(c, true, http.StatusCreated, chart_data)

}


// 현재 접속자 현황 (WOWZA) -- 해결 -- DB로 해결 (24시간 데이터만 맞는걸로 한다.)
func (api DashBoardAPI) GetApplicationConnectionCount(c echo.Context) error {

	claims 	 := apiModel.GetJWTClaims(c)
	logrus.Debug("#### TenantId : ", claims.TenantId)

	offSet, _ := time.ParseDuration("+09.00h")

	now := time.Now().UTC().Add(offSet)

	now_year 	:= now.Year()
	now_month 	:= int(now.Month())
	now_day     	:= now.Day()
	now_hour	:= now.Hour()

	nowDate := strconv.Itoa(now_year) + appLibs.StrPad(strconv.Itoa(now_month), 2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(now_day),2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(now_hour),2, "0","LEFT")
	fmt.Println(nowDate)

	monthAgo := now.AddDate(0, 0, -1)

	monthAgo_year	:= monthAgo.Year()
	monthAgo_month	:= int(monthAgo.Month())
	monthAgo_day	:= monthAgo.Day()
	monthAgo_hour	:= monthAgo.Hour()

	monthAgoDate 	:= strconv.Itoa(monthAgo_year) + appLibs.StrPad(strconv.Itoa(monthAgo_month), 2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(monthAgo_day),2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(monthAgo_hour),2, "0","LEFT")

	fmt.Println(monthAgoDate)

	tx := c.Get("Tx").(*gorm.DB)

	chartValueData := ChartValueData{}

	tx.Raw(`
		SELECT COUNT(*) value
		  FROM TB_VCMS_PLAY_LOG
		  WHERE cclientid NOT IN (
		    SELECT cclientid
		    FROM TB_VCMS_ACCESS_LOG
		    WHERE xevent = 'disconnect'
		  )
		  AND tenant_id = ?
		  AND group_year_month_day_hour between ? and ?
		`, claims.TenantId, monthAgoDate, nowDate).Scan(&chartValueData) // (*sql.Rows, error)


	fmt.Println(chartValueData)

	return handler.APIResultHandler(c, true, http.StatusCreated,
		map[string]interface{}{"totalConnections": chartValueData.Value})
}


// 채널별 접속자 현황 -- 1분당 채널 현황 조회로 해결
func (api DashBoardAPI) GetChannelConnectionCount(c echo.Context) error {

	channel_id := c.Param("channel_id")

	claims 	 := apiModel.GetJWTClaims(c)

	result, err := appLibs.WowzaHttpClient{APIUrl:appConfig.Config.WOWZA.ApiUrl}.GetChannelConnections(claims.TenantId, channel_id)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Wowza Server Error : " + err.Error())
	}


	return handler.APIResultHandler(c, true, http.StatusOK, result)
}


// 웨비나별 접속자 현황


// 총 채널수 -- 해결
func (api DashBoardAPI) GetChannelsCount(c echo.Context) error {

	claims 	 := apiModel.GetJWTClaims(c)

	channels := []rdbModel.Channel{}
	tx := c.Get("Tx").(*gorm.DB)

	var count = 0

	tx.Where("tenant_id = ? ", claims.TenantId).Find(&channels).Count(&count)


	return handler.APIResultHandler(c, true, http.StatusOK, count)
}

// 스케쥴 수 -- 해결 (최근 10개)
func (api DashBoardAPI) GetSchedules(c echo.Context) error {

	claims 	 := apiModel.GetJWTClaims(c)

	schedules := []rdbModel.Schedule{}
	tx := c.Get("Tx").(*gorm.DB)

	tx.Where("tenant_id = ? ", claims.TenantId).Order("idx desc").Limit(10).Find(&schedules)


	return handler.APIResultHandler(c, true, http.StatusOK, schedules)

}

func (api DashBoardAPI) GetNowConnectPieChart (c echo.Context) error {


	claims 	 := apiModel.GetJWTClaims(c)
	logrus.Debug("#### TenantId : ", claims.TenantId)

	offSet, _ := time.ParseDuration("+09.00h")

	now := time.Now().UTC().Add(offSet)

	now_year 	:= now.Year()
	now_month 	:= int(now.Month())
	now_day     	:= now.Day()
	now_hour	:= now.Hour()

	nowDate := strconv.Itoa(now_year) + appLibs.StrPad(strconv.Itoa(now_month), 2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(now_day),2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(now_hour),2, "0","LEFT")
	fmt.Println(nowDate)

	monthAgo := now.AddDate(0, 0, -1)

	monthAgo_year	:= monthAgo.Year()
	monthAgo_month	:= int(monthAgo.Month())
	monthAgo_day	:= monthAgo.Day()
	monthAgo_hour	:= monthAgo.Hour()

	monthAgoDate 	:= strconv.Itoa(monthAgo_year) + appLibs.StrPad(strconv.Itoa(monthAgo_month), 2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(monthAgo_day),2, "0","LEFT") + appLibs.StrPad(strconv.Itoa(monthAgo_hour),2, "0","LEFT")

	fmt.Println(monthAgoDate)

	tx := c.Get("Tx").(*gorm.DB)

	chartLabelIdValueDatas := []ChartLabelIdValueData{}

	tx.Raw(`
		SELECT COUNT(*) value,
		      content_id 'id',
		      'VOD' label
		FROM TB_VCMS_PLAY_LOG
		WHERE cclientid NOT IN (
		  SELECT cclientid
		  FROM TB_VCMS_ACCESS_LOG
		  WHERE xevent = 'disconnect'
		)
		  AND tenant_id = ?
		  AND content_id != '-'
		  AND group_year_month_day_hour between ? and ?
		UNION
		SELECT COUNT(*) value,
		  channel_id ,
		  concat
		  (
		      (
			SELECT name
			FROM TB_CHANNEL
			WHERE channel_id = a.channel_id
		      )
		  , " (채널)") channel_name
		  FROM TB_VCMS_PLAY_LOG a
		  WHERE cclientid NOT IN (
		    SELECT cclientid
		    FROM TB_VCMS_ACCESS_LOG
		    WHERE xevent = 'disconnect'
		  )
		  AND tenant_id = ?
		  AND channel_id != '-'
		  AND group_year_month_day_hour between ? and ?
		  GROUP BY channel_id
		`, claims.TenantId, monthAgoDate, nowDate, claims.TenantId, monthAgoDate, nowDate).Scan(&chartLabelIdValueDatas) // (*sql.Rows, error)


	var chart_data string

	chart_data += "["


	for i, val := range chartLabelIdValueDatas {


		chart_data += "["


		chart_data += `"` + val.Label + `",`
		chart_data += fmt.Sprintf("%.0f", val.Value) + ""

		chart_data += "]"


		if i != len(chartLabelIdValueDatas)-1 {
			chart_data += ","
		}

	}

	chart_data += "]"

	fmt.Println(chart_data)


	return handler.APIResultHandler(c, true, http.StatusCreated,
		chart_data)

}


// 등록된 컨텐츠 수 -- 해결
func (api DashBoardAPI) GetContentsCount(c echo.Context) error {

	claims 	 := apiModel.GetJWTClaims(c)

	contents := []rdbModel.Content{}
	tx := c.Get("Tx").(*gorm.DB)

	var count = 0

	tx.Where("tenant_id = ? ", claims.TenantId).Find(&contents).Count(&count)


	return handler.APIResultHandler(c, true, http.StatusOK, count)
}

// 스토리지 사용량 -- 해결
func (api DashBoardAPI) GetStorageSize(c echo.Context) error {

	claims 	 := apiModel.GetJWTClaims(c)

	storage_path := appConfig.Config.APP.ContentsRoot + string(filepath.Separator) + "vod" + string(filepath.Separator) + claims.TenantId

	storage_size := appLibs.DirSizeMB(storage_path)

	return handler.APIResultHandler(c, true, http.StatusOK, storage_size)
}


// 노출 많이된 컨텐츠 (top 10) -- 해결
type ContentsRanking struct {
	ContentCount 	int64 		`json:"content_count"`
	ContentId	string 		`json:"content_id"`
	ContentName 	string 		`json:"content_name"`
}

func (api DashBoardAPI) GetContentsRankingChart(c echo.Context) error {

	claims 	 := apiModel.GetJWTClaims(c)
	tx := c.Get("Tx").(*gorm.DB)

	contentsRankings := []ContentsRanking{}

	tx.Raw(`
	SELECT y.content_count, y.content_id, y.content_name
		FROM
		(
		  SELECT COUNT(*) content_count, x.content_id, x.content_name
		  FROM
		  (
		    SELECT a.content_id,
		    (
		      SELECT content_name FROM TB_CONTENT
		      WHERE content_id = a.content_id
		    ) as content_name
		    FROM TB_VCMS_PLAY_LOG a
		    WHERE tenant_id = ?
		    AND content_id != '-'
		  ) x
		  GROUP BY x.content_id, x.content_name
		) y
		ORDER BY content_count DESC
		LIMIT 5
		`, claims.TenantId).Scan(&contentsRankings) // (*sql.Rows, error)

	var chart_data string
	for i, val := range contentsRankings {


		chart_data += "["


		chart_data += `"` + val.ContentName + `",`
		chart_data += fmt.Sprintf("%d", val.ContentCount) + ""

		chart_data += "]"


		if i != len(contentsRankings)-1 {
			chart_data += ","
		}

	}

	return handler.APIResultHandler(c, true, http.StatusOK, chart_data)
}
//-------------------------------------------------------
