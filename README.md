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
title: Server Draft
---
classDiagram
    direction LR
    namespace server {
        class ServerHandler["handler"] {
            Handle(conn net.Conn) error
        }

        class ServerPowMiddleware["middleware"] {
            Handle(conn net.Conn) error
        }

        class ServerBook["book"] {
            <<interface>>
            GetRandomQuote(ctx) string, error
        }
        
        class ServerPOW["pow"] {
            <<interface>>
            GenerateChallenge() string, error
            SolveChallenge(challenge string, difficulty int) int
            ValidateSolution(challenge string, solution int) bool
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
            GetRandomQuote(ctx) string, error
        }
    }

    ServerHandler --> ServerBook
    ServerPowMiddleware --> ServerPOW
    ServerBook <|.. Book
    ServerPOW <|.. HashCash
```
