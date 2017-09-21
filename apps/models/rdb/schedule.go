package rdb

import "time"

type Schedule struct {
	Idx			int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 		`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantID 		string		`json:"tenant_id"`
	Channel			Channel		`gorm:"ForeignKey:ChannelId;AssociationForeignKey:ChannelId" json:"channel"`
	Stream			Stream		`gorm:"ForeignKey:ChannelId;AssociationForeignKey:ChannelId" json:"stream"`
	ChannelId		string		`json:"channel_id"`
	ScheduleId 		string		`json:"schedule_id"`
	Name			string 		`json:"schedule_name"`
	StartDateTime		string 		`json:"start_datetime"`
	TotalTime		string 		`json:"total_time"`
	ScheduleOrders   	[]ScheduleOrder	`gorm:"ForeignKey:ScheduleId;AssociationForeignKey:ScheduleId" json:"schedule_orders"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	UpdatedId		string		`json:"-"`
	DeletedAt 		*time.Time	`json:"-"`

}

type Streams struct {

}

func (Schedule) TableName() string {
	return "TB_SCHEDULE"
}