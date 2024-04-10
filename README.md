# faraway-test

[![codecov](https://codecov.io/gh/itimky/faraway-test/graph/badge.svg?token=EYGCJEHSDN)](https://codecov.io/gh/itimky/faraway-test)
[![Go Report Card](https://goreportcard.com/badge/github.com/itimky/faraway-test)](https://goreportcard.com/report/github.com/itimky/faraway-test)


## Test task for Server Engineer

```text
Design and implement “Word of Wisdom” tcp server.
 • TCP server should be protected from DDOS attacks with the
   Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work),
   the challenge-response protocol should be used.
 • The choice of the POW algorithm should be explained.
 • After Proof Of Work verification, server should send one of
   the quotes from “word of wisdom” book or any other collection
   of the quotes.
 • Docker file should be provided both for the server and for
   the client that solves the POW challenge
```

Chosen POW algorithm: [HashCash](https://en.wikipedia.org/wiki/Hashcash) as it is simple enough for a test.

<!-- https://mermaid.js.org/syntax/classDiagram.html -->

```mermaid
---
title: Code design
---
classDiagram
    direction LR
    namespace server {
        class ServerHandler["handler"] {
            Handle() error
        }

        class ServerPowMiddleware["middleware"] {
            Handle() error
        }

        class ServerBook["book"] {
            <<interface>>
            GetRandomQuote() string, error
        }

        class ServerPOW["POW"] {
            <<interface>>
            GenerateChallenge() string, error
            ValidateSolution(challenge string, solution int) error
        }
    }
    namespace client {
        class ClientTCP["Client"] {
            GetQuote() error
        }

        class ClientPOW["POW"] {
            <<interface>>
            SolveChallenge(challenge string, difficulty int) int
        }
    }
    namespace pow {
        class HashCash {
            GenerateChallenge() string, error
            SolveChallenge(challenge string, difficulty int) int
            ValidateSolution(challenge string, solution int) error
        }
    }
    namespace book {
        class Book {
            GetRandomQuote() string, error
        }
    }

    ServerHandler --> ServerBook
    ServerPowMiddleware --> ServerPOW
    ServerBook <|.. Book
    ServerPOW <|.. HashCash
    ClientTCP --> ClientPOW
    ClientPOW <|.. HashCash
```

```mermaid
---
title: Event design
---
sequenceDiagram
    participant Client
    participant Server

    Client ->> Server: ReqChallenge
    Server ->> Client: RepChallenge
    Client ->> Server: ReqQuote
    Server ->> Client: RepQuote
```
