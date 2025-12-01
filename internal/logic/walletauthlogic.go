package logic

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"wata-bot-BE/internal/model"
	"wata-bot-BE/internal/svc"
	"wata-bot-BE/internal/types"
	"wata-bot-BE/internal/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

type WalletAuthLogic struct {
	logger logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWalletAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WalletAuthLogic {
	return &WalletAuthLogic{
		logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WalletAuthLogic) WalletAuth(req *types.WalletAuthReq) (resp *types.WalletAuthResp, err error) {
	// Verify signature
	address, err := l.verifySignature(req.Message, req.Signature)
	if err != nil {
		l.logger.Errorf("Signature verification failed: %v", err)
		utils.WriteErrorLog("Signature verification failed", err)
		return nil, model.NewAPIError(model.ErrCodeInvalidSignature, model.ErrMsgInvalidSignature)
	}

	addressStr := address.String()

	// Get or create user (allow registration if not found)
	user, err := l.getOrCreateUser(addressStr, req.InviteCode)
	if err != nil {
		return nil, err
	}

	// Generate JWT tokens
	accessToken, refreshToken, expiresIn, err := l.generateTokens(addressStr, user.ReferralCode)
	if err != nil {
		l.logger.Errorf("Token generation failed: %v", err)
		utils.WriteErrorLog("Token generation failed", err)
		return nil, model.NewAPIError(model.ErrCodeTokenGenerationFailed, model.ErrMsgTokenGenerationFailed)
	}

	return &types.WalletAuthResp{
		Message: "success",
		Data: types.WalletAuthData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expiresIn,
			WataReward:   user.WataReward,
			Role:         user.Role,
		},
	}, nil
}

func (l *WalletAuthLogic) WalletAuthNotSign(req *types.WalletAuthNotSignReq) (resp *types.WalletAuthResp, err error) {
	// Validate and normalize address using go-ethereum
	addressStr, err := l.validateAndNormalizeAddress(req.Address)
	if err != nil {
		return nil, err
	}

	// Get or create user with referral_code from request (allow registration if not found)
	user, err := l.getOrCreateUserWithReferralCode(addressStr, req.ReferralCode)
	if err != nil {
		return nil, err
	}

	// Generate JWT tokens
	accessToken, refreshToken, expiresIn, err := l.generateTokens(addressStr, user.ReferralCode)
	if err != nil {
		l.logger.Errorf("Token generation failed: %v", err)
		utils.WriteErrorLog("Token generation failed", err)
		return nil, model.NewAPIError(model.ErrCodeTokenGenerationFailed, model.ErrMsgTokenGenerationFailed)
	}

	return &types.WalletAuthResp{
		Message: "success",
		Data: types.WalletAuthData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expiresIn,
			WataReward:   user.WataReward,
			Role:         user.Role,
		},
	}, nil
}

// getOrCreateUser gets existing user or creates new one if not found (allows registration)
func (l *WalletAuthLogic) getOrCreateUser(addressStr, inviteCode string) (*model.User, error) {
	referralCode := strings.ToUpper(addressStr[len(addressStr)-8:])

	// Check if user exists
	user, err := l.svcCtx.UserModel.FindOneByAddress(addressStr)
	if err != nil && err != model.ErrNotFound {
		l.logger.Errorf("Database error: %v", err)
		utils.WriteErrorLog("Database error when finding user by address", err)
		return nil, model.NewAPIError(model.ErrCodeDatabaseError, model.ErrMsgDatabaseError)
	}

	// If user not found, create new user (allow registration)
	if err == model.ErrNotFound {
		// New user registration - only save address and referral_code
		newUser := &model.User{
			Address:      addressStr,
			ReferralCode: referralCode,
		}
		_, err = l.svcCtx.UserModel.Insert(newUser)
		if err != nil {
			l.logger.Errorf("Failed to create user: %v", err)
			utils.WriteErrorLog("Failed to create user", err)
			return nil, model.NewAPIError(model.ErrCodeFailedToCreateUser, model.ErrMsgFailedToCreateUser)
		}

		// Retrieve the newly created user with retry (for cache consistency)
		// Retry up to 3 times with small delay between attempts
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			user, err = l.svcCtx.UserModel.FindOneByAddress(addressStr)
			if err == nil {
				break
			}
			if i < maxRetries-1 {
				// Small delay before retry (for cache consistency)
				time.Sleep(50 * time.Millisecond)
			}
		}

		// If still not found after retries, this is an edge case (cache not updated yet)
		// Since Insert succeeded, we can construct user object to continue
		if err != nil {
			l.logger.Errorf("Failed to find created user after insert and retries (address: %s): %v", addressStr, err)
			utils.WriteErrorLog("Failed to find created user after insert", err)
			// User was created successfully, so we construct user object to continue
			// This handles edge case where cache hasn't updated yet but DB insert succeeded
			user = &model.User{
				Address:      addressStr,
				ReferralCode: referralCode,
				WataReward:   0,
				Role:         "user",
			}
			l.logger.Infof("Using constructed user object for address: %s (insert succeeded but query failed)", addressStr)
		}
		l.logger.Infof("New user registered: %s", addressStr)
	}
	// Note: No longer updating invite_code for existing users

	return user, nil
}

