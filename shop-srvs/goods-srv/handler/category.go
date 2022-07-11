package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"shop-srvs/goods-srv/global"
	"shop-srvs/goods-srv/model"
	"shop-srvs/goods-srv/proto"
)

/*
这个接口如果只是想获取全部数据很简单，但可能前端希望我们整理好层级关系，就像这样：
    [
        {
            "id":xxx,
            "name":"",
            "level":1,
            "is_tab":false,
            "parent":13xxx,
            "sub_category":[
                "id":xxx,
                "name":"",
                "level":1,
                "is_tab":false,
                "sub_category":[]
            ]
        }
    ]
这就是Category表设计里加subCategory的原因，通过外键建立层级关系, 可以通过preloa加载subCategory。还有一个问题是因为category只用了一张表，所以用Find
查询的是所有的数据但我们希望的是列表里的是第一层的分类，其他二级三级的category都放在subcategory字段里，所以Find要加where条件指定leve=1，
然后要注意预加载只会加载一层需要指定加载subCategory的subcategory，所以应该用preload('subcategory.subcategory')
*/
func (s *GoodsServer) GetAllCategorysList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error){
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	b, _ := json.Marshal(&categorys)
	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

//获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListResponse := proto.SubCategoryListResponse{}

	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0{
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:       category.ID,
		Name: category.Name,
		Level: category.Level,
		IsTab: category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}

	var subCategorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse

	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id: subCategory.ID,
			Name: subCategory.Name,
			Level: subCategory.Level,
			IsTab: subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}

	categoryListResponse.SubCategorys = subCategoryResponse
	return &categoryListResponse, nil
}

func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}
	cMap := map[string]interface{}{}
	cMap["name"] = req.Name
	cMap["level"] = req.Level
	cMap["is_tab"] = req.IsTab
	if req.Level != 1 {
		//这里需要确认下父类目是否存在，可以我们确认，也可以交给调用端去做，我们默认已经确认过了
		cMap["parent_category_id"] = req.ParentCategory
	}
	tx := global.DB.Model(&model.Category{}).Create(cMap)
	fmt.Println(tx)
	return &proto.CategoryInfoResponse{Id:int32(category.ID)}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category

	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	global.DB.Save(&category)

	return &emptypb.Empty{}, nil
}
