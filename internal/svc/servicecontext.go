package svc

import (
	"wata-bot-BE/internal/config"
	"wata-bot-BE/internal/model"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Database.DataSource)

	// Use cache if configured, otherwise use empty cache config (cache disabled)
	cacheConf := c.Cache
	if len(cacheConf) == 0 {
		// Create empty cache config to disable caching
		// sqlc.NewConn requires a non-nil cache config, but empty slice works
		cacheConf = make([]cache.NodeConf, 0)
	}

	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlConn, cacheConf),
	}
}
