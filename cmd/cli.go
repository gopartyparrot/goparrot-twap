package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-co-op/gocron"
	twapConfig "github.com/gopartyparrot/goparrot-twap/config"
	"github.com/gopartyparrot/goparrot-twap/swap"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CliArgs struct {
	RPCUrl            string        `arg:"required,env" help:"rpc url"`
	RPCWs             string        `arg:"required,env" help:"rpc websocket"`
	WalletPK          string        `arg:"required,env,--wallet" help:"wallet private key"`
	StorePath         string        `arg:"env" help:"store successful swaps logs" default:"./logs/swaps.json"`
	Interval          string        `arg:"required,--interval" help:"run interval in time units (s, m, h)"`
	Pair              string        `arg:"required,--pair" help:"pair"`
	Side              swap.SwapSide `arg:"--side" help:"side of the swap can be buy or sell (default buy)" default:"buy"`
	Amount            float64       `arg:"required,--amount" help:"amount to buy or sell"`
	StopAmount        float64       `arg:"--stopAmount" help:"amount ro reach" default:"999999999999999"`
	TransferAddress   string        `arg:"--transferAddress" help:"address to transfer the balance when above the TransferThreshold"`
	TransferThreshold float64       `arg:"--transferThreshold" help:"threshold for transfer all balance to TransferAddress"`
}

func run() error {
	err := godotenv.Load()
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("loading environment: %w", err)
		}
	}

	var args CliArgs
	arg.MustParse(&args)

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		return fmt.Errorf("can't initialize zap logger: %w", err)
	}
	defer logger.Sync()
	logger.Info("using RPC",
		zap.String("http", args.RPCUrl),
	)

	s := gocron.NewScheduler(time.UTC)

	clientRPC := rpc.New(args.RPCUrl)

	swapper, err := swap.NewTokenSwapper(swap.TokenSwapperConfig{
		ClientRPC:  clientRPC,
		RPCWs:      args.RPCWs,
		PrivateKey: args.WalletPK,
		StorePath:  args.StorePath,
		Logger:     logger,
		Tokens:     twapConfig.GetTokens(),
		Pools:      twapConfig.GetPools(),
	})
	if err != nil {
		logger.Fatal("create swapper", zap.Error(err))
		return err
	}

	err = swapper.Init(
		context.Background(),
		args.Pair,
		args.Side,
		args.Amount,
		args.StopAmount,
		args.TransferAddress,
		args.TransferThreshold,
	)
	if err != nil {
		logger.Fatal("init swapper", zap.Error(err))
		return err
	}

	s.Every(args.Interval).Do(swapper.Start)

	s.StartBlocking()

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln("run error %w", err)
	}
}
