package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"shop-srvs/inventory-srv/global"
	"shop-srvs/inventory-srv/model"
	"shop-srvs/inventory-srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}


func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	//设置库存， 如果我要更新库存
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods:req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods:req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num: inv.Stocks,
	}, nil
}


func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {

	//对订单里的商品库存进行减扣，要么全部成功，要么全部失败，所以需要用到事务，目前还只是本机事务，最终要用分布式事务
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods:goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		//减扣操作是典型的并发问题，需要保证并发安全，不然会出现数据不一致的问题
		//数据不一致问题，不是事务解决的问题，是锁，最终要用分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
	}
    tx.Commit()
    return &emptypb.Empty{}, nil
}

func (s *InventoryServer) Reback(c context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//要么全部成功，要么全部失败
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods:goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		//减扣操作是典型的并发问题，需要保证并发安全，不然会出现数据不一致的问题
		//数据不一致问题，不是事务解决的问题，是锁，最终要用分布式锁
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}