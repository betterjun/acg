module {{ .Name }}

go 1.13

require (
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc
	github.com/ethereum/go-ethereum v1.9.24
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.1
	github.com/mattn/go-sqlite3 v2.0.1+incompatible // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.5.1
	go.etcd.io/bbolt v1.3.2
	go.uber.org/zap v1.15.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.0.3
	gorm.io/driver/postgres v1.0.5
	gorm.io/driver/sqlite v1.1.3
	gorm.io/gorm v1.20.5
	nanomsg.org/go/mangos/v2 v2.0.8
)
