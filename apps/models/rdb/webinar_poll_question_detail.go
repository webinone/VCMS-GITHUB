package rdb

import "time"

type WebinarPollQuestionDetail struct {
	Idx						int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant					Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 				string			`json:"tenant_id"`

	WebinarSiteId			string 			`json:"webinar_site_id"`
	WebinarPollId			string 			`json:"webinar_poll_id"`
	WebinarPollQuestionMasterId	string 			`json:"webinar_poll_question_master_id"`
	WebinarPollQuestionDetailId	string 			`json:"webinar_poll_question_detail_id"`

	Order				int			`gorm:"column:detail_order" json:"detail_order"`

	Title				string 			`json:"title"`

	ResultCount			string			`gorm:"-" json:"result_count"`
	ResultPercent		string			`gorm:"-" json:"result_percent"`

	CreatedAt			time.Time		`json:"created_at"`
	UpdatedAt			time.Time		`json:"updated_at"`
	UpdatedId			string			`json:"updated_id"`
	DeletedAt 			*time.Time		`json:"-"`
}


func (WebinarPollQuestionDetail) TableName() string {
	return "TB_WEBINAR_POLL_QUESTION_DETAIL"
}
