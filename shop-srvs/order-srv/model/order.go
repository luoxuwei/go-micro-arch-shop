package model

import "time"

//下单之前用户会将商品添加到购物车中，所以需要一张购物车的表
type ShoppingCart struct{
	BaseModel
	User int32 `gorm:"type:int;index"`  //加了索引，在购物车列表中我们需要查询当前用户的购物车记录，需要通过用户id查询
	Goods int32 `gorm:"type:int;index"` //加了索引，某些时候需要通过商品查购物车记录。
	                                    //如果没有必要不要加索引，不是越多越好，只有在我们需要查询时候才加。
	                                    //索引会带来负面问题：1. 会影响插入性能，在插入数据的时候索引也是需要更新的。2. 会占用磁盘
	Nums int32 `gorm:"type:int"`        //多少件商品
	Checked bool                        //是否选中
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

type OrderInfo struct{
	BaseModel

	User int32 `gorm:"type:int;index"`
	OrderSn string `gorm:"type:varchar(30);index"` //订单号，我们平台自己生成的订单号，订单号一般加上日期，这样一眼就能看出是什么时候生成的
	PayType string `gorm:"type:varchar(20) comment 'alipay(支付宝)， wechat(微信)'"` //便于查账

	//status大家可以考虑使用iota来做
	Status string `gorm:"type:varchar(20)  comment 'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
	TradeNo string `gorm:"type:varchar(100) comment '交易号'"` //交易号就是支付宝的订单号，查账时用
	OrderMount float32                             //订单金额
	PayTime *time.Time `gorm:"type:datetime"`      //支付时间

	Address string `gorm:"type:varchar(100)"`      //收货地址
	SignerName string `gorm:"type:varchar(20)"`    //收货人姓名
	SingerMobile string `gorm:"type:varchar(11)"`  //收货人手机
	Post string `gorm:"type:varchar(20)"`          //留言信息，备注
}

func (OrderInfo) TableName() string {
	return "orderinfo"
}

//订单商品，一个订单有多个商品，需要一个额外的表记录，不能做成一个切片字段放到订单表里，
//一是不变查询，如果要根据商品查询订单，放到list里查询效率低，而且不便于统计，想要看某一件商品被多个订单包含
type OrderGoods struct{
	BaseModel

	Order int32 `gorm:"type:int;index"`
	Goods int32 `gorm:"type:int;index"`

	//把商品的信息保存下来了，虽然字段冗余，但便于展示，不用再查一遍商品表，
	//还有，如果想通过商品名称查订单，要跨服务先找商品id，然后拿id再到这里查，
	//高并发系统中我们一般都不会遵循三范式，
	//还有一个，就是做镜像记录，就是用户当时在购买这件商品时，需要知道当时的名称图片是怎样的，当时的价格是什么
	//因为商家是可以改价的，名称图片都可以改，这里不保存，当时的价格就查不到了。假设产生了纠纷也能拿到当时的数据。
	GoodsName string `gorm:"type:varchar(100);index"`
	GoodsImage string `gorm:"type:varchar(200)"`
	GoodsPrice float32
	Nums int32 `gorm:"type:int"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}