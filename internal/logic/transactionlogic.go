package logic

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"wata-bot-BE/internal/model"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransactionLogic struct {
	logger logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionLogic {
	return &TransactionLogic{
		logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Deposit adds amount to user's balance
func (l *TransactionLogic) Deposit(req *types.DepositReq) (resp *types.TransactionResp, err error) {
	// Validate currency
	currency := strings.ToLower(req.Currency)
	if currency != "wata" && currency != "usdt" {
		return nil, model.NewAPIError(model.ErrCodeInvalidCurrency, model.ErrMsgInvalidCurrency)
	}

	// Validate amount
	amount, err := l.parseAmount(req.Amount)
	if err != nil || amount <= 0 {
		return nil, model.NewAPIError(model.ErrCodeInvalidAmount, model.ErrMsgInvalidAmount)
	}

	// Find user by address (no cache to get latest balance)
	user, err := l.svcCtx.UserModel.FindOneByAddressNoCache(req.Address)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.NewAPIError(model.ErrCodeFailedToFindUser, "User not found")
		}
		l.logger.Errorf("Failed to find user by address: %v", err)
		return nil, model.NewAPIError(model.ErrCodeFailedToFindUser, model.ErrMsgFailedToFindUser)
	}

	// Get current balance
	var balanceBefore, balanceAfter string
	if currency == "wata" {
		balanceBefore = user.WataBalance
		if balanceBefore == "" {
			balanceBefore = "0"
		}
		balanceAfter = l.addAmounts(balanceBefore, req.Amount)
		user.WataBalance = balanceAfter
	} else {
		balanceBefore = user.UsdtBalance
		if balanceBefore == "" {
			balanceBefore = "0"
		}
		balanceAfter = l.addAmounts(balanceBefore, req.Amount)
		user.UsdtBalance = balanceAfter
	}

	// Update user balance
	if err := l.svcCtx.UserModel.Update(user); err != nil {
		l.logger.Errorf("Failed to update user balance: %v", err)
		return nil, model.NewAPIError(model.ErrCodeFailedToUpdateBalance, model.ErrMsgFailedToUpdateBalance)
	}

	// Create transaction record
	transaction := &model.Transaction{
		UserId:        user.Id,
		Type:          "deposit",
		Currency:      currency,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Status:        "completed",
		TxHash:        req.TxHash,
	}

	result, err := l.svcCtx.TransactionModel.Insert(transaction)
	if err != nil {
		l.logger.Errorf("Failed to create transaction record: %v", err)
		// Don't fail the deposit if transaction record fails, but log it
	} else {
		transactionId, _ := result.LastInsertId()
		transaction.Id = transactionId
	}

	// Return response
	transactionData := types.TransactionData{
		Type:          transaction.Type,
		Currency:      transaction.Currency,
		Amount:        transaction.Amount,
		BalanceBefore: transaction.BalanceBefore,
		BalanceAfter:  transaction.BalanceAfter,
		Status:        transaction.Status,
		TxHash:        transaction.TxHash,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}

	return &types.TransactionResp{
		Message: "Deposit successful",
		Data:    transactionData,
	}, nil
}

// Withdraw subtracts amount from user's balance
func (l *TransactionLogic) Withdraw(req *types.WithdrawReq) (resp *types.TransactionResp, err error) {
	// Validate currency
	currency := strings.ToLower(req.Currency)
	if currency != "wata" && currency != "usdt" {
		return nil, model.NewAPIError(model.ErrCodeInvalidCurrency, model.ErrMsgInvalidCurrency)
	}

	// Validate amount
	amount, err := l.parseAmount(req.Amount)
	if err != nil || amount <= 0 {
		return nil, model.NewAPIError(model.ErrCodeInvalidAmount, model.ErrMsgInvalidAmount)
	}

	// Find user by address (no cache to get latest balance)
	user, err := l.svcCtx.UserModel.FindOneByAddressNoCache(req.Address)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, model.NewAPIError(model.ErrCodeFailedToFindUser, "User not found")
		}
		l.logger.Errorf("Failed to find user by address: %v", err)
		return nil, model.NewAPIError(model.ErrCodeFailedToFindUser, model.ErrMsgFailedToFindUser)
	}

	// Get current balance
	var balanceBefore, balanceAfter string
	if currency == "wata" {
		balanceBefore = user.WataBalance
		if balanceBefore == "" {
			balanceBefore = "0"
		}
		// Check sufficient balance
		if !l.hasSufficientBalance(balanceBefore, req.Amount) {
			return nil, model.NewAPIError(model.ErrCodeInsufficientBalance, model.ErrMsgInsufficientBalance)
		}
		balanceAfter = l.subtractAmounts(balanceBefore, req.Amount)
		user.WataBalance = balanceAfter
	} else {
		balanceBefore = user.UsdtBalance
		if balanceBefore == "" {
			balanceBefore = "0"
		}
		// Check sufficient balance
		if !l.hasSufficientBalance(balanceBefore, req.Amount) {
			return nil, model.NewAPIError(model.ErrCodeInsufficientBalance, model.ErrMsgInsufficientBalance)
		}
		balanceAfter = l.subtractAmounts(balanceBefore, req.Amount)
		user.UsdtBalance = balanceAfter
	}

	// Update user balance
	if err := l.svcCtx.UserModel.Update(user); err != nil {
		l.logger.Errorf("Failed to update user balance: %v", err)
		return nil, model.NewAPIError(model.ErrCodeFailedToUpdateBalance, model.ErrMsgFailedToUpdateBalance)
	}

	// Create transaction record
	transaction := &model.Transaction{
		UserId:        user.Id,
		Type:          "withdraw",
		Currency:      currency,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Status:        "completed",
		TxHash:        req.TxHash,
	}

	result, err := l.svcCtx.TransactionModel.Insert(transaction)
	if err != nil {
		l.logger.Errorf("Failed to create transaction record: %v", err)
		// Don't fail the withdraw if transaction record fails, but log it
	} else {
		transactionId, _ := result.LastInsertId()
		transaction.Id = transactionId
	}

	// Return response
	transactionData := types.TransactionData{
		Type:          transaction.Type,
		Currency:      transaction.Currency,
		Amount:        transaction.Amount,
		BalanceBefore: transaction.BalanceBefore,
		BalanceAfter:  transaction.BalanceAfter,
		Status:        transaction.Status,
		TxHash:        transaction.TxHash,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}

	return &types.TransactionResp{
		Message: "Withdraw successful",
		Data:    transactionData,
	}, nil
}

