package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	appConfig "VCMS/apps/config"
	rdbModel "VCMS/apps/models/rdb"
)

func Init() *gorm.DB {

	// Connection String
	rdb, err := gorm.Open(appConfig.Config.RDB[0].Product, appConfig.Config.RDB[0].ConnectString)
	if err != nil {
		panic(err)
	}

	rdb.DB().Ping()
	rdb.DB().SetMaxIdleConns(appConfig.Config.RDB[0].MaxIdleConns)
	rdb.DB().SetMaxOpenConns(appConfig.Config.RDB[0].MaxOpenConns)

	rdb.LogMode(appConfig.Config.RDB[0].Debug)

	if rdb.HasTable(&rdbModel.Tenant{}) {
		//rdb.DropTable(&rdbModel.Tenant{})
		//rdb.DropTable(&rdbModel.User{})
		//rdb.DropTable(&rdbModel.Category{})
		//rdb.DropTable(&rdbModel.Content{})
		//rdb.DropTable(&rdbModel.ContentTag{})
		//rdb.DropTable(&rdbModel.Stream{})
		//rdb.DropTable(&rdbModel.ThumbNail{})
		//rdb.DropTable(&rdbModel.Channel{})
		//rdb.DropTable(&rdbModel.Schedule{})
		//rdb.DropTable(&rdbModel.ScheduleOrder{})
		//rdb.DropTable(&rdbModel.WebinarSite{})
		//rdb.DropTable(&rdbModel.WebinarSiteTag{})
		//rdb.DropTable(&rdbModel.WebinarBanner{})
		//rdb.DropTable(&rdbModel.WebinarNotice{})
		//rdb.DropTable(&rdbModel.WebinarNoticeFile{})
		//rdb.DropTable(&rdbModel.WebinarPollQuestionMaster{})
		//rdb.DropTable(&rdbModel.WebinarPollQuestionDetail{})
		//rdb.DropTable(&rdbModel.WebinarPollMember{})
		//rdb.DropTable(&rdbModel.WebinarPollMemberResult{})
		//rdb.DropTable(&rdbModel.WebinarJoin{})
		//rdb.DropTable(&rdbModel.WebinarJoinMember{})

		//rdb.DropTable(&rdbModel.WebinarSiteBackImage{})
		//rdb.DropTable(&rdbModel.WebinarSiteAdmin{})
	}

	//rdb.AutoMigrate(&rdbModel.Tenant{})
	//rdb.AutoMigrate(&rdbModel.User{})
	//rdb.AutoMigrate(&rdbModel.Category{})
	//rdb.AutoMigrate(&rdbModel.Content{})
	//rdb.AutoMigrate(&rdbModel.ContentTag{})
	//rdb.AutoMigrate(&rdbModel.Stream{})
	//rdb.AutoMigrate(&rdbModel.ThumbNail{})
	//rdb.AutoMigrate(&rdbModel.Channel{})
	//rdb.AutoMigrate(&rdbModel.Schedule{})
	//rdb.AutoMigrate(&rdbModel.ScheduleOrder{})
	//rdb.AutoMigrate(&rdbModel.WebinarSite{})
	//rdb.AutoMigrate(&rdbModel.WebinarSiteTag{})
	//rdb.AutoMigrate(&rdbModel.WebinarBanner{})
	//
	//rdb.AutoMigrate(&rdbModel.WebinarNotice{})
	//rdb.AutoMigrate(&rdbModel.WebinarNoticeFile{})
	//rdb.AutoMigrate(&rdbModel.WebinarFrontQnA{})
	//rdb.AutoMigrate(&rdbModel.WebinarAdminQnA{})
	//rdb.AutoMigrate(&rdbModel.WebinarBanner{})

	//rdb.AutoMigrate(&rdbModel.MemberDefault{})
	//rdb.AutoMigrate(&rdbModel.MemberSub{})
	//rdb.AutoMigrate(&rdbModel.WebinarPoll{})
	//rdb.AutoMigrate(&rdbModel.WebinarPollMember{})
	//rdb.AutoMigrate(&rdbModel.WebinarPollMemberResult{})
	//rdb.AutoMigrate(&rdbModel.WebinarPollQuestionMaster{})
	//rdb.AutoMigrate(&rdbModel.WebinarPollQuestionDetail{})
	//rdb.AutoMigrate(&rdbModel.WebinarJoin{})
	//rdb.AutoMigrate(&rdbModel.WebinarJoinMember{})

	//rdb.AutoMigrate(&rdbModel.VcmsAccessLog{})
	//rdb.AutoMigrate(&rdbModel.VcmsPlayLog{})
	//rdb.AutoMigrate(&rdbModel.WowzaAccesslog{})
	//rdb.AutoMigrate(&rdbModel.WebinarComment{})

	//rdb.AutoMigrate(&rdbModel.WebinarSiteFile{})
	//rdb.AutoMigrate(&rdbModel.WebinarSiteThumbNail{})

	//rdb.AutoMigrate(&rdbModel.WebinarSiteBackImage{})
	//rdb.AutoMigrate(&rdbModel.WebinarSiteAdmin{})

	return rdb

}

// HelloT DB Connection
func GetHelloTDB() *gorm.DB {

	rdb, err := gorm.Open(appConfig.Config.RDB[1].Product, appConfig.Config.RDB[1].ConnectString)
	if err != nil {
		panic(err)
	}

	rdb.DB().Ping()
	//rdb.DB().SetMaxIdleConns(appConfig.Config.RDB[1].MaxIdleConns)
	//rdb.DB().SetMaxOpenConns(appConfig.Config.RDB[1].MaxOpenConns)

	rdb.LogMode(appConfig.Config.RDB[1].Debug)

	return rdb
}