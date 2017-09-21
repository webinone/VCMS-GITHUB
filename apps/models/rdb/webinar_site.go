package rdb

import "time"

type WebinarSite struct {
	Idx				int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string			`json:"tenant_id"`

	WebinarSiteId		string 		`json:"webinar_site_id"`

	Title 				string 		`gorm:"size:255" json:"title"`
	SubTitle 			string		`json:"sub_title"`

	StartDateTime		string		`json:"start_datetime"`
	TotalTime			string		`json:"total_time"`
	EndDateTime			string		`json:"end_datetime"`
	EnterBeforeTime		string		`json:"enter_before_time"`
	EnterDateTime		string		`json:"enter_datetime"`

	PresenterName		string		`json:"presenter_name"`
	PresenterCompany	string		`json:"presenter_company"`
	PresenterCompanyNo	string		`json:"presenter_company_no"`
	PresenterDep		string		`json:"presenter_dep"`
	PresenterPosition	string		`json:"presenter_position"`
	PresenterEmail		string		`json:"presenter_email"`
	PresenterPhone		string		`json:"presenter_phone"`
	PresenterFax		string		`json:"presenter_fax"`

	PreviewType			string 		`gorm:"size:10" json:"preview_type"`  // 1: 내부 , 2: 외부
	PreviewContent		Content		`gorm:"ForeignKey:ContentId;AssociationForeignKey:PreviewContentId" json:"preview_content"`
	PreviewContentId	string		`json:"preview_content_id"`
	PreviewExtraUrl		string		`json:"preview_extra_url"`

	Channel				Channel			`gorm:"ForeignKey:ChannelId;AssociationForeignKey:ChannelId" json:"channel"`
	ChannelId			string			`json:"channel_id"`

	Schedule			Schedule		`gorm:"ForeignKey:ScheduleId;AssociationForeignKey:ScheduleId" json:"schedule"`
	ScheduleId			string			`json:"schedule_id"`

	PostType			string 			`gorm:"size:10" json:"post_type"`  // 1: 내부 , 2: 외부
	PostContent			Content			`gorm:"ForeignKey:ContentId;AssociationForeignKey:PostContentId" json:"post_content"`
	PostContentId		string			`json:"post_content_id"`
	PostExtraUrl		string			`json:"post_extra_url"`

	HostName			string			`json:"host_name"`
	HostEmail			string			`json:"host_email"`
	HostPhone			string			`json:"host_phone"`
	HostFax				string			`json:"host_fax"`

	YoutubeLiveUrl		string			`json:"youtube_live_url"`

	WebinarBanners		[]WebinarBanner  		`gorm:"ForeignKey:WebinarSiteId;AssociationForeignKey:WebinarSiteId" json:"webinar_banners"`

	WebinarSiteAdmins	[]WebinarSiteAdmin  	`gorm:"ForeignKey:WebinarSiteId;AssociationForeignKey:WebinarSiteId" json:"webinar_admins"`

	WebinarSiteTags     []WebinarSiteTag 		`gorm:"ForeignKey:WebinarSiteId;AssociationForeignKey:WebinarSiteId" json:"webinar_tags"`

	WebinarSiteThumbNails []WebinarSiteThumbNail `gorm:"ForeignKey:WebinarSiteId;AssociationForeignKey:WebinarSiteId" json:"webinar_thumbnails"`

	WebinarSiteFiles []WebinarSiteFile 			`gorm:"ForeignKey:WebinarSiteId;AssociationForeignKey:WebinarSiteId" json:"webinar_site_files"`

	WebinarSiteBackImages []WebinarSiteBackImage `gorm:"ForeignKey:WebinarSiteId;AssociationForeignKey:WebinarSiteId" json:"webinar_site_backimages"`

	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	UpdatedId		string		`json:"updated_id"`
	DeletedAt 		*time.Time	`json:"-"`
}


func (WebinarSite) TableName() string {
	return "TB_WEBINAR_SITE"
}
