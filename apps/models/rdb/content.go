package rdb

import "time"

type Content struct {
	Idx					int64  			`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Tenant				Tenant 			`gorm:"ForeignKey:TenantID;AssociationForeignKey:TenantID" json:"-"`
	TenantID 			string			`json:"tenant_id"`
	Category			Category		`json:"-"`
	CategoryId			string			`json:"category_id"`
	ThumbNails			[]ThumbNail		`gorm:"ForeignKey:ContentId;AssociationForeignKey:ContentId" json:"thumbnails"`
	ContentId			string			`json:"content_id"`
	ContentName 		string			`json:"content_name"`
	ContentTags			[]ContentTag	`gorm:"ForeignKey:ContentId;AssociationForeignKey:ContentId" json:"tags"`
	OriginFileName 		string 			`gorm:"size:255" json:"origin_filename"`
	GeneratedFileName 	string 			`gorm:"size:255" json:"generated_filename"`
	Stream  			Stream			`gorm:"ForeignKey:ContentId;AssociationForeignKey:ContentId" json:"stream"`
	Type				string 			`json:"file_type"`
	Ext 				string 			`json:"file_ext"`
	FilePath			string 			`json:"file_path"`
	Duration 			string			`json:"duration"`
		Size			int64			`json:"file_size"`
	ThumbType			string 			`json:"thumb_type"`
	ThumbTime			string			`json:"thumb_time"`
	CreatedAt			time.Time		`json:"created_at"`
	UpdatedAt			time.Time		`json:"updated_at"`
	UpdatedId			string			`json:"-"`
	DeletedAt 			*time.Time		`json:"-"`
}

func (Content) TableName() string {
	return "TB_CONTENT"
}