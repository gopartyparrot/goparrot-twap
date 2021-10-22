package swap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gopartyparrot/goparrot-twap/config"
	"github.com/gopartyparrot/goparrot-twap/store"
	"go.uber.org/zap"
)

var (
	ErrSwapPoolNotFound     = errors.New("swap pool not found for given pair")
	ErrUpdateBalances       = errors.New("failed to update wallet balances")
	ErrFromBalanceNotEnough = errors.New("from balance not enough for swap")
	ErrTargetAmountReached  = errors.New("target balance amount reached")
)

type SwapSide string

const (
	SwapSide_Buy  SwapSide = "buy"
	SwapSide_Sell SwapSide = "sell"
)

type SwapStatus struct {
	TxID    string
	Pair    string
	Date    string
	Side    SwapSide
	Amount  uint64
	ErrLogs string `json:",omitempty"`
}

type SwapTaskConfig struct {
	pair         string
	side         SwapSide
	amount       float64
	targetAmount float64
	pool         config.PoolConfig
}

type TokenSwapperConfig struct {
	RPCEndpoint string
	PrivateKey  string
	StorePath   string
	Tokens      map[string]config.TokenInfo
	Pools       map[string]config.PoolConfig
	Logger      *zap.Logger
}

type TokenSwapper struct {
	clientRPC     *rpc.Client
	store         *store.JSONStore
	account       solana.PrivateKey
	logger        *zap.Logger
	raydiumSwap   *RaydiumSwap
	tokens        map[string]config.TokenInfo
	pools         map[string]config.PoolConfig
	tokenBalances map[string]uint64
	tokenAccounts map[string]solana.PublicKey
	swapTask      SwapTaskConfig
}

func (s *TokenSwapper) Init(
	ctx context.Context,
	pair string,
	side SwapSide,
	amount float64,
	targetAmount float64,
) error {

	s.swapTask = SwapTaskConfig{
		pair:         pair,
		side:         side,
		amount:       amount,
		targetAmount: targetAmount,
	}

	for k, v := range s.pools {
		if k == pair {
			s.swapTask.pool = v
		}
	}
	if s.swapTask.pool.FromToken == "" {
		return ErrSwapPoolNotFound
	}

	mints := []solana.PublicKey{
		solana.MustPublicKeyFromBase58(s.swapTask.pool.FromToken),
		solana.MustPublicKeyFromBase58(s.swapTask.pool.ToToken),
	}

	existingAccounts, missingAccounts, err := GetTokenAccountsFromMints(ctx, *s.clientRPC, s.account.PublicKey(), mints...)
	if err != nil {
		return err
	}

	if len(missingAccounts) != 0 {
		instrs := []solana.Instruction{}
		for mint, _ := range missingAccounts {
			if mint == config.NativeSOL {
				continue
			}
			s.logger.Info("need to create token account", zap.String("mint", mint))
			inst, err := associatedtokenaccount.NewCreateInstruction(
				s.account.PublicKey(),
				s.account.PublicKey(),
				solana.MustPublicKeyFromBase58(mint),
			).ValidateAndBuild()
			if err != nil {
				return err
			}
			instrs = append(instrs, inst)
		}
		sig, err := ExecuteInstructions(ctx, s.clientRPC, []solana.PrivateKey{s.account}, instrs...)
		if err != nil {
			return err
		}
		s.logger.Info("missing token accounts created", zap.String("txID", sig.String()))
		for k, v := range missingAccounts {
			existingAccounts[k] = v
		}
	}
	s.tokenAccounts = existingAccounts
	return nil
}

func (s *TokenSwapper) UpdateBalances(ctx context.Context) error {
	pks := []solana.PublicKey{}
	for _, v := range s.tokenAccounts {
		pks = append(pks, v)
	}
	res, err := GetTokenAccountsBalance(ctx, *s.clientRPC, pks...)
	if err != nil {
		return err
	}
	for address, amount := range res {
		s.tokenBalances[address] = amount
	}
	return nil
}

func (s *TokenSwapper) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	err := s.UpdateBalances(ctx)
	if err != nil {
		return ErrUpdateBalances
	}

	fromToken := s.swapTask.pool.FromToken
	toToken := s.swapTask.pool.ToToken

	if s.swapTask.side == SwapSide_Sell {
		fromToken = s.swapTask.pool.ToToken
		toToken = s.swapTask.pool.FromToken
	}

	fromAddress := s.tokenAccounts[fromToken]
	fromBalance := s.tokenBalances[fromAddress.String()]
	fromTokenInfo := s.tokens[fromToken]

	toAddress := s.tokenAccounts[toToken]
	toBalance := s.tokenBalances[toAddress.String()]
	toTokenInfo := s.tokens[toToken]

	amount := fromTokenInfo.FromFloat(s.swapTask.amount)
	targetAmount := toTokenInfo.FromFloat(s.swapTask.targetAmount)

	if toBalance > targetAmount {
		s.logger.Info("target amount reached for swap "+fromTokenInfo.Symbol+" to "+toTokenInfo.Symbol,
			zap.Uint64("targetAmount", targetAmount),
			zap.Uint64("currentBalance", toBalance),
		)
		return ErrTargetAmountReached
	}

	if amount > fromBalance {
		s.logger.Warn("not enough balance to swap "+fromTokenInfo.Symbol+" to "+toTokenInfo.Symbol,
			zap.Uint64("swapAmount", amount),
			zap.Uint64("currentBalance", fromBalance),
		)
		return ErrFromBalanceNotEnough
	}

	sig, err := s.raydiumSwap.Swap(
		ctx,
		&s.swapTask.pool.RaydiumPoolConfig,
		amount,
		fromToken,
		fromAddress,
		toToken,
		toAddress,
	)

	status := SwapStatus{
		Date:   time.Now().UTC().Format(time.UnixDate),
		Pair:   s.swapTask.pair,
		Side:   s.swapTask.side,
		Amount: amount,
	}
	if err != nil {
		s.logger.Warn("swap fail", zap.Error(err))
		status.ErrLogs = fmt.Sprintf("error: %v", err)
	} else {
		s.logger.Info("swap success", zap.String("txID", sig.String()))
		status.TxID = sig.String()
	}
	key := fmt.Sprintf("%s_%s", status.Pair, status.Date)
	s.store.Set(key, status)

	return nil
}

func NewTokenSwapper(cfg TokenSwapperConfig) (*TokenSwapper, error) {
	clientRPC := rpc.New(cfg.RPCEndpoint)

	store, err := store.OpenJSONStore(cfg.StorePath)
	if err != nil {
		return nil, err
	}

	privateKey, err := solana.PrivateKeyFromBase58(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	raydiumSwap := RaydiumSwap{
		clientRPC: clientRPC,
		account:   privateKey,
	}

	l := TokenSwapper{
		clientRPC:     clientRPC,
		store:         store,
		logger:        cfg.Logger,
		pools:         cfg.Pools,
		tokens:        cfg.Tokens,
		account:       privateKey,
		raydiumSwap:   &raydiumSwap,
		tokenBalances: map[string]uint64{},
	}

	return &l, nil
}
