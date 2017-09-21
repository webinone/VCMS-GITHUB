package rdb

import "time"

type WebinarSiteTag struct {
	Idx				int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 		`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string		`json:"tenant_id"`

	WebinarSiteId	string 		`json:"webinar_site_id"`
	Name 			string 		`gorm:"size:255" json:"name"`

	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	UpdatedId		string		`json:"updated_id"`
	DeletedAt 		*time.Time	`json:"-"`
}


func (WebinarSiteTag) TableName() string {
	return "TB_WEBINAR_SITE_TAG"
}
