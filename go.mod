module github.com/yn295636/MyGoPractice

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/agiledragon/gomonkey v0.0.0-20180923140027-0ba3ddf4a9d5
	github.com/alicebob/miniredis/v2 v2.13.2
	github.com/appleboy/gofight/v2 v2.0.0
	github.com/coreos/bbolt v1.3.5 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dankinder/httpmock v1.0.1
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/lint v0.0.0-20180702182130-06c8688daad7 // indirect
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.8 // indirect
	github.com/h2non/gock v1.0.15
	github.com/jonboulle/clockwork v0.2.0 // indirect
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/mongodb/mongo-go-driver v0.3.0
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/spf13/cobra v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.5.1
	github.com/tidwall/pretty v1.0.2 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.1-0.20180714160509-73f8eece6fdc // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/grpc v1.31.1
	google.golang.org/grpc/examples v0.0.0-20200902210233-8630cac324bf // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/h2non/gock v1.0.15 => gopkg.in/h2non/gock.v1 v1.0.15

replace github.com/coreos/bbolt v1.3.5 => go.etcd.io/bbolt v1.3.5

replace github.com/coreos/etcd v3.3.25+incompatible => go.etcd.io/etcd v3.3.25+incompatible
