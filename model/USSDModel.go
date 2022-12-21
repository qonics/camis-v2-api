package model

type MTNUSSDModel struct {
	Msisdn    string `form:"msisdn" binding:"required"`
	Input     string `form:"input" binding:"required"`
	SessionId string `form:"sessionId" binding:"required"`
	Operator  string `form:"operator" binding:"required"`
	Country   string `form:"country" binding:"required"`
}

type AfricasTalkingUSSDModel struct {
	Msisdn      string `form:"phoneNumber" binding:"required"`
	Input       string `form:"text" binding:"required"`
	SessionId   string `form:"sessionId" binding:"required"`
	Operator    string `form:"networkCode" binding:"required"`
	ServiceCode string `form:"serviceCode" binding:"required"`
}
