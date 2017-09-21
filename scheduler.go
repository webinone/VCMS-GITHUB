package main

import (
	appConfig "VCMS/apps/config"
	rdbModel "VCMS/apps/models/rdb"
	"fmt"
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/metakeule/fmtdate"
)

func init() {
	// config Loading
	appConfig.LoadAutoConfig()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

// 원본 로그 파일 삭제하기 최근 1달 이상 지난 것은 삭제 한다.
func deleteStatsLog() {
	logrus.Debug("deleteStatsLog !!")

	db, err := gorm.Open(appConfig.Config.RDB.Product, appConfig.Config.RDB.ConnectString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(appConfig.Config.RDB.Debug)

	vcms_play_log 	:= &rdbModel.VcmsPlayLog{}
	vcms_access_log := &rdbModel.VcmsAccessLog{}

	offSet, _ := time.ParseDuration("+09.00h")
	now := time.Now().UTC().Add(offSet)

	fmt.Println("Today : ", fmtdate.Format("YYYYMMDD", now))
	monthAgo := now.AddDate(0, -1, 0)

	monthAgoString := fmtdate.Format("YYYYMMDD", monthAgo)
	fmt.Println("Month ago : ", monthAgoString)

	tx := db.Begin()

	tx.Where("group_year_month_day < ? ", monthAgoString).Delete(vcms_play_log)
	tx.Where("group_year_month_day < ? ", monthAgoString).Delete(vcms_access_log)

	tx.Commit()

}

// 원본 플레이 로그 파일을 옮기기...
func moveVCMSPlayLog() {

	logrus.Debug("moveVCMSPlayLog !!")

	db, err := gorm.Open(appConfig.Config.RDB.Product, appConfig.Config.RDB.ConnectString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(appConfig.Config.RDB.Debug)


	tx := db.Begin()

	tx.Exec(`
		INSERT  INTO TB_VCMS_PLAY_LOG (date, time, year, month, day, hour, minute, group_year_month,
						group_year_month_day, group_year_month_day_hour, group_year_month_day_hour_minute, tz,
						tenant_id, channel_id, content_id, xsname,
						xduration, cip, cuseragent, cclientid, xfilename, cproto, xsuri
					)
		      SELECT  a.date, a.time, a.year,
			      a.month, a.day,
			      a.hour, a.minute,
			      (
				CONCAT(a.year, a.month)
			      ) group_year_month,
			      (
				CONCAT(a.year, a.month, a.day)
			      ) group_year_month_day,
			      (
				CONCAT(a.year, a.month, a.day, a.hour)
			      ) group_year_month_day_hour,
			      (
				CONCAT(a.year, a.month, a.day, a.hour, a.minute)
			      ) group_year_month_day_hour_minute,
			      a.tz,
			      a.tenant_id,
			      a.channel_id,
			      a.content_id,
			      a.xsname,
			      a.xduration,
			      a.cip,
			      a.cuseragent,
			      a.cclientid,
			      a.xfilename,
			      a.cproto,
			      a.xsuri
			FROM
			  (
			    SELECT
			      date,
			      time,
			      (
				left(date, 4)
			      ) AS year,
			      (
				substring(date, 6, 2)
			      ) AS month,
			      (
				substring(date, 9, 2)
			      ) AS day,
			      (
				left(time, 2)
			      ) AS hour,
			      (
				substring(time, 4, 2)
			      ) AS minute,
			      tz,
			      (
				CASE
				WHEN xapp != 'vcms_vod'
				  THEN xapp
				ELSE left(xctx, INSTR(xctx, '/') - 1)
				END
			      ) AS tenant_id,
			      (
				CASE
				WHEN xfileext != 'mp4'
				  THEN xsname
				ELSE '-'
				END
			      ) AS channel_id,
			      (
				CASE
				WHEN xfileext = 'mp4'
				  THEN REPLACE(substring(xsname, INSTR(xsname, '/') + 1), '.mp4', '')
				ELSE '-'
				END
			      ) AS content_id,
			      xsname,
			      xduration,
			      cip,
			      cuseragent,
			      cclientid,
			      xfilename,
			      cproto,
			      xsuri
			    FROM TB_WOWZA_ACCESSLOG
			    WHERE xevent = 'play'
				  AND xapp != 'vod'
			  ) a
	`)

	tx.Exec(`
		DELETE FROM TB_WOWZA_ACCESSLOG
		WHERE xevent = 'play'
	`)

	tx.Commit()

}

// 원본 플레이 로그 파일을 옮기기...
func moveVCMSAccessLog() {

	logrus.Debug("moveVCMSAccessLog !!")

	db, err := gorm.Open(appConfig.Config.RDB.Product, appConfig.Config.RDB.ConnectString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(appConfig.Config.RDB.Debug)

	//if db.HasTable(&rdbModel.Tenant{}) {
	//	db.DropTable(&rdbModel.VcmsLog{})
	//}
	//
	//db.AutoMigrate(&rdbModel.VcmsLog{})

	tx := db.Begin()

	tx.Exec(`
		INSERT  INTO TB_VCMS_ACCESS_LOG (date, time, year, month, day, hour, minute, group_year_month,
						group_year_month_day, group_year_month_day_hour, group_year_month_day_hour_minute, tz,
						tenant_id, channel_id, content_id, xevent, xsname,
						xduration, cip, cuseragent, cclientid, xfilename, cproto, xsuri
					)
		      SELECT  a.date, a.time, a.year,
			      a.month, a.day,
			      a.hour, a.minute,
			      (
				CONCAT(a.year, a.month)
			      ) group_year_month,
			      (
				CONCAT(a.year, a.month, a.day)
			      ) group_year_month_day,
			      (
				CONCAT(a.year, a.month, a.day, a.hour)
			      ) group_year_month_day_hour,
			      (
				CONCAT(a.year, a.month, a.day, a.hour, a.minute)
			      ) group_year_month_day_hour_minute,
			      a.tz,
			      a.tenant_id,
			      a.channel_id,
			      a.content_id,
			      a.xevent,
			      a.xsname,
			      a.xduration,
			      a.cip,
			      a.cuseragent,
			      a.cclientid,
			      a.xfilename,
			      a.cproto,
			      a.xsuri
			FROM
			  (
			    SELECT
			      date,
			      time,
			      (
				left(date, 4)
			      ) AS year,
			      (
				substring(date, 6, 2)
			      ) AS month,
			      (
				substring(date, 9, 2)
			      ) AS day,
			      (
				left(time, 2)
			      ) AS hour,
			      (
				substring(time, 4, 2)
			      ) AS minute,
			      tz,
			      (
				CASE
				WHEN xapp != 'vcms_vod'
				  THEN xapp
				ELSE left(xctx, INSTR(xctx, '/') - 1)
				END
			      ) AS tenant_id,
			      (
				CASE
				WHEN xfileext != 'mp4'
				  THEN xsname
				ELSE '-'
				END
			      ) AS channel_id,
			      (
				CASE
				WHEN xfileext = 'mp4'
				  THEN REPLACE(substring(xsname, INSTR(xsname, '/') + 1), '.mp4', '')
				ELSE '-'
				END
			      ) AS content_id,
			      xevent,
			      xsname,
			      xduration,
			      cip,
			      cuseragent,
			      cclientid,
			      xfilename,
			      cproto,
			      xsuri
			    FROM TB_WOWZA_ACCESSLOG
			    WHERE xevent != 'play'
				  AND xapp != 'vod'
			  ) a
	`)

	tx.Exec(`
		DELETE FROM TB_WOWZA_ACCESSLOG
		WHERE xevent = 'connect' OR xevent = 'disconnect'
	`)

	tx.Commit()

}

func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}

func main() {
	// Do jobs with params
	//gocron.Every(1).Second().Do(taskWithParams, 1, "hello")

	// Do jobs without params
	//gocron.Every(1).Second().Do(deleteAccessLog)
	// 로그 데이터 옮기기
	gocron.Every(1).Minute().Do(moveVCMSPlayLog)
	gocron.Every(1).Minute().Do(moveVCMSAccessLog)
	//gocron.Every(1).Second().Do(jsonTestExecute)
	//gocron.Every(2).Seconds().Do(task)
	//gocron.Every(1).Minute().Do(task)
	//gocron.Every(2).Minutes().Do(task)
	//gocron.Every(1).Hour().Do(task)
	//gocron.Every(2).Hours().Do(task)
	//gocron.Every(1).Day().Do(task)
	//gocron.Every(2).Days().Do(task)
	//
	//// Do jobs on specific weekday
	//gocron.Every(1).Monday().Do(task)
	//gocron.Every(1).Thursday().Do(task)
	//
	//// function At() take a string like 'hour:min'
	// 01:00 하루 한번 삭제
	gocron.Every(1).Day().At("01:00").Do(deleteStatsLog)
	//gocron.Every(1).Monday().At("18:30").Do(task)
	//
	//// remove, clear and next_run
	//_, time := gocron.NextRun()
	//fmt.Println(time)
	//
	//gocron.Remove(task)
	//gocron.Clear()

	// function Start start all the pending jobs
	<-gocron.Start()

	// also , you can create a your new scheduler,
	// to run two scheduler concurrently
	//s := gocron.NewScheduler()
	//s.Every(3).Seconds().Do(task)
	//<- s.Start()

}
