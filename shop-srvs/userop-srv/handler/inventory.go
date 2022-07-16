package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"shop-srvs/userop-srv/global"
	"shop-srvs/userop-srv/model"
	"shop-srvs/userop-srv/proto"
)

type UseropServer struct {
	proto.UnimplementedInventoryServer
}


