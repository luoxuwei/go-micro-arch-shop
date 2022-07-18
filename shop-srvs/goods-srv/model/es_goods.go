package model

//为了将es和mysql做一一对应，应该建一个唯一的id，所以就可以将mysql的id直接拿古来，
//不让es自动生成，否则一不小心同样的数据会存在多分，
//只将与搜索过滤相关的字段存到es
type EsGoods struct {
	ID int32 `json:"id"`
	//可能通过category和brand来过滤，OnSele是否在售，ShipFree是否免运费，
	//IsNew新品，IsHot热门，这些都有可能作为过滤条件
	CategoryID int32 `json:"category_id"`
	BrandsID int32 `json:"brands_id"`
	OnSale bool  `json:"on_sale"`
	ShipFree bool  `json:"ship_free"`
	IsNew bool `json:"is_new"`
	IsHot bool `json:"is_hot"`

	//name是重点查询字段
	Name  string `json:"name"`
	//点击数，如果有需要可以通过它来排序
	ClickNum int32  `json:"click_num"`
	SoldNum int32 `json:"sold_num"`
	FavNum int32 `json:"fav_num"`
	MarketPrice float32  `json:"market_price"`
	GoodsBrief string `json:"goods_brief"`
	ShopPrice float32 `json:"shop_price"`
}

//虽然es可以自动生成，但有些配置需要定制，最重要的是，设置中文分词的analyzer。
func (EsGoods) GetMapping() string {
	goodsMapping := `
	{
		"mappings" : {
			"properties" : {
				"brands_id" : {
					"type" : "integer"
				},
				"category_id" : {
					"type" : "integer"
				},
				"click_num" : {
					"type" : "integer"
				},
				"fav_num" : {
					"type" : "integer"
				},
				"id" : {
					"type" : "integer"
				},
				"is_hot" : {
					"type" : "boolean"
				},
				"is_new" : {
					"type" : "boolean"
				},
				"market_price" : {
					"type" : "float"
				},
				"name" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"goods_brief" : {
					"type" : "text",
					"analyzer":"ik_max_word"
				},
				"on_sale" : {
					"type" : "boolean"
				},
				"ship_free" : {
					"type" : "boolean"
				},
				"shop_price" : {
					"type" : "float"
				},
				"sold_num" : {
					"type" : "long"
				}
			}
		}
	}`
	return goodsMapping
}

func (EsGoods) GetIndexName() string {
	return "goods"
}
