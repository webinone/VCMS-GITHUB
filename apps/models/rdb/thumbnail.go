package rdb

import "time"

type ThumbNail struct {
	Idx			int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	ContentID 		string		`json:"content_id"`
	SavePath		string 		`json:"save_path"`
	WebPath			string 		`json:"web_path"`
	//Size 			int64		`json:"size"`
	Resolution  		string 		`json:"resolution"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	//DeletedAt 		*time.Time	`json:"-"`
}

func (ThumbNail) TableName() string {
	return "TB_THUMBNAIL"
}