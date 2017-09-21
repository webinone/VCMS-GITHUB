package rdb

import "time"

type Channel struct {
	Idx		int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant		Tenant 		`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantID 	string		`json:"tenant_id"`
	ChannelId 	string		`json:"channel_id"`
	Stream		Stream		`gorm:"ForeignKey:ChannelId;AssociationForeignKey:ChannelId" json:"stream"`
	Name		string 		`json:"channel_name"`
	Bitrate		string 		`json:"bitrate"`
	Shedules        []Schedule     	`gorm:"ForeignKey:ChannelId;AssociationForeignKey:ChannelId" json:"schedules"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	UpdatedId	string		`json:"-"`
	DeletedAt 	*time.Time	`json:"-"`
}

func (Channel) TableName() string {
	return "TB_CHANNEL"
}