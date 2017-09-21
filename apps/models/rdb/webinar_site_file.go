package rdb

import "time"

type WebinarSiteFile struct {
	Idx				int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string			`json:"tenant_id"`

	WebinarSiteId		string 			`json:"webinar_site_id"`
	WebinarSiteFileId	string		`json:"webinar_site_file_id"`

	OriginalName		string		`json:"original_name"`
	FileSize			int64		`json:"file_size"`
	SavePath			string 		`json:"save_path"`
	WebPath				string 		`json:"web_path"`

	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
	UpdatedId		string			`json:"updated_id"`
	DeletedAt 		*time.Time		`json:"-"`
}


func (WebinarSiteFile) TableName() string {
	return "TB_WEBINAR_SITE_FILE"
}
