package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
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

	sellDetail := model.StockSellDetail{
		OrderSn: req.OrderSn,
		Status:  1,
	}
	var details []model.GoodsDetail

	for _, goodInfo := range req.GoodsInfo {

		details = append(details, model.GoodsDetail{
			Goods: goodInfo.GoodsId,
			Num: goodInfo.Num,
		})

		mutex := global.RedisSync.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

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
		if ok, err := mutex.Unlock(); !ok || err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	sellDetail.Detail = details
	if result := tx.Create(&sellDetail); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "保存库存扣减历史失败")
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

func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	for i := range msgs {
		//确保幂等性，防止消息的重复发送导致一个订单的库存归还多次，还要防止没有扣减的库存被归还
		//新建一张表，记录详细的订单扣减细节，以及归还细节。就是订单的哪件商品扣了多少库存
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Errorf("解析json失败： %v\n", msgs[i].Body)
			return consumer.ConsumeSuccess, nil
		}

		//将库存加回去 将selldetail的status设置为2
		tx := global.DB.Begin()
		var sellDetail model.StockSellDetail
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn:orderInfo.OrderSn, Status:1}).First(&sellDetail); result.RowsAffected == 0 {
			//没查到，表示已归还
			return consumer.ConsumeSuccess, nil
		}

		//如果查询到那么逐个归还库存
		for _, orderGood := range sellDetail.Detail {
			if result := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods:orderGood.Goods}).Update("stocks", gorm.Expr("stocks+?", orderGood.Num)); result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}

		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn:orderInfo.OrderSn}).Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}