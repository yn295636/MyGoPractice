package main

type GreetReq struct {
	Name string `json:"name" form:"name" binding:"required" example:"Peter"`
}

type GreetResp struct {
	Message string `json:"message" example:"hello world"`
}

type StoreInMongoReq map[string]interface {
}

type StoreInMongoResp struct {
	Result int32 `json:"result" example:"1"`
}

type ErrorRsp struct {
	Code    int32  `json:"code" example:"404" format:"int32"`
	Message string `json:"message" example:"It's an error message for developer"`
}
