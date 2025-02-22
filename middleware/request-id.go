package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/songquanpeng/one-api/common/helper"
)

func RequestId() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		id := helper.GenRequestID(ctx)
		c.Set(helper.RequestIdKey, id)
		ctx = helper.SetRequestID(ctx, id)
		c.Request = c.Request.WithContext(ctx)
		c.Header(helper.RequestIdKey, id)
		c.Next()
	}
}
