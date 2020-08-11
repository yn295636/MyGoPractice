package sample_api

type MultiplyReq struct {
	A int32 `json:"a" form:"a" binding:"required" example:"2"`
	B int32 `json:"b" form:"b" binding:"required" example:"3"`
}

type MultiplyResp struct {
	Result int64 `json:"result" binding:"required" example:"6"`
}

type ErrorRsp struct {
	Code    int32  `json:"code" example:"404" format:"int32"`
	Message string `json:"message" example:"It's an error message for developer"`
}