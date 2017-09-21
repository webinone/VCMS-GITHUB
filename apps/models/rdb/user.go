package rdb

import "time"

type User struct {
	Idx		int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant		Tenant 		`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"tenant"`
	TenantID 	string		`json:"tenant_id"`
	UserId		string		`json:"user_id"`
	Name 		string 		`gorm:"size:255" json:"name"`
	Role 		string		`json:"role"`
	Telno		string		`json:"tel_no"`
	Email		string		`json:"email"`
	Password 	string		`json:"password"`
	UseYN		string 		`gorm:"size:5" json:"use_yn"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	UpdatedId	string		`json:"updated_id"`
}

func (User) TableName() string {
	return "TB_USER"
}