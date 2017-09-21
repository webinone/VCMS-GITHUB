package rdb

import "time"

type WebinarPoll struct {
	Idx			int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string			`json:"tenant_id"`

	WebinarSiteId		string 			`json:"webinar_site_id"`
	WebinarSite		WebinarSite		`gorm:"ForeignKey:WebinarSiteId;AssociationForeignKey:WebinarSiteId" json:"webinar_site_info"`
	WebinarPollId		string 			`json:"webinar_poll_id"`



	Title			string 			`json:"title"`
	Desc			string 			`json:"desc"`

	StartDate		string 			`json:"start_date"`
	EndDate			string 			`json:"end_date"`

	WebinarPollMembers		[]WebinarPollMember     	`gorm:"ForeignKey:WebinarPollId;AssociationForeignKey:WebinarPollId" json:"members"`
	WebinarPollQuestionMasters	[]WebinarPollQuestionMaster     `gorm:"ForeignKey:WebinarPollId;AssociationForeignKey:WebinarPollId" json:"questions"`



	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
	UpdatedId		string			`json:"updated_id"`
	DeletedAt 		*time.Time		`json:"-"`
}


func (WebinarPoll) TableName() string {
	return "TB_WEBINAR_POLL"
}
