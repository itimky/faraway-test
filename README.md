# faraway-test

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

<!-- https://mermaid.js.org/syntax/classDiagram.html -->

```mermaid
---
title: Code design
---
classDiagram
    direction LR
    namespace contract {
        class ReqRepChallenge {
        }
        class ReqRepQuote {
        }
    }
    namespace server {
        class ServerTCP {
            Start() error
        }

        class ServerHandler["handler"] {
            Handle() error
        }

        class ServerPowMiddleware["middleware"] {
            Handle() error
        }

        class ServerHandlerI["handler"] {
            <<interface>>
            Handle() error
        }

        class ServerBook["book"] {
            <<interface>>
            GetRandomQuote() string, error
        }

        class ServerPOW["POW"] {
            <<interface>>
            GenerateChallenge() string, error
            ValidateSolution(challenge string, solution int) bool
        }
    }
    namespace client {
        class ClientTCP["TCP"] {
            GetQuote() error
        }

        class ClientPOW["POW"] {
            SolveChallenge(challenge string, difficulty int) int
        }
    }
    namespace pow {
        class HashCash {
            GenerateChallenge() string, error
            SolveChallenge(challenge string, difficulty int) int
            ValidateSolution(challenge string, solution int) bool
        }
    }
    namespace book {
        class Book {
            GetRandomQuote() string, error
        }
    }

    ReqRepChallenge <-- ServerPowMiddleware
    ReqRepQuote <-- ServerHandler
    ServerTCP --> ServerHandlerI
    ServerHandlerI <|.. ServerHandler
    ServerHandlerI <|.. ServerPowMiddleware
    ServerHandler --> ServerBook
    ServerPowMiddleware --> ServerPOW
    ServerBook <|.. Book
    ServerPOW <|.. HashCash
    ReqRepChallenge <-- ClientTCP
    ReqRepQuote <-- ClientTCP
    ClientTCP --> ClientPOW
    ClientPOW <|.. HashCash
```

```mermaid
---
title: Event design
---
sequenceDiagram
    participant Client
    participant ServerProxy
    participant ServerWorker

    Client ->> ServerProxy: ReqChallenge
    ServerProxy ->> ServerWorker: ReqChallenge
    ServerWorker ->> ServerProxy: RepChallenge
    ServerProxy ->> Client: RepChallenge
    Client ->> ServerProxy: ReqQuote
    ServerProxy ->> ServerWorker: ReqQuote
    ServerWorker ->> ServerProxy: RepQuote
    ServerProxy ->> Client: RepQuote
```
