package main

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/rand"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/itimky/faraway-test/pkg/book"
	logadapter "github.com/itimky/faraway-test/pkg/log/adapter"
	"github.com/itimky/faraway-test/pkg/pow"
	"github.com/itimky/faraway-test/pkg/server"
)

func quotes() []string {
	return []string{
		`“I like to listen. I have learned a great deal from listening carefully. Most people never listen.”
― Ernest Hemingway`,
		`“Guard well your thoughts when alone and your words when accompanied.”
― Roy T. Bennett`,
		`“Quit being so hard on yourself. We are what we are; we love what we love.
We don't need to justify it to anyone... not even to ourselves.”
― Scott Lynch, The Republic of Thieves`,
		`“Voice is not just the sound that comes from your throat, but the feelings that come from your words.”
― Jennifer Donnelly, A Northern Light`,
	}
}

type CryptoRand struct {
}

func (r CryptoRand) Read(p []byte) (int, error) {
	n, err := cryptorand.Read(p)
	if err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}

	return n, nil
}

func main() { //nolint:funlen
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
	quotesBook := book.NewBook(rand.New(rand.NewSource(time.Now().UnixNano())), quotes()) //nolint:gosec

	powMiddleware := server.NewPOWMiddleware(hashCash)
	handler := server.NewHandler(quotesBook)

	listener, err := net.Listen("tcp", ":5678") //nolint:gosec
	if err != nil {
		logger.Errorf("listen: %s", err)

		return
	}

	go func() {
		<-ctx.Done()

		err := listener.Close()
		if err != nil {
			logger.Errorf("listener close: %s", err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() == nil {
				logger.Errorf("accept: %s", err)
			}

			return
		}

		go func(conn net.Conn) {
			defer func() {
				err := conn.Close()
				if err != nil {
					logger.Errorf("conn close: %s", err)
				}
			}()

			err := powMiddleware.Handle(ctx, conn)
			if err != nil {
				logger.Errorf("pow middleware handle: %s", err)

				return
			}

			err = handler.Handle(ctx, conn)
			if err != nil {
				logger.Errorf("handler handle: %s", err)
			}
		}(conn)
	}
}
