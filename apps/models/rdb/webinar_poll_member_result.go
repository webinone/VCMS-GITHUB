package rdb

import "time"

type WebinarPollMemberResult struct {
	Idx				int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant				Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 			string			`json:"tenant_id"`

	WebinarSiteId			string 			`json:"webinar_site_id"`
	WebinarPollId			string 			`json:"webinar_poll_id"`
	WebinarPollMemberId		string			`json:"webinar_poll_member_id"`
	WebinarPollQuestionMasterId	string 			`json:"webinar_poll_question_master_id"`

	WebinarPollQuestionMaster	WebinarPollQuestionMaster     `gorm:"ForeignKey:WebinarPollQuestionMasterId;AssociationForeignKey:WebinarPollQuestionMasterId" json:"question"`

	FrontUserId			string			`json:"front_user_id"`

	Answer				string			`json:"webinar_poll_answer"`	// detail id 또는 그냥 주관식 답변

	WebinarPollQuestionDetail	WebinarPollQuestionDetail	`gorm:"ForeignKey:WebinarPollQuestionDetailId;AssociationForeignKey:Answer" json:"question_detail"`

	CreatedAt			time.Time		`json:"created_at"`
	UpdatedAt			time.Time		`json:"updated_at"`
	UpdatedId			string			`json:"updated_id"`
	DeletedAt 			*time.Time		`json:"-"`
}


func (WebinarPollMemberResult) TableName() string {
	return "TB_WEBINAR_POLL_MEMBER_RESULT"
}
