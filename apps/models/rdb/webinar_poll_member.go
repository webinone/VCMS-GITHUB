package rdb

import "time"

type WebinarPollMember struct {
	Idx			int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string			`json:"tenant_id"`

	WebinarSiteId		string 			`json:"webinar_site_id"`
	WebinarPollId		string 			`json:"webinar_poll_id"`
	WebinarPollMemberId	string			`json:"webinar_poll_member_id"`

	WebinarPollMemberResults	[]WebinarPollMemberResult	`gorm:"ForeignKey:WebinarPollMemberId;AssociationForeignKey:WebinarPollMemberId" json:"webinar_poll_member_result"`

	FrontUserId		string			`json:"front_user_id"`
	FrontUserName		string			`json:"front_user_name"`

	//MemberDefault		MemberDefault		`gorm:"ForeignKey:FrontUserId;AssociationForeignKey:MemberId" json:"-"`

	WinYN			string			`json:"win_yn"`

	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
	UpdatedId		string			`json:"updated_id"`
	DeletedAt 		*time.Time		`json:"-"`
}


func (WebinarPollMember) TableName() string {
	return "TB_WEBINAR_POLL_MEMBER"
}
