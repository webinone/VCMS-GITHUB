package rdb

import "time"

type ContentTag struct {
	Idx			int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`

	Tenant			Tenant 		`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string		`json:"tenant_id"`

	ContentId		string		`json:"content_id"`

	Name 			string 		`gorm:"size:255" json:"name"`

	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	UpdatedId		string		`json:"updated_id"`
	DeletedAt 		*time.Time	`json:"-"`
}

func (ContentTag) TableName() string {
	return "TB_CONTENT_TAG"
}