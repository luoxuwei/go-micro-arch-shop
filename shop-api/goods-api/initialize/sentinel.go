package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"

	"shop-api/goods-api/global"
)

func InitSentinel(){
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf("初始化sentinel 异常: %v", err)
	}

	//配置限流规则
	sentinelInfo := global.ServerConfig.SentinelInfo
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               sentinelInfo.Resource,
			TokenCalculateStrategy: flow.TokenCalculateStrategy(sentinelInfo.Strategy),
			ControlBehavior:        flow.ControlBehavior(sentinelInfo.Behavior),
			Threshold:              sentinelInfo.Threshold,
			StatIntervalInMs:       sentinelInfo.Interval,
		},
	})

	if err != nil {
		zap.S().Fatalf("加载规则失败: %v", err)
	}
}