package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop-srvs/goods-srv/global"
	"shop-srvs/goods-srv/model"
	"shop-srvs/goods-srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

//商品的字段比较多，最好写一个公共的函数
func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse {
		Id:       goods.ID,
		CategoryId: goods.CategoryID,
		Name: goods.Name,
		GoodsSn: goods.GoodsSn,
		ClickNum: goods.ClickNum,
		SoldNum: goods.SoldNum,
		FavNum: goods.FavNum,
		MarketPrice: goods.MarketPrice,
		ShopPrice: goods.ShopPrice,
		GoodsBrief: goods.GoodsBrief,
		ShipFree: goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew: goods.IsNew,
		IsHot: goods.IsHot,
		OnSale: goods.OnSale,
		DescImages: goods.DescImages,
		Images: goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	localDB := global.DB.Model(model.Goods{})
	if req.KeyWords != "" {
		localDB = localDB.Where("name LIKE ?", "%" +req.KeyWords+"%")
	}

	if req.IsHot {
		localDB = localDB.Where(model.Goods{IsHot:true})
	}

	if req.IsNew {
		localDB = localDB.Where(model.Goods{IsNew:true})
	}

	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price >= ?", req.PriceMin)
	}

	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price <= ?", req.PriceMax)
	}

	if req.Brand > 0 {
		localDB = localDB.Where("brand_id = ?", req.Brand)
	}

	//通过category去查询商品
	var subQuery string
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}

		localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	}

   var total int64
	localDB.Count(&total)
	goodsListResponse.Total = int32(total)
    //要在分页之前拿到total

	var goods []model.Goods
	result := localDB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
    if result.Error != nil {
    	return nil, result.Error
	}

	for _, good := range goods {
        goodsInfoResponse := ModelToResponse(good)
        goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

    return goodsListResponse, nil
}