package main

import (
	"context"
	cryptorand "crypto/rand"
	"log"
	"math/rand"
	"os/signal"
	"syscall"
	"time"

	"github.com/itimky/faraway-test/pkg/book"
	"github.com/itimky/faraway-test/pkg/pow"
	"github.com/itimky/faraway-test/pkg/server"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/rep"
	// register transports
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

// randomizer is a function that implements io.Reader and uses a crypto rand
// to generate random bytes.

type CryptoRand struct {
}

func (r CryptoRand) Read(p []byte) (n int, err error) {
	return cryptorand.Read(p)
}

func main() {
	hashCash := pow.NewHashCash(CryptoRand{}, pow.DefaultDifficulty)

	quotesBook := book.NewBook(rand.New(rand.NewSource(time.Now().UnixNano())), []string{
		`“I like to listen. I have learned a great deal from listening carefully. Most people never listen.”
― Ernest Hemingway`,
		`“Guard well your thoughts when alone and your words when accompanied.”
― Roy T. Bennett`,
		`“Quit being so hard on yourself. We are what we are; we love what we love. We don't need to justify it to anyone... not even to ourselves.”
― Scott Lynch, The Republic of Thieves`,
		`“Voice is not just the sound that comes from your throat, but the feelings that come from your words.”
― Jennifer Donnelly, A Northern Light`,
	})

	var sock mangos.Socket

	sock, err := rep.NewSocket()
	if err != nil {
		log.Fatal("new socket:", err)
	}
	defer func() {
		err := sock.Close()
		if err != nil {
			log.Fatal("sock close:", err)
		}
	}()

	powMiddleware := server.NewPOWMiddleware(sock, hashCash)
	handler := server.NewHandler(sock, quotesBook)

	if err := sock.Listen("tcp://0.0.0.0:5555"); err != nil {
		log.Fatal("listen:", err)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
		syscall.SIGTERM,
	)
	defer cancel()

	go func() {
		for {
			err := powMiddleware.Handle(ctx)
			if err != nil {
				log.Println("pow middleware handle:", err)
				continue
			}

			err = handler.Handle(ctx)
			if err != nil {
				log.Println("handler handle:", err)
				continue
			}
		}
	}()

	<-ctx.Done()
}
