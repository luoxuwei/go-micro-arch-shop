package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"shop-srvs/order-srv/global"
	"shop-srvs/order-srv/model"
	"shop-srvs/order-srv/proto"
	"time"
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

func (*OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp proto.OrderListResponse

	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)

	//分页
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &rsp, nil
}

func (*OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp proto.OrderInfoDetailResponse

	//需要检查一下权限，确认这个订单的id是否是当前用户的订单，这是必须的，比如有可能是爬虫在爬取订单数据。
	//如果在web层在查询订单详情时，应该先查询一下订单id是否是当前用户的，但这样需要提供一个检查这个订单id是不是这个用户的接口，
	//底层服务可以简单的做一下，web层把订单id和用户id一起传过来，
	//在个人中心可以这样做，但是如果是后台管理系统，web层如果是后台管理系统 那么只传递order的id，如果是电商系统还需要一个用户的id
	//在底层可以不管，gorm 的 where会忽略是默认值的字段
	if result := global.DB.Where(&model.OrderInfo{BaseModel:model.BaseModel{ID:req.Id}, User:req.UserId}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	orderInfo := proto.OrderInfoResponse{}
	orderInfo.Id = order.ID
	orderInfo.UserId = order.User
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.Post = order.Post
	orderInfo.Total = order.OrderMount
	orderInfo.Address = order.Address
	orderInfo.Name = order.SignerName
	orderInfo.Mobile = order.SingerMobile

	rsp.OrderInfo = &orderInfo

	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order:order.ID}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}

	for _, orderGood := range orderGoods {
		rsp.Goods = append(rsp.Goods, &proto.OrderItemResponse{
			GoodsId: orderGood.Goods,
			GoodsName: orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			GoodsImage: orderGood.GoodsImage,
			Nums: orderGood.Nums,
		})
	}

	return &rsp, nil
}

func (*OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {

	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	//获取购物车中选中的商品, 有可能用户没有选中任何商品就点击结算，所以需要判断一下result
	if result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
        return nil, status.Errorf(codes.InvalidArgument, "没有选中的商品")
	}

	goodsNumMap := make(map[int32]int32)
	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumMap[shopCart.Goods] = shopCart.Nums
	}

	//从商品微服务中拉取商品列表，需要商品信息中的价格算订单金额，和后续插入订单商品信息
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	goods, err := global.GoodsClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "批量查询商品信息失败")
	}

	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods: good.Id,
			GoodsName: good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums: goodsNumMap[good.Id],
		})

		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num: goodsNumMap[good.Id],
		})
	}

	//扣减库存
	if _, err = global.InvertoryClient.Sell(context.Background(), &proto.SellInfo{OrderSn:orderInfo.OrderSn, GoodsInfo: goodsInvInfo}); err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败")
	}

	tx := global.DB.Begin()
	order := model.OrderInfo{
		OrderSn: GenerateOrderSn(req.UserId),
		Address: req.Address,
		SignerName: req.Name,
		SingerMobile: req.Mobile,
		Post: req.Post,
		User: req.UserId,
	}
	if result := global.DB.Save(&order); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建订单失败")
	}

	//批量插入订单商品信息
	for _, orderGood := range orderGoods {
		orderGood.Order = order.ID
	}
	if result := global.DB.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建订单失败")
	}

	if result := global.DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建订单失败")
	}
	tx.Commit()
    return &proto.OrderInfoResponse{Id: order.ID, OrderSn: order.OrderSn, Total: order.OrderMount}, nil
}

//订单号的生成规则, 仅仅有时间戳不够，高并发情况下可能多个请求拿到的时间戳一样，
//加上用户id已经有一定保证，同一个用户不会有那么高的并发，再加上两位随机数，就更有保障
func GenerateOrderSn(userId int32) string{
	/*
	   年月日时分秒+用户id+2位随机数
	*/
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
	return orderSn
}

func (*OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &emptypb.Empty{}, nil
}