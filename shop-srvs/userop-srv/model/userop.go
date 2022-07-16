package model

//留言
type LeavingMessages struct{
	BaseModel

	User int32 `gorm:"type:int;index"`
	MessageType int32 `gorm:"type:int comment '留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)'"`
	Subject string `gorm:"type:varchar(100)"`

	Message string //长度比较大的话，不用加类型，默认是text
	File string `gorm:"type:varchar(200)"` //保存的是url，真正的文件是会保存到阿里云中的
}

func (LeavingMessages) TableName() string {
	return "leavingmessages"
}

//省、市、区域、地址、名称、手机号码
type Address struct{
	BaseModel

	User int32 `gorm:"type:int;index"`
	Province string `gorm:"type:varchar(10)"`
	City string `gorm:"type:varchar(10)"`
	District string `gorm:"type:varchar(20)"`
	Address string `gorm:"type:varchar(100)"`
	SignerName string `gorm:"type:varchar(20)"`
	SignerMobile string `gorm:"type:varchar(11)"`
}

//用户id和商品id，用联合唯一索引，指定相同的索引名称就行了。而且必须是唯一的
type UserFav struct{
	BaseModel

	User int32 `gorm:"type:int;index:idx_user_goods,unique"`
	Goods int32 `gorm:"type:int;index:idx_user_goods,unique"`
}

func (UserFav) TableName() string {
	return "userfav"
}