// getOrCreateUserWithReferralCode gets existing user or creates new one with referral_code from request
func (l *WalletAuthLogic) getOrCreateUserWithReferralCode(addressStr, referralCode string) (*model.User, error) {
	// Validate and normalize referral_code
	referralCode = strings.TrimSpace(referralCode)
	referralCode = strings.ToUpper(referralCode)

	// If referral_code is empty, generate from address
	if referralCode == "" {
		referralCode = strings.ToUpper(addressStr[len(addressStr)-8:])
	}

	// Check if user exists
	user, err := l.svcCtx.UserModel.FindOneByAddress(addressStr)
	if err != nil && err != model.ErrNotFound {
		l.logger.Errorf("Database error: %v", err)
		utils.WriteErrorLog("Database error when finding user by address", err)
		return nil, model.NewAPIError(model.ErrCodeDatabaseError, model.ErrMsgDatabaseError)
	}

	// If user not found, create new user (allow registration)
	if err == model.ErrNotFound {
		// New user registration - only save address and referral_code
		newUser := &model.User{
			Address:      addressStr,
			ReferralCode: referralCode,
		}
		_, err = l.svcCtx.UserModel.Insert(newUser)
		if err != nil {
			l.logger.Errorf("Failed to create user: %v", err)
			utils.WriteErrorLog("Failed to create user", err)
			return nil, model.NewAPIError(model.ErrCodeFailedToCreateUser, model.ErrMsgFailedToCreateUser)
		}

		// Retrieve the newly created user with retry (for cache consistency)
		// Retry up to 3 times with small delay between attempts
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			user, err = l.svcCtx.UserModel.FindOneByAddress(addressStr)
			if err == nil {
				break
			}
			if i < maxRetries-1 {
				// Small delay before retry (for cache consistency)
				time.Sleep(50 * time.Millisecond)
			}
		}

		// If still not found after retries, this is an edge case (cache not updated yet)
		// Since Insert succeeded, we can construct user object to continue
		if err != nil {
			l.logger.Errorf("Failed to find created user after insert and retries (address: %s): %v", addressStr, err)
			utils.WriteErrorLog("Failed to find created user after insert", err)
			// User was created successfully, so we construct user object to continue
			// This handles edge case where cache hasn't updated yet but DB insert succeeded
			user = &model.User{
				Address:      addressStr,
				ReferralCode: referralCode,
				WataReward:   0,
				Role:         "user",
			}
			l.logger.Infof("Using constructed user object for address: %s (insert succeeded but query failed)", addressStr)
		}
		l.logger.Infof("New user registered: %s with referral_code: %s", addressStr, referralCode)
	}

	return user, nil
}