// Helper functions for amount calculations
func (l *TransactionLogic) parseAmount(amountStr string) (float64, error) {
	amountStr = strings.TrimSpace(amountStr)
	return strconv.ParseFloat(amountStr, 64)
}

func (l *TransactionLogic) addAmounts(balanceStr, amountStr string) string {
	balance, err1 := l.parseAmount(balanceStr)
	amount, err2 := l.parseAmount(amountStr)
	if err1 != nil || err2 != nil {
		return balanceStr // Return original if parse fails
	}
	result := balance + amount
	return l.formatAmount(result)
}

func (l *TransactionLogic) subtractAmounts(balanceStr, amountStr string) string {
	balance, err1 := l.parseAmount(balanceStr)
	amount, err2 := l.parseAmount(amountStr)
	if err1 != nil || err2 != nil {
		return balanceStr // Return original if parse fails
	}
	result := balance - amount
	if result < 0 {
		result = 0
	}
	return l.formatAmount(result)
}

func (l *TransactionLogic) hasSufficientBalance(balanceStr, amountStr string) bool {
	balance, err1 := l.parseAmount(balanceStr)
	amount, err2 := l.parseAmount(amountStr)
	if err1 != nil || err2 != nil {
		return false
	}
	return balance >= amount
}

func (l *TransactionLogic) formatAmount(amount float64) string {
	// Format to remove trailing zeros but keep up to 8 decimal places
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.8f", amount), "0"), ".")
}
