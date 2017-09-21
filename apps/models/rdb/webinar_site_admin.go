package rdb

import "time"

type WebinarSiteAdmin struct {
	Idx				int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string			`json:"tenant_id"`

	WebinarSiteId		string 		`json:"webinar_site_id"`

	MemberId		string			`json:"member_id"`
	MemberName		string			`json:"member_name"`

	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
}


func (WebinarSiteAdmin) TableName() string {
	return "TB_WEBINAR_SITE_ADMIN"
}
