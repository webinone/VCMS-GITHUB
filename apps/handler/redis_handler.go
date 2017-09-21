package handler

import (
	"VCMS/apps/libs"
	"fmt"
	"strings"
)

// 사용자 로그인 JWT 토큰 발급
func  StartRedisHandler() error {

	redis_sub_client := libs.NewRedisClient()
	//redis_pub_client := libs.NewSSORedisClient()
	//
	sub := redis_sub_client.PSubscribe("/sso/*")

	for {
		msg, err := sub.ReceiveMessage()
		if err != nil {
			panic(err)
		}

		fmt.Println(msg.Pattern, msg.Channel, msg.Payload)

		fmt.Println(">>> msg.Channel : ", msg.Channel)

		if strings.Index(msg.Channel, "expired") != -1 {



		//logrus.Debug("expired !!!!")
		//logrus.Debug("msg.Payload : ", msg.Payload)
		//
		//session_id := strings.Split(msg.Payload, ":")[1]
		//
		//logrus.Debug("session_id : ", session_id)
		//
		//redis_pub_client.Publish("/sso/expired", session_id)
		}
	}

}
