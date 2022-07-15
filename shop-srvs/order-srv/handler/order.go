package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"shop-srvs/order-srv/global"
	"shop-srvs/order-srv/model"
	"shop-srvs/order-srv/proto"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

func (*OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
    var shopCarts []model.ShoppingCart

    result := global.DB.Where(&model.OrderInfo{User: req.Id}).Find(&shopCarts)
    if result.Error != nil {
        return nil, result.Error
	}

    rsp := proto.CartItemListResponse{
    	Total: int32(result.RowsAffected),
	}

	for _, shopCar := range shopCarts {
		rsp.Data = append(rsp.Data, &proto.ShopCartInfoResponse{
			Id: shopCar.ID,
			UserId: shopCar.User,
			GoodsId: shopCar.Goods,
			Nums: shopCar.Nums,
			Checked: shopCar.Checked,
		})
	}

	return &rsp, nil
}

func (*OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//将商品添加到购物车:
	//    1. 购物车中原本没有这件商品, 新建一个记录.
	//    2. 这个商品之前添加到了购物车, 合并，只把商品数加上
	var shopCart model.ShoppingCart

	if result := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart); result.RowsAffected == 1 {
		//如果已经存在，只需加上商品数，更新操作
		shopCart.Nums += req.Nums
	}else{
		//插入操作
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

//购物车中选中和加减商品数操作调用的接口，更新购物车记录，更新数量和选中状态
func (*OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	var shopCart model.ShoppingCart

	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&shopCart); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	shopCart.Checked = req.Checked
	//如果前端只进行选中操作，更新的只有checked状态，num就是默认值0，所以这里要判断一下，不为0才更新num
	//如果购物车中商品数减为0了，前端判断一下调用删除接口就行了
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	global.DB.Save(&shopCart)

	return &emptypb.Empty{}, nil
}

func (*OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{}); result.RowsAffected == 0{
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}