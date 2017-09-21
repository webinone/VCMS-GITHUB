package rdb


type WowzaAccesslog struct {
	Logid			int64  		`gorm:"AUTO_INCREMENT;primary_key;column:logid"`
	Date            	string		`gorm:"column:date;type:varchar(100)"`
	Time            	string		`gorm:"column:time;type:varchar(100)"`
	Tz              	string		`gorm:"column:tz;type:varchar(100)"`
	Tevent          	string		`gorm:"column:xevent;type:varchar(20)"`
	Tcategory       	string		`gorm:"column:xcategory;type:varchar(20)"`
	Xseverity       	string		`gorm:"column:xseverity;type:varchar(100)"`
	Xstatus       		string		`gorm:"column:xstatus;type:varchar(100)"`
	Xctx       		string		`gorm:"column:xctx;type:varchar(100)"`
	Xcomment       		string		`gorm:"column:xcomment;type:varchar(255)"`
	Xhost       		string		`gorm:"column:xvhost;type:varchar(100)"`
	Xapp       		string		`gorm:"column:xapp;type:varchar(100)"`
	Xappinst       		string		`gorm:"column:xappinst;type:varchar(100)"`
	Xduration       	string		`gorm:"column:xduration;type:varchar(100)"`
	Sip       		string		`gorm:"column:sip;type:varchar(100)"`
	Sport       		string		`gorm:"column:sport;type:varchar(100)"`
	Suri       		string		`gorm:"column:suri;type:varchar(255)"`
	Cip       		string		`gorm:"column:cip;type:varchar(100)"`
	Cproto       		string		`gorm:"column:cproto;type:varchar(100)"`
	Creferrer       	string		`gorm:"column:creferrer;type:varchar(255)"`
	Cuseragent      	string		`gorm:"column:cuseragent;type:varchar(100)"`
	Cclientid      		string		`gorm:"column:cclientid;type:varchar(25)"`
	Csbytes      		string		`gorm:"column:csbytes;type:varchar(20)"`
	Scbytes      		string		`gorm:"column:scbytes;type:varchar(20)"`
	Xstreamid      		string		`gorm:"column:xstreamid;type:varchar(20)"`
	Xspos      		string		`gorm:"column:xspos;type:varchar(20)"`
	Csstreambytes      	string		`gorm:"column:csstreambytes;type:varchar(20)"`
	Scstreambytes      	string		`gorm:"column:scstreambytes;type:varchar(20)"`
	Xsname      		string		`gorm:"column:xsname;type:varchar(100)"`
	Xsnamequery    		string		`gorm:"column:xsnamequery;type:varchar(100)"`
	Xfilename    		string		`gorm:"column:xfilename;type:varchar(100)"`
	Xfileext    		string		`gorm:"column:xfileext;type:varchar(100)"`
	Xfilesize    		string		`gorm:"column:xfilesize;type:varchar(100)"`
	Xfilelength    		string		`gorm:"column:xfilelength;type:varchar(100)"`
	Xsuri    		string		`gorm:"column:xsuri;type:varchar(100)"`
	Xsuristem    		string		`gorm:"column:xsuristem;type:varchar(255)"`
	Xsuriquery    		string		`gorm:"column:xsuriquery;type:varchar(255)"`
	Csuristem    		string		`gorm:"column:csuristem;type:varchar(255)"`
	Csuriquery    		string		`gorm:"column:csuriquery;type:varchar(255)"`

}

func (WowzaAccesslog) TableName() string {
	return "TB_WOWZA_ACCESSLOG"
}