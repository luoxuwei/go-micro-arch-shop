package pay

import (
	"context"
	"github.com/gin-gonic/gin"
    "github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"net/http"
	"shop-api/order-api/global"
	"shop-api/order-api/proto"
)

func Notify(ctx *gin.Context) {
	//第三个参数表示是不是生产环境，会影响网关，开发时必须设置为FALSE，生成的支付url就会使用沙箱的支付网关
	client, err := alipay.New(global.ServerConfig.AliPayInfo.AppID, global.ServerConfig.AliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey((global.ServerConfig.AliPayInfo.AliPublicKey))
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":err.Error(),
		})
		return
	}

	//获取订单状态，失败了可能是验证失败，是伪造的请求
	noti, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	//更新订单支付状态
	_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: noti.OutTradeNo, //这个订单号，是我们自己生成的，在生成支付url的时候带过去的。TradeNo是支付宝的订单号。
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	//这样就行了，支付宝收到就能确认。
	ctx.String(http.StatusOK, "success")
}