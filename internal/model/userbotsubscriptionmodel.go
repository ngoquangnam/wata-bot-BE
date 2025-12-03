package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	cacheSubscriptionIdPrefix   = "cache:subscription:id:"
	cacheSubscriptionUserPrefix = "cache:subscription:user:"
)

type (
	UserBotSubscriptionModel interface {
		Insert(data *UserBotSubscription) (sql.Result, error)
		FindOne(id int64) (*UserBotSubscription, error)
		FindByUserIdAndBotId(userId int64, botId string) (*UserBotSubscription, error)
		FindByUserId(userId int64) ([]*UserBotSubscription, error)
		Delete(id int64) error
		DeleteByUserIdAndBotId(userId int64, botId string) error
		CountByBotId(botId string) (int64, error)
	}

	defaultUserBotSubscriptionModel struct {
		sqlc.CachedConn
		table string
	}

	UserBotSubscription struct {
		Id        int64     `db:"id"`
		UserId    int64     `db:"user_id"`
		BotId     string    `db:"bot_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
)

func NewUserBotSubscriptionModel(conn sqlx.SqlConn, c cache.CacheConf) UserBotSubscriptionModel {
	cachedConn := sqlc.NewConn(conn, c)

	return &defaultUserBotSubscriptionModel{
		CachedConn: cachedConn,
		table:      "`user_bot_subscription`",
	}
}

func (m *defaultUserBotSubscriptionModel) Insert(data *UserBotSubscription) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`user_id`, `bot_id`) values (?, ?)", m.table)
	ret, err := m.ExecNoCache(query, data.UserId, data.BotId)
	return ret, err
}

func (m *defaultUserBotSubscriptionModel) FindOne(id int64) (*UserBotSubscription, error) {
	subscriptionIdKey := fmt.Sprintf("%s%v", cacheSubscriptionIdPrefix, id)
	var resp UserBotSubscription
	err := m.QueryRow(&resp, subscriptionIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select * from %s where `id` = ? limit 1", m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserBotSubscriptionModel) FindByUserIdAndBotId(userId int64, botId string) (*UserBotSubscription, error) {
	var resp UserBotSubscription
	query := fmt.Sprintf("select * from %s where `user_id` = ? and `bot_id` = ? limit 1", m.table)
	err := m.QueryRowNoCache(&resp, query, userId, botId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserBotSubscriptionModel) FindByUserId(userId int64) ([]*UserBotSubscription, error) {
	query := fmt.Sprintf("select * from %s where `user_id` = ? order by `created_at` desc", m.table)
	var resp []*UserBotSubscription
	err := m.QueryRowsNoCache(&resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserBotSubscriptionModel) Delete(id int64) error {
	subscriptionIdKey := fmt.Sprintf("%s%v", cacheSubscriptionIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, subscriptionIdKey)
	return err
}

func (m *defaultUserBotSubscriptionModel) DeleteByUserIdAndBotId(userId int64, botId string) error {
	query := fmt.Sprintf("delete from %s where `user_id` = ? and `bot_id` = ?", m.table)
	_, err := m.ExecNoCache(query, userId, botId)
	return err
}

func (m *defaultUserBotSubscriptionModel) CountByBotId(botId string) (int64, error) {
	var count int64
	query := fmt.Sprintf("select count(*) from %s where `bot_id` = ?", m.table)
	err := m.QueryRowNoCache(&count, query, botId)
	return count, err
}
