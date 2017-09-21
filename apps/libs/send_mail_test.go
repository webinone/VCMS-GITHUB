package libs

import (
	"testing"
)

func TestSendMailClient_SendWebinarQnAMail(t *testing.T) {
	send_mail_client := &SendMailClient {
		SmtpURL		: "smtp.mailplug.co.kr",
		SmtpPort	: 465,
		User		: "help@hellot.net",
		Password	: "chomdan4151",
		From		: "help@hellot.net",
		To			: []string{"youngjo.jang@gmaill.com"},
		Subject		: "눼메~ 웨비나 질문 답변이당 씨봉",
	}

	send_mail_client.SendWebinarQnAMail(
		"member_id",
		"김븅신",
		"질문이다 시밤 웨비나 왜 하냐??",
		"2017-10-11 12:30:11",
		"내맘이다 어쩔래 하던말던 시봉아 니가 먼 상관인데",
		"2017-10-11 12:30:11",
		"http://www.naver.com",
		"테스트 웨비나",
		"2017-12-11 12:30:11",
	)

	send_mail_client := &SendMailClient {
		SmtpURL		: "smtp.mailplug.co.kr",
		SmtpPort	: 465,
		User		: "help@hellot.net",
		Password	: "chomdan4151",
		From		: "help@hellot.net",
		To			: []string{"youngjo.jang@gmaill.com"},
		Subject		: "눼메~ 웨비나 질문 답변이당 씨봉",
	}

	send_mail_client.SendWebinarQnAMail(
		"member_id",
		"김븅신",
		"질문이다 시밤 웨비나 왜 하냐??",
		"2017-10-11 12:30:11",
		"내맘이다 어쩔래 하던말던 시봉아 니가 먼 상관인데",
		"2017-10-11 12:30:11",
		"http://www.naver.com",
		"테스트 웨비나",
		"2017-12-11 12:30:11",
	)
}

