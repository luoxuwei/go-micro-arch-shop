package model


//虽然有三级分类，但分类的信息都差不多，所以用一张表，通过指定父类和子类来表示层级关系
type Category struct{
	BaseModel
	Name  string `gorm:"type:varchar(20);not null" json:"name"`
	ParentCategoryID int32 `json:"parent"`
	ParentCategory *Category `json:"-"`
	SubCategory []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level int32 `gorm:"type:int;not null;default:1" json:"level"`
	IsTab bool `gorm:"default:false;not null" json:"is_tab"`
}

type Brands struct {
	BaseModel
	Name  string `gorm:"type:varchar(20);not null"`
	Logo  string `gorm:"type:varchar(200);default:'';not null"`
}

//品牌和分类是多对多的关系，通过分类可以过滤品牌，品牌也有多个分类，所以需要建一张表，保存两者关系
//需要建立品牌id和分类id的联合索引, 设置索引名一样就是联合索引
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands Brands
}

//gorm的默认表名用下划线连接每个单词，如果不喜欢可以重载Tablename方法
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	//跳转商品的url
	Url string `gorm:"type:varchar(200);not null"`
	//轮播图是有先后顺序的，所以加一个index表名顺序
	Index int32 `gorm:"type:int;default:1;not null"`
}

type Goods struct {
	BaseModel

	//两个外键，分类和品牌
	CategoryID int32 `gorm:"type:int;not null"`
	Category Category
	BrandsID int32 `gorm:"type:int;not null"`
	Brands Brands

	//是否已上架，上架了才能被搜索到
	OnSale bool `gorm:"default:false;not null"`
	//是否免运费
	ShipFree bool `gorm:"default:false;not null"`
	//是否是新品，根据需求来的，一般都有一个新品栏，做推广。
	IsNew bool `gorm:"default:false;not null"`
	//是否是热卖商品，这是用来做广告位的，比如你给了钱，那就把它做成热卖商品
	IsHot bool `gorm:"default:false;not null"`

	//商品名称要长一些，做过优化的都知道，很多人都会加关键字做优化
	Name  string `gorm:"type:varchar(50);not null"`
	//商品编号，这个编号不是数据库自己生成的编号，二是商家自己内部的编号，一旦你下单了，他会拿这个编号去仓库里找商品
	//商家自己有一套仓库管理系统，如果没这个编号，商家自己都不知道去怎么找，
	GoodsSn string `gorm:"type:varchar(50);not null"`
	//点击数，可以点击数越高就表示越畅销，这些对后期做数据分析很重要，所以一般情况下对商品的信息记录的详细一点
	ClickNum int32 `gorm:"type:int;default:0;not null"`
	//销售数，已经买了多少件了，就是销量
	SoldNum int32 `gorm:"type:int;default:0;not null"`
	//收藏数量，收藏数高，说明受欢迎的程度高
	FavNum int32 `gorm:"type:int;default:0;not null"`
	//市场价，商品价格, 商品一般有两个价格，一个是市场价，一个是本地的价格，或者说叫平时价和活动价, 所以这里提供两个价格字段
	MarketPrice float32 `gorm:"not null"`
	//本地价
	ShopPrice float32 `gorm:"not null"`
	//商品的简介
	GoodsBrief string `gorm:"type:varchar(100);not null"`
	//商品的图片，有三个地方，一个是封面图，一个是简介里的图片列表，一个是详情里的图片列表
	//代码中图片列表就是一个string类型的切片[]string，但在这里不能用[]string, 在数据库里没有数组类型，只有json类型
	//有两种方式，一个定义一个新类型，typ GormList []string ，另外一个方式是建一张新表goods_image，商品图片的表，把商品id和图片url对应起来
	//如果新建一张表，在查询商品信息时肯定需要将图片列表拉出来的，就需要用表join，当数量越来越大join带来的性能降低是非常明显的，所以这里用定义新类型的方式
	Images GormList `gorm:"type:varchar(1000);not null"`
	DescImages GormList `gorm:"type:varchar(1000);not null"`
	//封面图
	GoodsFrontImage string `gorm:"type:varchar(200);not null"`
}