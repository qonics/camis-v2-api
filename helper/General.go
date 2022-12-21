package helper

import (
	"camis-v2-api/config"
	"camis-v2-api/model"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()
var SessionExpirationTime time.Duration = 1800
var RedisPrefix string = "CAMIS_"
var CachePrefix string = "CAMIS_CACHE_"

func SecurePath(c *gin.Context) *model.UserPayload {
	token := c.GetHeader("Authorization")

	var userData model.UserPayload

	token = RedisPrefix + token
	// fmt.Println("TOKEN: ", token)
	client := []byte(config.Redis.Get(ctx, token).Val())
	if client == nil || len(string(client)) == 0 {
		c.JSON(401, gin.H{"message": "Token not found or expired", "status": 401})
		panic("Token not found or expired")
	}
	// fmt.Println("User data:", string(client))
	var logger model.UserPayload
	err := json.Unmarshal(client, &logger)
	if err != nil {
		c.JSON(401, gin.H{"message": "Authentication failed, invalid token", "status": 401})
		panic("done, secure path failed #unmarshal" + err.Error())
	}
	userAgent := c.Request.UserAgent()
	// userIp := c.ClientIP()
	if len(c.GetHeader("uag")) > 0 {
		userAgent = c.GetHeader("uag")
	}
	if logger.Uag != userAgent {
		//destroy this token, it is altered
		config.Redis.Del(ctx, token)
		c.JSON(401, gin.H{"message": "Authentication failed, invalid token", "status": 401})
		panic("done, secure path failed #unmarshal" + err.Error())
	}
	// if len(c.GetHeader("ip")) > 0 {
	// 	userIp = c.GetHeader("ip")
	// }

	//check if it is current active token
	activeToken := string([]byte(config.Redis.Get(ctx, RedisPrefix+"user_"+logger.Uid+"_active_token").Val()))
	if token != RedisPrefix+activeToken {
		//destroy this token, it is not the current
		config.Redis.Del(ctx, token)
		c.JSON(401, gin.H{"message": "Your account has be signed in on other computer", "status": 401})
		panic("Your account has be signed in on other computer:" + activeToken + " - " + token)
	}
	config.Redis.Expire(ctx, token, time.Duration(SessionExpirationTime*time.Minute))
	return &userData
}
func RemoveCachedItem(key string) {
	config.Redis.Del(ctx, CachePrefix+key)
}

func CompressJsonIndexing(data string) string {
	data = strings.ReplaceAll(data, "[", "")
	data = strings.ReplaceAll(data, "]", "")
	data = strings.ReplaceAll(data, "{", "")
	data = strings.ReplaceAll(data, "}", "")
	data = strings.ReplaceAll(data, ",", "")
	data = strings.ReplaceAll(data, ":", "")
	data = strings.ReplaceAll(data, " ", "")
	data = strings.ReplaceAll(data, "\"", "")
	fmt.Println("simplified string: ", data)
	return data
}
