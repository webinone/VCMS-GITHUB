package rdb


type VcmsAccessLog struct {
	Idx				int64  		`gorm:"AUTO_INCREMENT;primary_key" json:"idx"`
	Date            		string		`gorm:"column:date;type:varchar(100)" json:"date"`
	Time            		string		`gorm:"column:time;type:varchar(100)" json:"time"`
	Year				string 		`gorm:"column:year;type:varchar(25)" json:"year"`
	Month				string 		`gorm:"column:month;type:varchar(25)" json:"month"`
	Day				string 		`gorm:"column:day;type:varchar(25)" json:"day"`
	Hour				string 		`gorm:"column:hour;type:varchar(25)" json:"hour"`
	Minute				string 		`gorm:"column:minute;type:varchar(25)" json:"minute"`
	GroupYearMonth          	string		`gorm:"column:group_year_month;type:varchar(255)" json:"group_year_month"`
	GroupYearMonthDay       	string		`gorm:"column:group_year_month_day;type:varchar(255)" json:"group_year_month_day"`
	GroupYearMonthDayHour   	string		`gorm:"column:group_year_month_day_hour;type:varchar(255)" json:"group_year_month_day_hour"`
	GroupYearMonthDayHourMinute     string		`gorm:"column:group_year_month_day_hour_minute;type:varchar(255)" json:"group_year_month_day_hour_minute"`
	Tz              		string		`gorm:"column:tz;type:varchar(100)" json:"tz"`
	TenantId 			string		`json:"tenant_id"`
	ChannelId 			string		`json:"channel_id"`
	ContentId			string		`json:"content_id"`
	Xsname      			string		`gorm:"column:xsname;type:varchar(100)" json:"xsname"`
	Xduration       		string		`gorm:"column:xduration;type:varchar(100)" json:"xduration"`
	Cip       			string		`gorm:"column:cip;type:varchar(100)" json:"cip"`
	Cuseragent			string		`gorm:"column:cuseragent;type:varchar(100)" json:"cuseragent"`
	Cclientid			string		`gorm:"column:cclientid;type:varchar(25)" json:"cclientid"`
	Xfilename    			string		`gorm:"column:xfilename;type:varchar(100)" json:"xfilename"`
	Cproto       			string		`gorm:"column:cproto;type:varchar(100)" json:"cproto"`
	Xsuri    			string		`gorm:"column:xsuri;type:varchar(100)" json:"xsuri"`
}

func (VcmsAccessLog) TableName() string {
	return "TB_VCMS_ACCESS_LOG"
}