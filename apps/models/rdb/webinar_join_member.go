package rdb

import "time"

type WebinarJoinMember struct {
	Idx			int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string			`json:"tenant_id"`

	WebinarSiteId		string 			`json:"webinar_site_id"`

	FrontUserId		string			`json:"front_user_id"`
	FrontUserName		string			`json:"front_user_name"`

	//MemberDefault		MemberDefault		`gorm:"ForeignKey:FrontUserId;AssociationForeignKey:MemberId" json:"front_user_info"`

	MobilePhoneNum		string		`json:"mobile_phone_num"`
	Email			string			`json:"email"`

	CompanyName		string			`json:"company_name"`
	Department		string			`json:"department"`
	Position		string			`json:"position"`
	PhoneNum		string			`json:"phone_num"`
	ZipCode			string			`json:"zip_code"`
	Address1		string			`json:"address1"`
	Address2		string			`json:"address2"`

	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
	UpdatedId		string			`json:"updated_id"`
	DeletedAt 		*time.Time		`json:"-"`
}


func (WebinarJoinMember) TableName() string {
	return "TB_WEBINAR_JOIN_MEMBER"
}
