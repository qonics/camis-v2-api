package controller

import (
	"camis-v2-api/model"
	"errors"

	"github.com/gin-gonic/gin"
)

func getUSSDdata(c *gin.Context) (*model.AfricasTalkingUSSDModel, error) {
	var data model.AfricasTalkingUSSDModel
	if err := c.ShouldBind(data); err != nil {
		return nil, errors.New("invalid data found " + err.Error())
	}
	return &data, nil
}
func ExecuteUSSD(c *gin.Context) {
	USSDData, err := getUSSDdata(c)
	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "USSD session: " + USSDData.SessionId})
}
