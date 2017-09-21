package rdb

import "time"

type Stream struct {
	Idx			int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"-"`
	ChannelId		string 		`json:"-"`
	ContentId		string		`json:"-"`
	MpegDash		string 		`json:"mpeg_dash"`
	RTMP			string  	`json:"rtmp"`
	HLS			string		`json:"hls"`
	HDS             	string		`json:"hds"`
	IOS         		string 		`json:"ios"`
	Android         	string  	`json:"rtsp"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
	DeletedAt 		*time.Time	`json:"-"`
}

func (Stream) TableName() string {
	return "TB_STREAM"
}

