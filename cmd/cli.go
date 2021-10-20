package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/go-co-op/gocron"
	"github.com/gopartyparrot/goparrot_buy/swap"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CliArgs struct {
	RPCUrl    string `arg:"required,env" help:"rpc url"`
	WSUrl     string `arg:"required,env" help:"ws url"`
	WalletPK  string `arg:"required,env,-w" help:"wallet private key"`
	StorePath string `arg:"env" help:"store successful swaps infos" default:"./logs/swaps.json"`
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
		log.Panicln("can't initialize zap logger: %w", err)
	}
	defer logger.Sync()
	logger.Info("Using RPC",
		zap.String("http", args.RPCUrl),
		zap.String("ws", args.WSUrl),
	)

	s := gocron.NewScheduler(time.UTC)

	swapper, err := swap.NewTokenSwapper(swap.TokenSwapperConfig{
		RPCEndpoint: args.RPCUrl,
		WSEndpoint:  args.WSUrl,
		PrivateKey:  args.WalletPK,
		StorePath:   args.StorePath,
		Logger:      logger,
	})
	if err != nil {
		return err
	}

	err = swapper.Init(context.Background())
	if err != nil {
		return err
	}

	s.Every(5).Seconds().Do(swapper.Start)

	s.StartBlocking()

	return nil
}

func main() {
	err := run()

	if err != nil {
		log.Fatalln(err)
	}
}
