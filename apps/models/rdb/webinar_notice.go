package rdb

import "time"

type WebinarNotice struct {
	Idx			int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string			`json:"tenant_id"`

	WebinarSiteId		string 			`json:"webinar_site_id"`
	WebinarNoticeId		string 			`json:"webinar_notice_id"`

	Title 			string 			`gorm:"size:255" json:"title"`
	Content 		string			`json:"content"`

	WebinarNoticeFiles	[]WebinarNoticeFile	`gorm:"ForeignKey:WebinarNoticeId;AssociationForeignKey:WebinarNoticeId" json:"webinar_notice_files"`

	UseYN			string			`json:"use_yn"`

	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
	UpdatedId		string			`json:"updated_id"`
	DeletedAt 		*time.Time		`json:"-"`
}


func (WebinarNotice) TableName() string {
	return "TB_WEBINAR_NOTICE"
}
