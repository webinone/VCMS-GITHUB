package rdb

import "time"

type WebinarNoticeFile struct {
	Idx			int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 		`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string		`json:"tenant_id"`

	WebinarSiteId		string 		`json:"webinar_site_id"`
	WebinarNoticeId		string 		`json:"webinar_notice_id"`

	WebinarNoticeFileId	string		`json:"webinar_notice_file_id"`

	OriginalName		string		`json:"original_name"`
	FileSize		int64		`json:"file_size"`
	SavePath		string 		`json:"save_path"`
	WebPath			string 		`json:"web_path"`

	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	UpdatedId		string		`json:"updated_id"`
}


func (WebinarNoticeFile) TableName() string {
	return "TB_WEBINAR_NOTICE_FILE"
}
