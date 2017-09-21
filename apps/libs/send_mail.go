package libs

import (
	appConfig "VCMS/apps/config"
	"gopkg.in/gomail.v2"
	"net"
	"crypto/tls"
	"html/template"
	"bytes"
	//"strings"
	"fmt"
)

type SendMailClient struct {
	SmtpURL		string
	SmtpPort	int
	User		string
	Password	string
	From		string
	To			string
	Subject		string
}

func (client SendMailClient) SendWebinarQnAMail(QuestionMemberId 	string,
												QuestionMemberName	string,
												Question 			string,
												QuestionCreatedAt 	string,
												Reply 				string,
												ReplyCreatedAt 		string,
												WebinarSiteUrl 		string,
												WebinarSiteTitle 	string,
												WebinarSiteJoinDate string) error {


	type templateData struct {
		QuestionMemberId 		string
		QuestionMemberName  	string
		Question  				string
		QuestionCreatedAt  		string
		Reply  					string
		ReplyCreatedAt  		string
		WebinarSiteUrl  		string
		WebinarSiteTitle  		string
		WebinarSiteJoinDate  	string
	}

	htmlTemplate := templateData{
		QuestionMemberId : QuestionMemberId,
		QuestionMemberName : QuestionMemberName,
		Question : Question,
		QuestionCreatedAt : QuestionCreatedAt,
		Reply : Reply,
		ReplyCreatedAt : ReplyCreatedAt,
		WebinarSiteUrl : WebinarSiteUrl,
		WebinarSiteTitle : WebinarSiteTitle,
		WebinarSiteJoinDate : WebinarSiteJoinDate,
	}

	htmlPath := appConfig.Config.EMAIL.Template
	bodyHtml, _ := client.ParseTemplate(htmlPath, htmlTemplate)

	fmt.Println(bodyHtml)

	m := gomail.NewMessage()
	m.SetHeader("From", client.From)
	m.SetHeader("To", client.To)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", client.Subject)
	m.SetBody("text/html", bodyHtml)
	//m.Attach("/home/Alex/lolcat.jpg")

	// Connect to the SMTP Server
	servername := client.SmtpURL + ":" + fmt.Sprintf("%v",client.SmtpPort)
	host, _, _ := net.SplitHostPort(servername)

	d := gomail.NewDialer(client.SmtpURL, client.SmtpPort, client.User, client.Password)
	d.SSL = true

	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName: host,
	}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}

func (client SendMailClient) ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return  buf.String(), nil
}