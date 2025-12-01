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
	cacheUserIdPrefix      = "cache:user:id:"
	cacheUserAddressPrefix = "cache:user:address:"
)

type (
	UserModel interface {
		Insert(data *User) (sql.Result, error)
		FindOne(id int64) (*User, error)
		FindOneByAddress(address string) (*User, error)
		Update(data *User) error
		Delete(id int64) error
	}

	defaultUserModel struct {
		sqlc.CachedConn
		table string
	}

	User struct {
		Id           int64     `db:"id"`
		Address      string    `db:"address"`
		ReferralCode string    `db:"referral_code"`
		InviteCode   string    `db:"invite_code"` // Can be NULL in DB, use COALESCE in queries
		WataReward   int       `db:"wata_reward"`
		Role         string    `db:"role"`
		CreatedAt    time.Time `db:"created_at"`
		UpdatedAt    time.Time `db:"updated_at"`
	}
)

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	// Create CachedConn - if cache config is empty, it will work without caching
	// Empty cache config is valid and won't cause "no cache nodes" error
	cachedConn := sqlc.NewConn(conn, c)

	return &defaultUserModel{
		CachedConn: cachedConn,
		table:      "`user`",
	}
}

func (m *defaultUserModel) Insert(data *User) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`address`, `referral_code`) values (?, ?)", m.table)
	ret, err := m.ExecNoCache(query, data.Address, data.ReferralCode)
	return ret, err
}

func (m *defaultUserModel) FindOne(id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRow(&resp, userIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select `id`, `address`, `referral_code`, COALESCE(`invite_code`, '') as `invite_code`, `wata_reward`, `role`, `created_at`, `updated_at` from %s where `id` = ? limit 1", m.table)
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

func (m *defaultUserModel) FindOneByAddress(address string) (*User, error) {
	userAddressKey := fmt.Sprintf("%s%v", cacheUserAddressPrefix, address)
	var resp User
	err := m.QueryRowIndex(&resp, userAddressKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select `id`, `address`, `referral_code`, COALESCE(`invite_code`, '') as `invite_code`, `wata_reward`, `role`, `created_at`, `updated_at` from %s where `address` = ? limit 1", m.table)
		if err := conn.QueryRow(&resp, query, address); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) Update(data *User) error {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	userAddressKey := fmt.Sprintf("%s%v", cacheUserAddressPrefix, data.Address)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `address`=?, `referral_code`=?, `invite_code`=?, `wata_reward`=?, `role`=? where `id` = ?", m.table)
		return conn.Exec(query, data.Address, data.ReferralCode, data.InviteCode, data.WataReward, data.Role, data.Id)
	}, userIdKey, userAddressKey)
	return err
}

func (m *defaultUserModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	userAddressKey := fmt.Sprintf("%s%v", cacheUserAddressPrefix, data.Address)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, userIdKey, userAddressKey)
	return err
}

func (m *defaultUserModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
}

func (m *defaultUserModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select `id`, `address`, `referral_code`, COALESCE(`invite_code`, '') as `invite_code`, `wata_reward`, `role`, `created_at`, `updated_at` from %s where `id` = ? limit 1", m.table)
	return conn.QueryRow(v, query, primary)
}
