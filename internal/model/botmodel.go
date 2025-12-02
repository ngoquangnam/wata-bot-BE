package model

import (
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	cacheBotIdPrefix = "cache:bot:id:"
)

type (
	BotModel interface {
		Insert(data *Bot) (sql.Result, error)
		FindOne(id string) (*Bot, error)
		FindAll() ([]*Bot, error)
		FindAllActive() ([]*Bot, error)
		Update(data *Bot) error
		Delete(id string) error
	}

	defaultBotModel struct {
		sqlc.CachedConn
		table string
	}

	Bot struct {
		Id                   string  `db:"id"`
		Name                 string  `db:"name"`
		IconLetter           string  `db:"icon_letter"`
		RiskLevel            string  `db:"risk_level"`
		DurationDays         int     `db:"duration_days"`
		ExpectedReturnPercent int    `db:"expected_return_percent"`
		AprDisplay           string  `db:"apr_display"`
		MinInvestment        int     `db:"min_investment"`
		MaxInvestment        int     `db:"max_investment"`
		InvestmentRange      string  `db:"investment_range"`
		Subscribers          int     `db:"subscribers"`
		Author               string  `db:"author"`
		Description          string  `db:"description"`
		IsActive             bool    `db:"is_active"`
		LockupPeriod         string  `db:"lockup_period"`
		ExpectedReturn       string  `db:"expected_return"`
		MinInvestmentDisplay string  `db:"min_investment_display"`
		MaxInvestmentDisplay string  `db:"max_investment_display"`
		Roi30d               string  `db:"roi30d"`
		WinRate              string  `db:"win_rate"`
		TradingPair          string  `db:"trading_pair"`
		TotalTrades          int     `db:"total_trades"`
		Pnl30d               float64 `db:"pnl30d"`
	}
)

func NewBotModel(conn sqlx.SqlConn, c cache.CacheConf) BotModel {
	cachedConn := sqlc.NewConn(conn, c)

	return &defaultBotModel{
		CachedConn: cachedConn,
		table:      "`bot`",
	}
}

func (m *defaultBotModel) Insert(data *Bot) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (`id`, `name`, `icon_letter`, `risk_level`, `duration_days`, `expected_return_percent`, `apr_display`, `min_investment`, `max_investment`, `investment_range`, `subscribers`, `author`, `description`, `is_active`, `lockup_period`, `expected_return`, `min_investment_display`, `max_investment_display`, `roi30d`, `win_rate`, `trading_pair`, `total_trades`, `pnl30d`) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table)
	
	ret, err := m.ExecNoCache(query,
		data.Id, data.Name, data.IconLetter, data.RiskLevel, data.DurationDays,
		data.ExpectedReturnPercent, data.AprDisplay, data.MinInvestment, data.MaxInvestment,
		data.InvestmentRange, data.Subscribers, data.Author, data.Description, data.IsActive,
		data.LockupPeriod, data.ExpectedReturn, data.MinInvestmentDisplay, data.MaxInvestmentDisplay,
		data.Roi30d, data.WinRate, data.TradingPair, data.TotalTrades, data.Pnl30d,
	)
	return ret, err
}

func (m *defaultBotModel) FindOne(id string) (*Bot, error) {
	botIdKey := fmt.Sprintf("%s%v", cacheBotIdPrefix, id)
	var resp Bot
	err := m.QueryRow(&resp, botIdKey, func(conn sqlx.SqlConn, v interface{}) error {
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

func (m *defaultBotModel) FindAll() ([]*Bot, error) {
	query := fmt.Sprintf("select * from %s order by `id`", m.table)
	var resp []*Bot
	err := m.QueryRowsNoCache(&resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultBotModel) FindAllActive() ([]*Bot, error) {
	query := fmt.Sprintf("select * from %s where `is_active` = 1 order by `id`", m.table)
	var resp []*Bot
	err := m.QueryRowsNoCache(&resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultBotModel) Update(data *Bot) error {
	botIdKey := fmt.Sprintf("%s%v", cacheBotIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `name`=?, `icon_letter`=?, `risk_level`=?, `duration_days`=?, `expected_return_percent`=?, `apr_display`=?, `min_investment`=?, `max_investment`=?, `investment_range`=?, `subscribers`=?, `author`=?, `description`=?, `is_active`=?, `lockup_period`=?, `expected_return`=?, `min_investment_display`=?, `max_investment_display`=?, `roi30d`=?, `win_rate`=?, `trading_pair`=?, `total_trades`=?, `pnl30d`=? where `id` = ?", m.table)
		return conn.Exec(query,
			data.Name, data.IconLetter, data.RiskLevel, data.DurationDays, data.ExpectedReturnPercent,
			data.AprDisplay, data.MinInvestment, data.MaxInvestment, data.InvestmentRange, data.Subscribers,
			data.Author, data.Description, data.IsActive, data.LockupPeriod, data.ExpectedReturn,
			data.MinInvestmentDisplay, data.MaxInvestmentDisplay, data.Roi30d, data.WinRate,
			data.TradingPair, data.TotalTrades, data.Pnl30d, data.Id,
		)
	}, botIdKey)
	return err
}

func (m *defaultBotModel) Delete(id string) error {
	botIdKey := fmt.Sprintf("%s%v", cacheBotIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, botIdKey)
	return err
}

