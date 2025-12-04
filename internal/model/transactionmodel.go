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
	cacheTransactionIdPrefix = "cache:transaction:id:"
)

type (
	TransactionModel interface {
		Insert(data *Transaction) (sql.Result, error)
		FindOne(id int64) (*Transaction, error)
		FindByUserId(userId int64, limit int) ([]*Transaction, error)
	}

	defaultTransactionModel struct {
		sqlc.CachedConn
		table string
	}

	Transaction struct {
		Id            int64     `db:"id"`
		UserId        int64     `db:"user_id"`
		Type          string    `db:"type"`
		Currency      string    `db:"currency"`
		Amount        string    `db:"amount"`
		BalanceBefore string    `db:"balance_before"`
		BalanceAfter  string    `db:"balance_after"`
		Status        string    `db:"status"`
		TxHash        string    `db:"tx_hash"`
		CreatedAt     time.Time `db:"created_at"`
		UpdatedAt     time.Time `db:"updated_at"`
	}
)

func NewTransactionModel(conn sqlx.SqlConn, c cache.CacheConf) TransactionModel {
	cachedConn := sqlc.NewConn(conn, c)

	return &defaultTransactionModel{
		CachedConn: cachedConn,
		table:      "`transaction`",
	}
}

func (m *defaultTransactionModel) Insert(data *Transaction) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`user_id`, `type`, `currency`, `amount`, `balance_before`, `balance_after`, `status`, `tx_hash`) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table)
	ret, err := m.ExecNoCache(query, data.UserId, data.Type, data.Currency, data.Amount, data.BalanceBefore, data.BalanceAfter, data.Status, data.TxHash)
	return ret, err
}

func (m *defaultTransactionModel) FindOne(id int64) (*Transaction, error) {
	transactionIdKey := fmt.Sprintf("%s%v", cacheTransactionIdPrefix, id)
	var resp Transaction
	err := m.QueryRow(&resp, transactionIdKey, func(conn sqlx.SqlConn, v interface{}) error {
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

func (m *defaultTransactionModel) FindByUserId(userId int64, limit int) ([]*Transaction, error) {
	if limit <= 0 {
		limit = 50
	}
	query := fmt.Sprintf("select * from %s where `user_id` = ? order by `created_at` desc limit ?", m.table)
	var resp []*Transaction
	err := m.QueryRowsNoCache(&resp, query, userId, limit)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

