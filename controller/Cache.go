package controller

import (
	"camis-v2-api/config"
	"camis-v2-api/helper"
	"camis-v2-api/model"
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var ctx = context.Background()
var MOMO_PERSONAL_INFORMATION_URL = "/mot/mm/getaccountholderpersonalinformation"
var MOMO_PERSONAL_IDENTIFICATION_URL = "/mot/mm/getaccountholderidentification"
var MOMO_DEBIT = "/mot/mm/debit"
var MOMO_CREDIT = "/mot/mm/sptransfer"
var MOMO_TRANSACTION_STATUS = "/mot/mm/gettransactionstatus"
var SP_ACCOUNT = "mopay.sp"

const (
	Pending           = 201
	Processing        = 202
	Success           = 200
	Not_enough_fund   = 301
	Account_not_fund  = 404
	Unauthorized      = 401
	Invalid_user_data = 400
	System_error      = 500
)

func SocketConnection() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/paymentMonitor", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

}

/*
Receive deleteCache request
*/
func DeleteCache(c *gin.Context) {
	helper.SecurePath(c)
	go deleteCacheProccess(c)
	c.JSON(200, gin.H{"status": 200,
		"message": "Delete cache worker is processing",
	})
}

/*
Receive addCache request
*/
func AddCache(c *gin.Context) {
	helper.SecurePath(c)
	var data model.CacheModel
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"status": 400, "message": "Invalid data" + err.Error()})
		return
	}
	go addCacheProccess(data)
	c.JSON(200, gin.H{"status": 200,
		"message": "Add cache worker is processing",
	})
}

/*
Delete cache based on item id search
*/
func deleteCacheProccess(c *gin.Context) {
	helper.SecurePath(c)
	id := c.Param("id")
	var (
		url        string
		cacheCount uint
	)
	iter := config.SESSION.Query("select url from cache where data like '%" + id + "%'").Iter()
	for iter.Scan(&url) {
		decoded, err := base64.StdEncoding.DecodeString(url)
		if err != nil {
			helper.Warning("cache delete failed, base64 decoding: "+err.Error(), nil)
		}
		helper.RemoveCachedItem(url)
		err = config.SESSION.Query("delete from cache where url = ?", url).Exec()
		if err != nil {
			helper.Warning("cache delete failed: "+err.Error(), nil)
		}
		fmt.Println("Cache: ", string(decoded))
		url = ""
		cacheCount++
	}
	if err := iter.Close(); err != nil {
		helper.Warning("cache delete: "+err.Error(), nil)
		// panic(err.Error())
	}
}
func addCacheProccess(data model.CacheModel) {
	data.Url = base64.StdEncoding.EncodeToString([]byte((data.Url)))
	newData, err := base64.StdEncoding.DecodeString(data.Data)
	if err != nil {
		helper.Warning("cache saving failed, unable to decode data: "+err.Error(), nil)
		// panic(err.Error())
	}
	compressedIndexData := helper.CompressJsonIndexing(string(newData))
	err = config.SESSION.Query("insert into cache (url,data,created_at,updated_at) values (?,?,toTimestamp(now()),toTimestamp(now()))", data.Url, compressedIndexData).Exec()
	if err != nil {
		helper.Warning("cache saving failed: Url:"+data.Url+","+err.Error(), nil)
		// panic(err.Error())
	}
	var minutes time.Duration = time.Duration(data.Duration)
	var ctx = context.Background()
	config.Redis.Set(ctx, helper.CachePrefix+data.Url, data.Data, minutes*time.Minute)
}

func ServiceStatusCheck(c *gin.Context) {
	c.JSON(400, gin.H{"status": 200, "message": "CA-MIS API service is running"})
}
