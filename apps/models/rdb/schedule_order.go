package rdb

import "time"

type ScheduleOrder struct {
	Idx			int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 		`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantID 		string		`json:"tenant_id"`
	ScheduleId 		string		`json:"schedule_id"`
	ScheduleOrderId		string 		`gorm:"size:10" json:"schedule_order_id"`
	StartSec		string 		`gorm:"size:10" json:"start_sec"`
	EndSec			string 		`json:"end_sec"`
	Content			Content		`gorm:"ForeignKey:ContentId;AssociationForeignKey:ContentId" json:"content"`
	ContentId		string		`json:"content_id"`
	Order                   int		`json:"order"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	UpdatedId		string		`json:"updated_id"`
}

func (ScheduleOrder) TableName() string {
	return "TB_SCHEDULE_ORDER"
}