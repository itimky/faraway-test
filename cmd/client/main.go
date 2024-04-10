package main

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/itimky/faraway-test/pkg/client"
	logadapter "github.com/itimky/faraway-test/pkg/log/adapter"
	"github.com/itimky/faraway-test/pkg/pow"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/req"
	_ "go.nanomsg.org/mangos/v3/transport/tcp"
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

	var sock mangos.Socket

	sock, err := req.NewSocket()
	if err != nil {
		logger.Errorf("new socket: %s", err)

		return
	}

	defer func() {
		err := sock.Close()
		if err != nil {
			logger.Errorf("close socket: %s", err)

			return
		}
	}()

	if err := sock.Dial("tcp://server:5678"); err != nil {
		logger.Errorf("dial: %s", err)

		return
	}

	bookClient := client.NewClient(sock, hashCash)

	quote, err := bookClient.GetRandomQuote(ctx)
	if err != nil {
		logger.Errorf("get random quote: %s", err)

		return
	}

	logger.Infof("quote: %s", quote)
}
