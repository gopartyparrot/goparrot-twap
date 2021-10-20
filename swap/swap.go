package swap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/gopartyparrot/goparrot_buy/config"
	"github.com/gopartyparrot/goparrot_buy/json"
	"go.uber.org/zap"
)

var (
	ErrUpdateBalances       = errors.New("failed to update wallet balances")
	ErrFromBalanceNotEnough = errors.New("from balance not enough for swap")
	ErrTargetAmountReached  = errors.New("target balance amount reached")
)

type TokenSwapperConfig struct {
	RPCEndpoint string
	WSEndpoint  string
	PrivateKey  string
	StorePath   string
	Logger      *zap.Logger
}

type SwapParams struct {
	Reverse      bool    `json:"reverse,omitempty"`
	BuyAmount    float64 `json:"buy_amount,omitempty"`
	TargetAmount float64 `json:"target_amount,omitempty"`
}

type SwapStatus struct {
	TxID              string
	Pool              string
	Date              string
	BuyAmount         uint64
	PreBalanceAmount  uint64
	PostBalanceAmount uint64
	ErrLogs           string `json:",omitempty"`
}

type TokenSwapper struct {
	clientRPC     *rpc.Client
	clientWS      *ws.Client
	store         *json.JSONStore
	params        *json.JSONStore
	account       solana.PrivateKey
	logger        *zap.Logger
	raydiumSwap   *RaydiumSwap
	tokenBalances map[string]uint64
	tokenAccounts map[string]solana.PublicKey
}

func (s *TokenSwapper) Init(ctx context.Context) error {
	mints := []solana.PublicKey{}
	for _, c := range config.PoolConfigs {
		mints = append(mints, solana.MustPublicKeyFromBase58(c.FromToken))
		mints = append(mints, solana.MustPublicKeyFromBase58(c.ToToken))
	}

	existingAccounts, missingAccounts, err := GetTokenAccountsFromMints(ctx, *s.clientRPC, s.account.PublicKey(), mints...)
	if err != nil {
		return err
	}

	if len(missingAccounts) != 0 {
		instrs := []solana.Instruction{}
		for mint, _ := range missingAccounts {
			s.logger.Info("-- Need token account for mint", zap.String("mint", mint))
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
		log.Println("-- Missing token accounts created in txID:", sig)
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

	for name, c := range config.PoolConfigs {

		var params SwapParams
		ok, err := s.params.Get(name, &params)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		fromToken := c.FromToken
		toToken := c.ToToken
		if params.Reverse {
			fromToken = c.ToToken
			toToken = c.FromToken
		}

		fromAddress := s.tokenAccounts[fromToken]
		fromBalance := s.tokenBalances[fromAddress.String()]
		fromTokenInfo := config.TokenInfos[fromToken]

		toAddress := s.tokenAccounts[toToken]
		toBalance := s.tokenBalances[toAddress.String()]
		toTokenInfo := config.TokenInfos[toToken]

		buyAmount := fromTokenInfo.FromFloat(params.BuyAmount)
		targetAmount := toTokenInfo.FromFloat(params.TargetAmount)

		if toBalance > targetAmount {
			s.logger.Info("reach max amount for swap "+fromTokenInfo.Symbol+" to "+toTokenInfo.Symbol,
				zap.Uint64("balance", toBalance),
				zap.Uint64("target", targetAmount),
			)
			return ErrTargetAmountReached
		}

		if buyAmount > fromBalance {
			s.logger.Warn("not enough balance to swap "+fromTokenInfo.Symbol+" to "+toTokenInfo.Symbol,
				zap.Uint64("balance", fromBalance),
				zap.Uint64("buyAmount", buyAmount),
			)
			return ErrFromBalanceNotEnough
		}

		sig, err := s.raydiumSwap.Swap(
			ctx,
			&c.RaydiumPoolConfig,
			buyAmount,
			fromToken,
			fromAddress,
			toToken,
			toAddress,
		)

		status := SwapStatus{
			Pool:              name,
			Date:              time.Now().UTC().Format(time.UnixDate),
			BuyAmount:         buyAmount,
			PreBalanceAmount:  fromBalance,
			PostBalanceAmount: fromBalance + buyAmount,
		}
		if err != nil {
			s.logger.Fatal("swap failled", zap.Error(err))
			status.ErrLogs = fmt.Sprintf("error: %v", err)
		} else {
			s.logger.Info("swap success", zap.String("txID", sig.String()))
			status.TxID = sig.String()
		}
		key := fmt.Sprintf("%s:%s", name, time.Now().UTC().Format(time.UnixDate))
		s.store.Set(key, status)
	}

	return nil
}

func NewTokenSwapper(cfg TokenSwapperConfig) (*TokenSwapper, error) {

	// TODO: change to yaml
	params, err := json.OpenJSONStore("params.json")
	if err != nil {
		return nil, err
	}

	store, err := json.OpenJSONStore(cfg.StorePath)
	if err != nil {
		return nil, err
	}

	clientRPC := rpc.New(cfg.RPCEndpoint)

	clientWS, err := ws.Connect(context.Background(), cfg.WSEndpoint)
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
		clientWS:      clientWS,
		store:         store,
		logger:        cfg.Logger,
		params:        params,
		account:       privateKey,
		raydiumSwap:   &raydiumSwap,
		tokenBalances: map[string]uint64{},
	}

	return &l, nil
}
