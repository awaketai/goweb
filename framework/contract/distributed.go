package contract

import "time"

const DistributedKey = "web:distributed"

// Distributed 分布式服务
type Distributed interface {
	// Select
	// serviceName 服务名称
	// appID 当前AppID
	// holdTime 分布式选择器hold住的时间
	// selectAppId 分布式选择器最终选择的App
	Select(serviceName string, appID string, holdTime time.Duration) (selectAppID string, err error)
}
