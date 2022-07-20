package model

import (
	"database/sql/driver"
	"encoding/json"
)

type Inventory struct{
	BaseModel
	Goods int32 `gorm:"type:int;index"`
	Stocks int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` //分布式锁的乐观锁
}

type GoodsDetail struct {
	Goods int32
	Num int32
}
type GoodsDetailList []GoodsDetail

func (g GoodsDetailList) Value() (driver.Value, error){
	return json.Marshal(g)
}

func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type StockSellDetail struct {
	OrderSn string `gorm:"type:varchar(200);index:idx_order_sn,unique;"`
	Status int32 `gorm:"type:varchar(200)"` //1 表示已扣减 2. 表示已归还
	Detail GoodsDetailList `gorm:"type:varchar(200)"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}