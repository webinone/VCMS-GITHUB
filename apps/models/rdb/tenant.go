package rdb

import (
	"time"
)

type Tenant struct {
	Idx		int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	TenantId	string 		`gorm:"not null;unique" json:"tenant_id"`
	Name 	 	string 		`gorm:"size:255" json:"name"`
	TenantDesc 	string		`json:"tenant_desc"`
	TimeZone	string		`json:"timezone"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	DeletedAt 	*time.Time	`json:"-"`
}


func (Tenant) TableName() string {
	return "TB_TENANT"
}