module obwallet

go 1.20

require (
	github.com/emersion/go-sasl v0.0.0-20220912192320-0145f2c60ead
	github.com/emersion/go-smtp v0.18.0
	github.com/lightningnetwork/lnd/cert v1.1.0
	github.com/lithammer/shortuuid/v4 v4.0.0
	github.com/mojocn/base64Captcha v1.3.5
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/gorm v1.25.3
)

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/image v0.0.0-20190501045829-6d32002ffd75 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
)

replace github.com/lightningnetwork/lnd/cert => ../lnd/cert
