package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"shop-srvs/inventory-srv/global"
	"shop-srvs/inventory-srv/model"
	"shop-srvs/inventory-srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedGoodsServer
}

