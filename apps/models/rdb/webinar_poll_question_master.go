package rdb

import "time"

type WebinarPollQuestionMaster struct {
	Idx				int64  				`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant				Tenant 				`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 			string				`json:"tenant_id"`

	WebinarSiteId			string 				`json:"webinar_site_id"`



	WebinarPollId			string 				`json:"webinar_poll_id"`
	WebinarPollQuestionMasterId	string 				`json:"webinar_poll_question_master_id"`

	WebinarPollQuestionDetails	[]WebinarPollQuestionDetail	`gorm:"ForeignKey:WebinarPollQuestionMasterId;AssociationForeignKey:WebinarPollQuestionMasterId" json:"webinar_question_detail"`

	Title				string 				`json:"title"`
	QuestionType			string 				`json:"question_type"`  // 1:객관식 2:주관식

	QuestionCount			string				`json:"question_count"`

	Answer			string			`gorm:"-" json:"webinar_poll_answer"`

	CreatedAt			time.Time			`json:"created_at"`
	UpdatedAt			time.Time			`json:"updated_at"`
	UpdatedId			string				`json:"updated_id"`
	DeletedAt 			*time.Time			`json:"-"`
}


func (WebinarPollQuestionMaster) TableName() string {
	return "TB_WEBINAR_POLL_QUESTION_MASTER"
}
