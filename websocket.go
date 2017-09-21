package main

import (
	//"os"
	//"path/filepath"
	"fmt"
	//"strings"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/rs/cors"
	"net/http"
	"log"
	appLibs "VCMS/apps/libs"
	appConfig "SSO/apps/config"
	"github.com/Sirupsen/logrus"
	//"time"
	"strings"
	"time"
)

type RequestJoinMessage struct {
	WebinarSiteId  	string	`json:"webinar_site_id"`
	TenantId	   	string	`json:"tenant_id"`
	ChannelId		string	`json:"channel_id"`
}

type ResultJoinMessage struct {
	TotalConnections int	`json:"total_connections"`
}

// 관리자 답변 Request.......
type RequestQnAMessage struct {
	WebinarSiteId  	string	`json:"webinar_site_id"`
	MemberId	   	string	`json:"member_id"`
	MemberName		string	`json:"member_name"`
	Message			string	`json:"message"`
}

type ResultQnAMessage struct {
	WebinarSiteId  	string		`json:"webinar_site_id"`
	MemberId	   	string		`json:"member_id"`
	MemberName		string		`json:"member_name"`
	Message			string		`json:"message"`
	CreatedAt		time.Time	`json:"created_at"`
}

func init() {
	// config Loading
	appConfig.LoadAutoConfig()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}


func main() {

	//create
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	//server.

	//handle connected
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel, arg interface{}) {
		fmt.Println("New client connected")

		fmt.Println(arg)
	})

	server.On("join", func(c *gosocketio.Channel, msg RequestJoinMessage) int {
		//send event to all in room
		c.Join(msg.WebinarSiteId)

		fmt.Println("JOIN !!!")
		fmt.Println("msg.WebinarSiteId : ", msg.WebinarSiteId)
		fmt.Println("Client ID : ", c.Id())

		c.Join(msg.WebinarSiteId)

		redis_client := appLibs.NewRedisClient()

		redis_key := msg.WebinarSiteId + ":" + c.Id()

		redis_client.Set(redis_key,
			c.Id(),
			-1 )

		result := redis_client.Keys(msg.WebinarSiteId + "*")

		defer redis_client.Close()

		resultMessage := ResultJoinMessage {
			TotalConnections: len(result.Val()),
		}

		c.BroadcastTo(msg.WebinarSiteId, "connection_count", resultMessage)

		return len(result.Val())
	})


	server.On("reply_send_message", func(c *gosocketio.Channel, msg RequestQnAMessage) {

		resultMessage := ResultQnAMessage {
			WebinarSiteId: msg.WebinarSiteId,
			MemberId : msg.MemberId,
			MemberName: msg.MemberName,
			Message: msg.Message,
			CreatedAt:time.Now(),
		}
		c.BroadcastTo(msg.WebinarSiteId, "reply_receive_message", resultMessage)
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		//caller is not necessary, client will be removed from rooms
		//automatically on disconnect
		//but you can remove client from room whenever you need to
		//c.Leave("room name")

		fmt.Println("Disconnected")
		fmt.Println("Client ID : ", c.Id())

		redis_client := appLibs.NewRedisClient()
		defer redis_client.Close()

		redis_key := "*:" + c.Id()

		result := redis_client.Keys(redis_key)

		fmt.Println(result.Name())
		fmt.Println(result.String())
		fmt.Println(result.Result())
		fmt.Println(result.Val())

		results, _ := result.Result()

		if (len(results) > 0) {
			redis_key = results[0]

			redis_client.Del(redis_key)

			webinar_site_id := strings.Split(redis_key, ":")[0]

			c.Leave(webinar_site_id)


			result = redis_client.Keys(webinar_site_id + "*")

			resultMessage := ResultJoinMessage {
				TotalConnections: len(result.Val()),
			}
			c.BroadcastTo(webinar_site_id, "message", resultMessage)

		}

	})

	//setup http server
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	handler := cors.Default().Handler(serveMux)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	handler = c.Handler(handler)

	log.Panic(http.ListenAndServe(":1324", handler))

}

func DirSizeMB(path string) {


}