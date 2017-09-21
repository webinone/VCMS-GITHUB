package rdb

//seqNo int(10) unsigned not null,
//apply_type varchar(20) default 'GENERAL' null,
//member_gubun enum('personal', 'company') null,
//member_name varchar(40) not null,
//member_id varchar(20) not null,
//member_pwd varchar(50) not null,
//emailAddr varchar(50) null,
//eletter enum('Y', 'N') default 'N' null,
//phoneNum varchar(20) null,
//biznumber varchar(30) null,
//com_name varchar(50) null,
//registDate int(9) not null,
//loginDate int(9) default '0' null,
//ipAddr varchar(20) null,
//infoAdd enum('Y', 'N') default 'N' not null comment '업체정보등록여부',
//updDate int(9) null,
//modifyAddr int(9) default '0' null comment '도로명주소 :$ntime, 지번주소 :0',
//temp_yn varchar(1) default 'N' null comment '임시정보 여부',
//sms varchar(1) default 'N' null,
//OS varchar(30) null comment ' 사용자 OS',
//browser varchar(30) null comment ' 사용자 브라우저',
//chgPwd enum('Y', 'N') default 'N' null
type MemberDefault struct {

	SeqNo			int  		`gorm:"primary_key;column:seqNo" json:"seqNo"`
	ApplyType            	string		`gorm:"column:apply_type;type:varchar(20)" json:"apply_type"`
	MemberGubun            	string		`gorm:"column:member_gubun" json:"member_gubun"`
	MemberName             	string		`gorm:"column:member_name;type:varchar(40)" json:"member_name"`
	MemberId          	string		`gorm:"column:member_id;type:varchar(20)" json:"member_id"`
	MemberPwd       	string		`gorm:"column:member_pwd;type:varchar(50)" json:"member_pwd"`
	EmailAddr       	string		`gorm:"column:emailAddr;type:varchar(50)" json:"emailAddr"`
	ELetter       		string		`gorm:"column:eletter" json:"eletter"`
	PhoneNum       		string		`gorm:"column:phoneNum;type:varchar(20)" json:"phoneNum"`
	BizNumber      		string		`gorm:"column:biznumber;type:varchar(30)" json:"biznumber"`
	ComName      		string		`gorm:"column:com_name;type:varchar(30)" json:"com_name"`
	RegistDate      	int		`gorm:"column:registDate;type:int(9)" json:"registDate"`
	LoginDate      		int		`gorm:"column:loginDate;type:int(9)" json:"loginDate"`

	IpAddr          	string		`gorm:"column:ipAddr;type:varchar(20)" json:"ipAddr"`
	InfoAdd          	string		`gorm:"column:infoAdd" json:"infoAdd"`
	UpdDate          	string		`gorm:"column:updDate" json:"updDate"`
	ModifyAddr          	string		`gorm:"column:modifyAddr" json:"modifyAddr"`
	TempYN          	string		`gorm:"column:temp_yn;type:varchar(1)" json:"temp_yn"`
	SMS          		string		`gorm:"column:sms;type:varchar(1)" json:"sms"`
	OS          		string		`gorm:"column:os;type:varchar(30)" json:"os"`
	Browser          	string		`gorm:"column:browser;type:varchar(30)" json:"browser"`
	ChPwd	          	string		`gorm:"column:chgPwd" json:"chgPwd"`

	MemberSub		MemberSub	`gorm:"ForeignKey:MemberId;AssociationForeignKey:MemberId" json:"member_sub"`

}

func (MemberDefault) TableName() string {
	return "member_default"
}