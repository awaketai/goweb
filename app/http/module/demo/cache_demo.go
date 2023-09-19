package demo

import (
	"fmt"
	"github.com/awaketai/goweb/framework/contract"
	"github.com/awaketai/goweb/framework/gin"
	"time"
)

func (api *DemoApi) DemoCache(c *gin.Context) {
	logger := c.MustMake(contract.LogKey).(contract.Log)
	logger.Info(c, "request start", nil)
	cacheService := c.MustMake(contract.CacheKey).(contract.Cache)
	// currently,only basic types can be stored
	err := cacheService.Set(c, "foo", 1.23, 1*time.Hour)
	fmt.Println("ser err:", err)
	if err != nil {

		c.AbortWithError(500, err)
		return
	}
	val, err := cacheService.Get(c, "foo")
	fmt.Println("val-3:", val)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	logger.Info(c, "cache get", map[string]any{
		"val": val,
	})
	if err := cacheService.Del(c, "foo"); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, "ok")

}
