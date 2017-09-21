package rdb

import "time"

type Category struct {
	Idx			int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant			Tenant 		`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantId 		string		`json:"tenant_id"`
	CategoryId		string		`json:"category_id"`
	Name 			string 		`gorm:"size:255" json:"name"`
	Desc 			string		`json:"desc"`
	ParentIdx		int64		`json:"parent_idx"`
	UseYN			string 		`gorm:"size:5" json:"use_yn"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	UpdatedId		string		`json:"updated_id"`
}


func (Category) TableName() string {
	return "TB_CATEGORY"
}
