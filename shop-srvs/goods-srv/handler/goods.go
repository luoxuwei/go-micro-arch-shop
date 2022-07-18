package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
	//match bool 复合查询
	q := elastic.NewBoolQuery()

	if req.KeyWords != "" {
		//localDB = localDB.Where("name LIKE ?", "%" +req.KeyWords+"%")
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}

	if req.IsHot {
		//localDB = localDB.Where(model.Goods{IsHot:true})
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}

	if req.IsNew {
		//localDB = localDB.Where(model.Goods{IsNew:true})
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}

	if req.PriceMin > 0 {
		//localDB = localDB.Where("shop_price >= ?", req.PriceMin)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}

	if req.PriceMax > 0 {
		//localDB = localDB.Where("shop_price <= ?", req.PriceMax)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}

	if req.Brand > 0 {
		//localDB = localDB.Where("brand_id = ?", req.Brand)
		q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}

	//通过category去查询商品
	var subQuery string
	categoryIds := make([]interface{}, 0)
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

		type Result struct {
			ID int32
		}
		var results []Result
		global.DB.Model(model.Category{}).Raw(subQuery).Scan(&results)
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}
		//生成terms查询 相当于sql里的in
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
		//localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	}

	//对分页参数做安全校验，
	//   1.防止pages和nums都为0，这样查不出任何结果
	//   2.防止pagesize过大，一是性能考虑，二是如果是爬虫，一下就拿到大把数据
	if req.Pages == 0 {
		req.Pages = 1
	}

	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}

	//GET goods/_search
	result, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil, err
	}

	goodsIds := make([]int32, 0)
	goodsListResponse.Total = int32(result.Hits.TotalHits.Value)
	for _, value := range result.Hits.Hits {
		goods := model.EsGoods{}
		_ = json.Unmarshal(value.Source, &goods)
		goodsIds = append(goodsIds, goods.ID)
	}

   //var total int64
	//localDB.Count(&total)
	//goodsListResponse.Total = int32(total)
   // //要在分页之前拿到total

	var goods []model.Goods
	re := localDB.Preload("Category").Preload("Brands").Find(&goods, goodsIds)
	if re.Error != nil {
		return nil, re.Error
	}
	//result := localDB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
    //if result.Error != nil {
    //	return nil, result.Error
	//}

	for _, good := range goods {
        goodsInfoResponse := ModelToResponse(good)
        goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

    return goodsListResponse, nil
}

//批量获取商品
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error){
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods

	result := global.DB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}

func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error){
	var goods model.Goods

	if result := global.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods := model.Goods{
		Brands: brand,
		BrandsID: brand.ID,
		Category: category,
		CategoryID: category.ID,
		Name: req.Name,
		GoodsSn: req.GoodsSn,
		MarketPrice: req.MarketPrice,
		ShopPrice: req.ShopPrice,
		GoodsBrief: req.GoodsBrief,
		ShipFree: req.ShipFree,
		Images: req.Images,
		DescImages: req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew: req.IsNew,
		IsHot: req.IsHot,
		OnSale: req.OnSale,
	}

	//global.DB.Save(&goods)
    //要确保es和mysql操作的一致性，不能在入库是发生一个成功一个失败的情况
    //所以这里加事务，Save方法会调用afterCreate，可以通过返回值判断是否成功。
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	return &proto.GoodsInfoResponse{
		Id:  goods.ID,
	}, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	tx := global.DB
	tx.Begin()
	result := tx.Delete(&model.Goods{BaseModel:model.BaseModel{ID:req.Id}}, req.Id)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	//if result := global.DB.Delete(&model.Goods{BaseModel:model.BaseModel{ID:req.Id}}, req.Id); result.Error != nil {
	//	return nil, status.Errorf(codes.NotFound, "商品不存在")
	//}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error){
	var goods model.Goods

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice= req.MarketPrice
	goods.ShopPrice= req.ShopPrice
	goods.GoodsBrief= req.GoodsBrief
	goods.ShipFree= req.ShipFree
	goods.Images= req.Images
	goods.DescImages= req.DescImages
	goods.GoodsFrontImage= req.GoodsFrontImage
	goods.IsNew= req.IsNew
	goods.IsHot= req.IsHot
	goods.OnSale= req.OnSale

	//global.DB.Save(&goods)
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}