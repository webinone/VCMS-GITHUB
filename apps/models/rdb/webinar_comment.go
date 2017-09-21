package rdb

import "time"

type WebinarComment struct {
	Idx			int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string			`json:"tenant_id"`

	WebinarSiteId		string 			`json:"webinar_site_id"`
	WebinarCommentId	string 			`json:"webinar_comment_id"`

	FrontUserId		string			`json:"front_user_id"`
	FrontUserName		string			`json:"front_user_name"`

	MemberDefault		MemberDefault		`gorm:"ForeignKey:FrontUserId;AssociationForeignKey:MemberId" json:"front_user_info"`

	Comment			string			`json:"comment"`

	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
	UpdatedId		string			`json:"updated_id"`
	DeletedAt 		*time.Time		`json:"-"`
}


func (WebinarComment) TableName() string {
	return "TB_WEBINAR_COMMENT"
}
