package rdb

import "time"

type WebinarFrontQnA struct {
	Idx			int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string			`json:"tenant_id"`

	WebinarSiteId		string 			`json:"webinar_site_id"`
	WebinarQnaId		string 			`json:"webinar_qna_id"`

	QuestionType		string			`json:"question_type"`	// 1: 웨비나, 2: (1:1문의)

	FrontUserId		string			`json:"front_user_id"`
	FrontUserName		string			`json:"front_user_name"`

	//MemberDefault		MemberDefault		`gorm:"ForeignKey:FrontUserId;AssociationForeignKey:MemberId" json:"front_user_info"`

	WebinarJoin		WebinarJoin		`gorm:"ForeignKey:FrontUserId;AssociationForeignKey:FrontUserId" json:"webinar_join_info"`

	QuestionContent 	string 			`json:"question_content"`
	QuestionVideoTime	string			`json:"question_video_time"`

	WebinarAdminQnA		WebinarAdminQnA 	`gorm:"ForeignKey:WebinarQnaId;AssociationForeignKey:WebinarQnaId" json:"reply"`

	ReplyYN			string			`sql:"DEFAULT:'N'" json:"reply_yn"`	// Y/N

	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
	UpdatedId		string			`json:"updated_id"`
	DeletedAt 		*time.Time		`json:"-"`
}


func (WebinarFrontQnA) TableName() string {
	return "TB_WEBINAR_FRONT_QNA"
}
