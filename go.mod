module github.com/yn295636/MyGoPractice

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/agiledragon/gomonkey v0.0.0-20180923140027-0ba3ddf4a9d5
	github.com/alicebob/miniredis/v2 v2.13.2
	github.com/appleboy/gofight/v2 v2.0.0
	github.com/dankinder/httpmock v1.0.1
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/h2non/gock v1.0.15
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/mongodb/mongo-go-driver v0.3.0
	github.com/nsqio/go-nsq v1.0.8
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/pretty v1.0.2 // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.1-0.20180714160509-73f8eece6fdc // indirect
	go.etcd.io/etcd/api/v3 v3.5.5
	go.etcd.io/etcd/client/v3 v3.5.5
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.26.0
)

replace github.com/h2non/gock v1.0.15 => gopkg.in/h2non/gock.v1 v1.0.15
