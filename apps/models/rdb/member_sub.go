package rdb


type MemberSub struct {

	MemberId          	string		`gorm:"column:member_id;type:varchar(20)" json:"member_id"`
	ZipCode       		string		`gorm:"column:zipcode;type:char(7)" json:"zipcode"`

	Address1       		string		`gorm:"column:address1;type:varchar(200)" json:"address1"`
	Address2       		string		`gorm:"column:address2;type:varchar(100)" json:"address2"`

	MobileNum      		string		`gorm:"column:mobileNum;type:varchar(20)" json:"mobileNum"`
	CompanyName    		string		`gorm:"column:company_name;type:varchar(40)" json:"company_name"`

	BType    		string		`gorm:"column:bType;type:varchar(30)" json:"bType"`
	BKind    		string		`gorm:"column:bKind;type:varchar(30)" json:"bKind"`
	Staff    		string		`gorm:"column:staff;type:varchar(30)" json:"staff"`

	Interest    		string		`gorm:"column:interest;type:varchar(30)" json:"interest"`
	etcText    		string		`gorm:"column:etcText;type:text" json:"etcText"`
	temp_yn    		string		`gorm:"column:temp_yn;type:varchar(1)" json:"temp_yn"`
}

func (MemberSub) TableName() string {
	return "member_sub"
}