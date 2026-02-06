module github.com/sonramos/microservices/payment

go 1.25.1

require (
	github.com/sonramos/microservices-proto/golang/payment v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.78.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/sonramos/microservices-proto/golang/payment => ../../microservices-proto/github/sonramos/microservices-proto/golang/payment
