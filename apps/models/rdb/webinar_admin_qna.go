package rdb

import "time"

type WebinarAdminQnA struct {
	Idx			int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string			`json:"tenant_id"`

	WebinarSiteId		string 			`json:"webinar_site_id"`
	WebinarQnaId		string 			`json:"webinar_qna_id"`

	ReplyContent 		string			`json:"reply_content"`

	EmailSendYN		string			`json:"email_send_yn"` // 1:예, 0:아니오

	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
	UpdatedId		string			`json:"updated_id"`
	DeletedAt 		*time.Time		`json:"-"`
}


func (WebinarAdminQnA) TableName() string {
	return "TB_WEBINAR_ADMIN_QNA"
}
