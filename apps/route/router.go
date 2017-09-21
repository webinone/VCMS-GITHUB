package route

import (
	"github.com/labstack/echo"
	echoMw "github.com/labstack/echo/middleware"
	"VCMS/apps/api"
	apiModel "VCMS/apps/models/api"
	"gopkg.in/go-playground/validator.v9"
	appConfig "VCMS/apps/config"
	"VCMS/apps/handler"
	"VCMS/apps/db"
	front_api "VCMS/apps/api/front"
	"net/http"
)

func Init() *echo.Echo {

	e := echo.New()
	// validator 등록
	e.Validator = &apiModel.APIValidator{validator.New()}

	e.Use(echoMw.Logger())
	e.Use(echoMw.Gzip())

	e.Static("/thumbnail", "statics/thumbnail")
	e.Static("/banner",    "statics/banner")
	e.Static("/download",  "statics/download")
	e.Static("/excel",  "statics/excel")

	// CORS 설정
	e.Use(echoMw.CORSWithConfig(echoMw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderAuthorization,
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAcceptEncoding,
			"x-request-with",
			"responseType",
		},
	}))

	// JWT 설정
	jwtMw := echoMw.JWTWithConfig(echoMw.JWTConfig{
		SigningKey: []byte(appConfig.Config.AUTH.JwtKey),
		ContextKey: "jwt",
		Claims:&apiModel.JWTClaims{},
		//TokenLookup: "header:token" ,
	})

	e.HTTPErrorHandler = handler.JSONHTTPErrorHandler

	e.Use(handler.TransactionHandler(db.Init()))

	e.GET("/api/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, VCMS Rest API Server!")
	})


	// Routing 설정 밑으로 개발자들은 밑에만 신경 쓰면 된다.
	//-----------------------------------------------------------------

	e.GET("/api/webinar/notice/download/:idx", 	api.WebinarNoticeAPI{}.DownloadNoticeFile)
	e.GET("/api/webinar/site/download/:idx", 		api.WebinarSiteAPI{}.DownloadSiteFile)
	e.GET("/api/webinar/qna/download", 			api.WebinarAdminQnAAPI{}.DownloadExcelFile)

	// 참여자 엑셀 다운로드 정보
	e.GET("/api/webinar/join/download", 	api.WebinarAdminJoinAPI{}.DownloadExcelFile)

	// 설문조사 엑셀 다운로드 정보
	e.GET("/api/webinar/poll/download", 	api.WebinarPollMemberAPI{}.DownloadExcelFile)

	// Login
	e.POST("/api/login", api.AuthAPI{}.PostLogin)

	// JWT 체크 부분
	v1 := e.Group("/api/v1", jwtMw)
	{
		// 테넌트
		//--------------------------------------------------------------
		v1.POST("/tenant", 		api.TenantAPI{}.PostTenant)
		v1.PUT("/tenant/:idx", 		api.TenantAPI{}.PutTenant)
		v1.DELETE("/tenant/:idx", 		api.TenantAPI{}.DeleteTenant)
		v1.GET("/tenant/:idx", 		api.TenantAPI{}.GetTenant)
		v1.GET("/tenant", 			api.TenantAPI{}.GetTenants)
		//--------------------------------------------------------------

		// 사용자
		//--------------------------------------------------------------
		v1.POST("/user", 			api.UserAPI{}.PostUser)
		v1.GET("/user/:idx", 		api.UserAPI{}.GetUserByIdx)
		v1.GET("/user/user_id/:user_id", 	api.UserAPI{}.GetUserByUserId)
		v1.GET("/user", 			api.UserAPI{}.GetUsers)
		v1.PUT("/user/:idx", 		api.UserAPI{}.PutUser)
		v1.DELETE("/user/:idx", 	api.UserAPI{}.DeleteUser)

		v1.GET("/user/hellot", 	api.UserAPI{}.GetHelloTUsers)
		//--------------------------------------------------------------

		// 카테고리
		//--------------------------------------------------------------
		v1.GET("/category",		api.CategoryAPI{}.GetCategories)
		v1.GET("/category/:idx",		api.CategoryAPI{}.GetCategory)
		v1.GET("/category/paths",		api.CategoryAPI{}.GetCategoryPaths)
		v1.POST("/category",		api.CategoryAPI{}.PostCategory)
		v1.PUT("/category/:idx",		api.CategoryAPI{}.PutCategory)
		v1.DELETE("/category/:idx",	api.CategoryAPI{}.DeleteCategory)
		//--------------------------------------------------------------

		// 컨텐츠
		//--------------------------------------------------------------
		v1.POST("/content/upload", 	api.ContentAPI{}.UploadContent)
		v1.POST("/content", 		api.ContentAPI{}.PostContent)
		v1.GET("/content", 		api.ContentAPI{}.GetContents)
		v1.GET("/content/:idx", 		api.ContentAPI{}.GetContent)
		v1.DELETE("/content/:idx", 	api.ContentAPI{}.DeleteContent)

		v1.POST("/content/thumbnail", 	api.ContentAPI{}.UploadThumbnail)
		v1.PUT("/content/:idx", 		api.ContentAPI{}.PutContent)
		//------------------------------------------------------------

		// 채널 관리
		//----------------------------------------------------------------------
		v1.POST("/channel", 		api.ChannelAPI{}.PostChannel)
		v1.GET("/channel", 		api.ChannelAPI{}.GetChannels)
		v1.GET("/channel/:idx", 		api.ChannelAPI{}.GetChannel)
		v1.DELETE("/channel/:idx", 	api.ChannelAPI{}.DeleteChannel)
		v1.PUT("/channel/:idx", 		api.ChannelAPI{}.PutChannel)
		//-------------------------------------------------------------------------

		// 스케쥴 관리
		//----------------------------------------------------------------------
		v1.POST("/schedule", 		api.ScheduleAPI{}.PostSchedule)
		v1.GET("/schedule", 		api.ScheduleAPI{}.GetSchedules)
		v1.GET("/schedule/:idx", 		api.ScheduleAPI{}.GetSchedule)
		v1.DELETE("/schedule/:idx", 	api.ScheduleAPI{}.DeleteSchedule)
		v1.PUT("/schedule/:idx", 		api.ScheduleAPI{}.PutSchedule)
		//-------------------------------------------------------------------------

		// 대시보드
		//----------------------------------------------------------------------
		v1.GET("/dashboard/monthconnect", 			api.DashBoardAPI{}.GetMonthConnections)
		v1.GET("/dashboard/weeksconnect", 			api.DashBoardAPI{}.GetWeeksConnections)
		v1.GET("/dashboard/todayconnect", 			api.DashBoardAPI{}.GetTodayConnections)
		v1.GET("/dashboard/contentcount", 			api.DashBoardAPI{}.GetContentsCount)
		v1.GET("/dashboard/channelcount", 			api.DashBoardAPI{}.GetChannelsCount)
		v1.GET("/dashboard/contentsrankingchart", 		api.DashBoardAPI{}.GetContentsRankingChart)
		v1.GET("/dashboard/schedules", 			api.DashBoardAPI{}.GetSchedules)
		v1.GET("/dashboard/storagesize", 			api.DashBoardAPI{}.GetStorageSize)
		v1.GET("/dashboard/appconnect", 			api.DashBoardAPI{}.GetApplicationConnectionCount)
		v1.GET("/dashboard/channelconnnect/:channel_id", 	api.DashBoardAPI{}.GetChannelConnectionCount)
		v1.GET("/dashboard/nowconnectpiechart", 		api.DashBoardAPI{}.GetNowConnectPieChart)

		//----------------------------------------------------------------------

		// Webinar
		//----------------------------------------------------------------------------------

		webinar := v1.Group("/webinar")
		{
			// SITE
			//----------------------------------------------------------------------------
			webinar.GET("/site",				api.WebinarSiteAPI{}.GetWebinarSites)
			webinar.GET("/site/:idx",			api.WebinarSiteAPI{}.GetWebinarSite)
			webinar.POST("/site/thumbnail", 	api.WebinarSiteAPI{}.UploadThumbNail)
			webinar.POST("/site/file", 		api.WebinarSiteAPI{}.UploadSiteFile)
			webinar.POST("/site/backimage", 	api.WebinarSiteAPI{}.UploadBackImage)
			webinar.POST("/site",				api.WebinarSiteAPI{}.PostWebinarSite)
			webinar.DELETE("/site/:idx",		api.WebinarSiteAPI{}.DeleteWebinarSite)
			webinar.PUT("/site/:idx",			api.WebinarSiteAPI{}.PutWebinarSite)
			//----------------------------------------------------------------------------

			// BANNER
			//----------------------------------------------------------------------------
			webinar.POST("/banner", 			api.WebinarBannerAPI{}.PostWebinarBanner)
			webinar.GET("/banner", 			api.WebinarBannerAPI{}.GetWebinarBanners)
			webinar.POST("/banner/upload", 		api.WebinarBannerAPI{}.UploadBanner)
			webinar.PUT("/banner/:banner_type", 	api.WebinarBannerAPI{}.PutWebinarBanner)

			//----------------------------------------------------------------------------

			// 공지사항
			//----------------------------------------------------------------------------
			webinar.POST("/notice/upload", 		api.WebinarNoticeAPI{}.UploadNoticeFile)
			webinar.GET("/notice/download/:idx", 	api.WebinarNoticeAPI{}.DownloadNoticeFile)
			webinar.POST("/notice", 			api.WebinarNoticeAPI{}.PostWebinarNotice)
			webinar.PUT("/notice/:idx", 		api.WebinarNoticeAPI{}.PutWebinarNotice)
			webinar.GET("/notice", 			api.WebinarNoticeAPI{}.GetWebinarNotices)
			webinar.GET("/notice/:idx", 		api.WebinarNoticeAPI{}.GetWebinarNotice)

			webinar.DELETE("/notice/:idx",		api.WebinarNoticeAPI{}.DeleteWebinarNotice)
			//----------------------------------------------------------------------------

			// Q&A
			//----------------------------------------------------------------------------
			webinar.POST("/qna", 			api.WebinarAdminQnAAPI{}.PostWebinarQnA)
			webinar.DELETE("/qna/:idx", 		api.WebinarAdminQnAAPI{}.DeleteWebinarQnA)
			webinar.GET("/qna", 			api.WebinarAdminQnAAPI{}.GetWebinarQnAs)
			webinar.GET("/qna/:idx", 			api.WebinarAdminQnAAPI{}.GetWebinarQnA)

			//-----------------------------------------------------------------------------

			// JOIN
			//----------------------------------------------------------------------------
			webinar.GET("/join", 			api.WebinarAdminJoinAPI{}.GetWebinarJoins)
			//----------------------------------------------------------------------------

			// 응원 댓글
			//-----------------------------------------------------------------------
			webinar.GET("/comment", 			api.WebinarCommentAPI{}.GetWebinarComments)
			webinar.DELETE("/comment/:idx", 		api.WebinarCommentAPI{}.DeleteWebinarComment)

			//-----------------------------------------------------------------------

			// 설문
			//-----------------------------------------------------------------------
			webinar.POST("/poll", 							api.WebinarPollAPI{}.PostWebinarPoll)
			webinar.GET("/poll", 							api.WebinarPollAPI{}.GetWebinarPolls)
			webinar.GET("/poll/:idx", 						api.WebinarPollAPI{}.GetWebinarPoll)
			webinar.DELETE("/poll/:idx",					api.WebinarPollAPI{}.DeleteWebinarPoll)
			webinar.PUT("/poll/:idx",						api.WebinarPollAPI{}.PutWebinarPoll)

			webinar.POST("/poll/question",					api.WebinarPollQuestionAPI{}.PostWebinarPollQuestion)
			webinar.GET("/poll/question",					api.WebinarPollQuestionAPI{}.GetWebinarPollQuestions)
			webinar.GET("/poll/question/:idx",				api.WebinarPollQuestionAPI{}.GetWebinarPollQuestion)

			webinar.GET("/poll/member",					api.WebinarPollMemberAPI{}.GetWebinarPollMembers)
			webinar.GET("/poll/member/statistics",			api.WebinarPollMemberAPI{}.GetWebinarPollMemberStatistics)
			webinar.PUT("/poll/member/win",				api.WebinarPollMemberAPI{}.PutWebinarPollMemberWinYN)

			//-----------------------------------------------------------------------

		}
	}

	frontJwtMw := echoMw.JWTWithConfig(echoMw.JWTConfig{
		SigningKey: []byte(appConfig.Config.AUTH.JwtKey),
		ContextKey: "jwt-front",
		Claims:&apiModel.FrontJWTClaims{},
		//TokenLookup: "header:token" ,
	})

	// 권한 없는 영역
	// ---------------------------------------------------------------------------------------
	e.POST("/api/front/webinar/newsletter", 		front_api.NewsLetterAPI{}.PostNewsLetter)
	e.GET("/api/front/webinar/news", 				front_api.WebinarNewsAPI{}.GetWebinarNews)
	e.GET("/api/webinar/banner/redirect",			front_api.WebinarBannerAPI{}.GetWebinarBannerRedirect)
	e.GET("/api/front/webinar/site",				front_api.WebinarSiteAPI{}.GetWebinarIfSites)
	e.GET("/api/front/webinar/site/relations",		front_api.WebinarSiteAPI{}.GetWebinarSitesByCompanyNo)
	e.POST("/api/front/webinar/adduser",			front_api.AuthAPI{}.PostAddUser)

	// 아카데미 연동 웨비나 리스트
	e.GET("/api/front/webinar/site/academy",	front_api.WebinarSiteAPI{}.GetWebinarSites)

	// Webinar site MAIN 관련
	//-----------------------------------------------------------------------
	e.GET("/api/front/webinar/site/:webinar_site_id",		front_api.WebinarSiteAPI{}.GetWebinarSite)

	// NOTICE
	//-----------------------------------------------------------------------
	e.GET("/api/front/webinar/notice",		front_api.WebinarNoticeAPI{}.GetWebinarNotices)
	//-----------------------------------------------------------------------

	// Banner
	//-----------------------------------------------------------------------
	e.GET("/api/front/webinar/banner",			front_api.WebinarBannerAPI{}.GetWebinarBanners)

	//-----------------------------------------------------------------------

	e.GET("/api/front/webinar/comment", 		front_api.WebinarCommentAPI{}.GetWebinarComments)

	//-----------------------------------------------------------------------

	// front (사용자) API
	frontApi := e.Group("/api/front/webinar", frontJwtMw)
	{
		// frontAPI
		//-----------------------------------------------------------------------
		frontApi.POST("/auth/check", 	front_api.AuthAPI{}.PostSessionCheck)
		//-----------------------------------------------------------------------

		// Q&A
		//-----------------------------------------------------------------------
		frontApi.POST("/qna", 					front_api.WebinarFrontQnAAPI{}.PostWebinarQnA)
		frontApi.PUT("/qna/:idx", 				front_api.WebinarFrontQnAAPI{}.PutWebinarQnA)
		frontApi.DELETE("/qna/:idx", 			front_api.WebinarFrontQnAAPI{}.DeleteWebinarQnA)
		frontApi.GET("/qna", 					front_api.WebinarFrontQnAAPI{}.GetWebinarQnAs)
		frontApi.POST("/qna/reply", 			front_api.WebinarFrontQnAAPI{}.PostWebinarAdminQnA)
		frontApi.DELETE("/qna/reply/:idx", 	front_api.WebinarFrontQnAAPI{}.DeleteWebinarAdminQnA)
		//-----------------------------------------------------------------------

		// JOIN
		//-----------------------------------------------------------------------
		frontApi.POST("/join", 				front_api.WebinarFrontJoinAPI{}.PostWebinarJoin)
		frontApi.GET("/join/check", 			front_api.WebinarFrontJoinAPI{}.GetWebinarJoinCheck)

		//-----------------------------------------------------------------------

		// 응원 댓글
		//-----------------------------------------------------------------------
		frontApi.POST("/comment", 				front_api.WebinarCommentAPI{}.PostWebinarComment)

		frontApi.PUT("/comment/:idx", 			front_api.WebinarCommentAPI{}.PutWebinarComment)
		frontApi.DELETE("/comment/:idx", 		front_api.WebinarCommentAPI{}.DeleteWebinarComment)

		//-----------------------------------------------------------------------

		// 설문조사
		//-----------------------------------------------------------------------
		frontApi.POST("/poll/member", 			front_api.WebinarPollMemberAPI{}.PostWebinarPollMember)
		frontApi.GET("/poll", 					front_api.WebinarPollAPI{}.GetWebinarPoll)
		//-----------------------------------------------------------------------

	}


	return e
}