// validateAndNormalizeAddress validates and normalizes Ethereum address using go-ethereum
func (l *WalletAuthLogic) validateAndNormalizeAddress(addressInput string) (string, error) {
	// Trim whitespace
	addressInput = strings.TrimSpace(addressInput)
	if addressInput == "" {
		l.logger.Errorf("Empty address provided")
		utils.WriteErrorLog("Empty address", fmt.Errorf("address is empty"))
		return "", model.NewAPIError(model.ErrCodeInvalidAddressFormat, model.ErrMsgInvalidAddressFormat)
	}

	// Basic format validation: must start with 0x
	if !strings.HasPrefix(strings.ToLower(addressInput), "0x") {
		l.logger.Errorf("Invalid address format: %s (missing 0x prefix)", addressInput)
		utils.WriteErrorLog("Invalid address format", fmt.Errorf("address missing 0x prefix: %s", addressInput))
		return "", model.NewAPIError(model.ErrCodeInvalidAddressFormat, model.ErrMsgInvalidAddressFormat)
	}

	// Normalize to lowercase for format check
	addressLower := strings.ToLower(addressInput)

	// Validate hex characters (after 0x prefix)
	hexPart := addressLower[2:]
	for _, char := range hexPart {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			l.logger.Errorf("Invalid address format: %s (contains invalid hex characters)", addressInput)
			utils.WriteErrorLog("Invalid address format", fmt.Errorf("address contains invalid hex characters"))
			return "", model.NewAPIError(model.ErrCodeInvalidAddressFormat, model.ErrMsgInvalidAddressFormat)
		}
	}

	// Use go-ethereum's HexToAddress to parse and validate
	// HexToAddress accepts both checksum and lowercase addresses
	address := common.HexToAddress(addressInput)

	// Check if address is zero address (0x0000...)
	zeroAddress := common.HexToAddress("0x0")
	if address == zeroAddress {
		l.logger.Errorf("Zero address not allowed: %s", addressInput)
		utils.WriteErrorLog("Zero address", fmt.Errorf("zero address is not allowed"))
		return "", model.NewAPIError(model.ErrCodeInvalidAddressFormat, model.ErrMsgInvalidAddressFormat)
	}

	// Return checksum address (Hex() returns EIP-55 checksummed address)
	normalizedAddress := address.Hex()

	// Verify the normalized address is valid (should always be 42 chars)
	if len(normalizedAddress) != 42 {
		l.logger.Errorf("Invalid address normalization: %s -> %s", addressInput, normalizedAddress)
		utils.WriteErrorLog("Invalid address normalization", fmt.Errorf("normalized address has invalid length"))
		return "", model.NewAPIError(model.ErrCodeInvalidAddressFormat, model.ErrMsgInvalidAddressFormat)
	}

	return normalizedAddress, nil
}

// verifySignature verifies the Ethereum signature
func (l *WalletAuthLogic) verifySignature(message, signature string) (common.Address, error) {
	// Remove 0x prefix if present
	sig := strings.TrimPrefix(signature, "0x")
	if len(sig) < 128 {
		return common.Address{}, errors.New("invalid signature length")
	}

	// Decode signature
	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to decode signature: %v", err)
	}

	// Ethereum signature recovery
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	// Hash message with Ethereum prefix
	msgHash := crypto.Keccak256Hash(
		[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)),
	)

	// Recover public key
	pubKey, err := crypto.SigToPub(msgHash.Bytes(), sigBytes)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to recover public key: %v", err)
	}

	// Get address from public key
	address := crypto.PubkeyToAddress(*pubKey)
	return address, nil
}

// generateTokens generates JWT access token and refresh token
func (l *WalletAuthLogic) generateTokens(address, inviteCode string) (string, string, int64, error) {
	// JWT secret key from config
	secretKey := []byte(l.svcCtx.Config.JWTSecret)

	// Expires in 1 year (31536000 seconds)
	expiresIn := int64(31536000)
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// Generate referral code from address (last 8 chars)
	referralCode := strings.ToUpper(address[len(address)-8:])

	// Create JWT claims
	claims := jwt.MapClaims{
		"aud":           address,
		"exp":           expiresAt.Unix(),
		"iat":           time.Now().Unix(),
		"iss":           "prod-aibot-backend-issuer",
		"sub":           "auth",
		"user_id":       fmt.Sprintf("%x", crypto.Keccak256Hash([]byte(address)).Bytes()[:16]),
		"address":       address,
		"referral_code": referralCode,
		"role":          "user",
	}

	// Generate access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", "", 0, err
	}

	// Generate refresh token (base64 encoded)
	refreshTokenBytes := crypto.Keccak256Hash([]byte(address + time.Now().String())).Bytes()
	refreshToken := base64.StdEncoding.EncodeToString(refreshTokenBytes)

	return accessToken, refreshToken, expiresIn, nil
}
