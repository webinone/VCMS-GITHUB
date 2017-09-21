package rdb

import "time"

type WebinarBanner struct {
	Idx			int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 		`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string		`json:"tenant_id"`

	WebinarSiteId		string 		`json:"webinar_site_id"`

	BannerId		string		`json:"banner_id"`
	BannerTitle		string		`json:"banner_title"`
	BannerDesc		string		`json:"banner_desc"`
	BannerType		string		`json:"banner_type"` // 1: 경품배너 2:업체배너 3: 진행페이지 배너

	BannerOrder             int		`json:"banner_order"`
	SavePath		string 		`json:"save_path"`
	WebPath			string 		`json:"web_path"`

	LinkUrl			string		`json:"link_url"`

	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}

func (WebinarBanner) TableName() string {
	return "TB_WEBINAR_BANNER"
}
