package api

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

type StoreInRedisReq struct {
	Key   string `json:"key" form:"key" binding:"required" example:"name"`
	Value string `json:"value" form:"value" binding:"required" example:"Peter"`
}

type StoreInRedisResp struct {
	Result int32 `json:"result" example:"1"`
}

type StoreUserInDbReq struct {
	Username    string `json:"username" form:"username" binding:"required" example:"peter123"`
	Gender      int8   `json:"gender" form:"gender" binding:"required" example:"1"`
	Age         int8   `json:"age" form:"age" binding:"required" example:"18"`
	ExternalUid int32  `json:"external_uid" form:"external_uid" binding:"required" example:"1"`
}

type StoreUserInDbResp struct {
	Uid int64 `json:"uid" example:"1"`
}

type StoreUserPhoneInDbReq struct {
	Phone string `json:"phone" form:"phone" binding:"required" example:"+6512345678"`
}

type StoreUserPhoneInDbResp struct {
}

type ErrorRsp struct {
	Code    int32  `json:"code" example:"404" format:"int32"`
	Message string `json:"message" example:"It's an error message for developer"`
}
