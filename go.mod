module it.elfire.ru/lekovr/grpcdemo

go 1.12

replace SELF => ./

require (
	SELF v0.0.0-00010101000000-000000000000
	github.com/birkirb/loggers-mapper-logrus v0.0.0-20180326232643-461f2d8e6f72
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/go-pg/pg v8.0.4+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/jessevdk/go-flags v1.4.0
	github.com/jinzhu/gorm v1.9.10
	github.com/lib/pq v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.2.2
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107
	google.golang.org/grpc v1.22.1
	gopkg.in/birkirb/loggers.v1 v1.1.0
)
