package data

import (
	"context"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/omalloc/contrib/kratos/orm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/omalloc/kratos-layout/internal/biz"
	"github.com/omalloc/kratos-layout/internal/conf"
)

var driverMap = map[string]func(ds string) gorm.Dialector{
	"sqlite": sqlite.Open,
	"mysql":  mysql.Open,
}

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	dialector, ok := driverMap[c.Database.Driver]
	if !ok {
		panic("unsupported database driver")
	}

	db, err := orm.New(
		orm.WithDriver(dialector(c.Database.Source)),
		orm.WithTracingOpts(orm.WithDatabaseName("test")), // <- change me.
		orm.WithLogger(
			orm.WithDebug(),
			orm.WIthSlowThreshold(time.Second*2),
			orm.WithSkipCallerLookup(true),
			orm.WithSkipErrRecordNotFound(true),
			orm.WithLogHelper(log.NewFilter(log.GetLogger(), log.FilterLevel(log.LevelDebug))),
		),
	)

	// open database failed.
	if err != nil {
		panic(err)
	}

	// 自动创建表结构，如果不用自动建表，就删掉这行
	_ = db.Session(&gorm.Session{SkipHooks: true}).AutoMigrate(&biz.Greeter{})

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		if db != nil {
			if sql, _ := db.DB(); sql != nil {
				_ = sql.Close()
			}
		}
	}

	return &Data{
		db: db,
	}, cleanup, nil
}

func (d *Data) Check(ctx context.Context) error {
	return nil
}
