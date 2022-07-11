package handler

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"strings"
	"time"

	"context"

	"shop-srvs/goods-srv/global"
	"shop-srvs/goods-srv/model"
	"shop-srvs/goods-srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedUserServer
}

