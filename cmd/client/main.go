package main

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	"github.com/itimky/faraway-test/pkg/client"
	logadapter "github.com/itimky/faraway-test/pkg/log/adapter"
	"github.com/itimky/faraway-test/pkg/pow"
)

type CryptoRand struct {
}

func (r CryptoRand) Read(p []byte) (int, error) {
	n, err := cryptorand.Read(p)
	if err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}

	return n, nil
}

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGTERM,
	)
	defer cancel()

	logger := logadapter.NewZeroLog()

	hashCash := pow.NewHashCash(CryptoRand{}, pow.DefaultDifficulty)

	bookClient := client.NewClient(hashCash)

	conn, err := net.Dial("tcp", "server:5678")
	if err != nil {
		logger.Errorf("dial: %s", err)

		return
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Errorf("close: %s", err)
		}
	}()

	quote, err := bookClient.GetRandomQuote(ctx, conn)
	if err != nil {
		logger.Errorf("get random quote: %s", err)

		return
	}

	logger.Infof("quote: %s", quote)
}
