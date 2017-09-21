package handler

import (
	"github.com/labstack/echo"
	apiModel "VCMS/apps/models/api"
	"github.com/valyala/fasthttp"
	"github.com/jinzhu/gorm"
	"github.com/Sirupsen/logrus"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"fmt"
	echoSession "github.com/ipfans/echo-session"
	"net/http"
	"io/ioutil"
	"bytes"
)

func APIResultHandler(c echo.Context, httpSuccess bool,  httpStatus int, data interface{}) error {

	apiResult := &apiModel.APIResult{
		Success : httpSuccess,
		ResultCode : httpStatus,
		ResultData: data,
	}

	return c.JSON(httpStatus, apiResult)
}

func JSONHTTPErrorHandler(err error, c echo.Context) {
	code := fasthttp.StatusInternalServerError
	msg := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}

	APIResultHandler(c, false, code, msg)
}

// transaction middleware
func TransactionHandler(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {

			tx := db.Begin()

			c.Set("Tx", tx)

			if err := next(c); err != nil {
				tx.Rollback()
				logrus.Debug("Transction Rollback: ", err)
				return err
			}
			logrus.Debug("Transaction Commit")
			tx.Commit()

			return nil
		})
	}
}

// 세션 체크
// 세션이 존재하지 않으면...
func SessionHandler () echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			session := echoSession.Default(ctx)

			if session != nil && session.Get("session_id") != nil {

				session_id 	:= session.Get("session_id").(string)
				token 		:= session.Get("token").(string)

				url := "http://localhost:8081/api/v1/sso/auth/active"

				reqBody := []byte(`
					{
						"session_id" : "`+session_id+`",
						"token" : "`+token+`"
					}
				`)

				var jsonStr = []byte(reqBody)

				fmt.Println(string(jsonStr))
				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
				req.Header.Set("Content-Type","application/json")
				req.Header.Set("Accept","application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {

					fmt.Println(err.Error())

				}
				defer resp.Body.Close()

				fmt.Println("response Status:", resp.Status)
				fmt.Println("response Headers:", resp.Header)
				respBody, _ := ioutil.ReadAll(resp.Body)
				fmt.Println("response Body:", string(respBody))

				return h(ctx)

				//return true
			} else {
				// 로그인 페이지로 이동


				return echo.NewHTTPError(http.StatusUnauthorized, "SESSION_EXPIRED")


			}
		}
	}
}

func SocketIOHandler (c echo.Context) error {

	type RequestMessage struct {
		WebinarSiteId  string	`json:"webinar_site_id"`
	}

	type ResultMessage struct {
		TotalJoinUser string	`json:"total_join_user"`
	}

	//create
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel, arg interface{}) {
		fmt.Println("New client connected")

		fmt.Println(arg)
	})


	//handle custom event
	server.On("send", func(c *gosocketio.Channel, msg RequestMessage) string {
		//send event to all in room

		// WOWZA 에게 물어본다. 현재 접속자 수를

		resultMessage := ResultMessage {
			TotalJoinUser: "100",
		}
		c.BroadcastTo(msg.WebinarSiteId, "message", resultMessage)
		return "OK"
	})

	server.On("join", func(c *gosocketio.Channel, msg RequestMessage) string {
		//send event to all in room
		c.Join(msg.WebinarSiteId)

		return "JOIN OK"
	})


	return nil

}


//func setUpRequest ( c echo.HandlerFunc ) echo.HandlerFunc {
//	return func(ctx echo.Context) error {
//		req := ctx.Request()
//
//		Logger := logger.NewLogger()
//		// add some default fields to the logger ~ on all messages
//		logger := api.log.WithFields(logrus.Fields{
//			"method":     req.Method(),
//			"path":       req.URL().Path(),
//			"request_id": uuid.NewRandom().String(),
//		})
//		ctx.Set(loggerKey, logger)
//		startTime := time.Now()
//
//		defer func() {
//			rsp := ctx.Response()
//			// at the end we will want to log a few more interesting fields
//			logger.WithFields(logrus.Fields{
//				"status_code":  rsp.Status(),
//				"runtime_nano": time.Since(startTime).Nanoseconds(),
//			}).Info("Finished request")
//		}()
//
//		// now we will log out that we have actually started the request
//		logger.WithFields(logrus.Fields{
//			"user_agent":     req.UserAgent(),
//			"content_length": req.ContentLength(),
//		}).Info("Starting request")
//
//		err := f(ctx)
//		if err != nil {
//			ctx.Error(err)
//		}
//
//		return err
//	}
//}
