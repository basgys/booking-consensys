# Booking Consensys

## How To

### Dependencies

Install go 1.6

### Test

```shell
make dev
```

There will be a JWT displayed on stdout. It is a valid session for the seed user.
It must be used for endpoints, which require an authentication.

Open [API.md](API.md) to see how to use the REST API.

### Project structure

- **main.go** - It all starts with a main function
- **app** − It contains the domain logic (DDD style)
- **config** - It contains app configuration code
- **pkg** - It contains libraries that could be shared with other projects

## Source code

Besides some inspirations here and there, I've implemented every packages you see
on this project. Also, I created all repositories on github.com/deixis.

## Assumptions

Due to the time constraint, I've implemented a solution with several assumptions
instead of asking Nako.

- A room can be booked at any time without any restrictions
- A room can be booked by only one user at a time
- A room can be booked by the same user for several hours
- A room can only be booked every hour, on the hour (e.g. 6:00am to 7:00am)
- A room availability can be queried by anybody
- A room can only be reserved by an authenticated user
- A room can be cancelled by an authenticated user (not just the owner)

## Features

### Requested

- ✅ Backend in Go
- ✅ Users can see meeting rooms availability
- ✅ Users can book meeting rooms by the hour (first come first served)
- ✅ Users can cancel their own reservations

## Extra features

- ✅ Users can list meeting rooms
- ✅ Users can authenticate with JWT
- 🛑 Users can request a challenge a sign it with Metamask to authenticate (not finished)

## Possible improvements

- Validations
- More tests
- Frontend (Typescript/React)
- Booking restrictions/quotas (user can book only one room at a time)

## Key decisions

### No external database

To simplify testing, I have decided to use an embedded database in Go (Badger),
so you can easily test the booking system. This KV storage could be plugged
to an external database for production, such as FoundationDB.

I thought that a "raw" implementation would give more room for discussions than
using PostgreSQL for example. And Consensys works mainly on decentralised systems after all :)

Also, I have a strong interest in databases, so that is why I have fun
implementing storage layers in Go as a hobby. There is a KV and Event Sourcing implementation.
See [deixis/storage](https://github.com/deixis/storage)

(centralised/decentralised, SQL/NoSQL, Blockchain, Conflict resolution, Privacy, ...)

### REST API

To allow you to easily test the system, I have implemented a REST API. But a GraphQL/gRPC API could be
used as well. No strong motivations here besides making things easy.

### Authentication

I thought it would be great to test the authentication with a tiny frontend and Metamask. This would make
the "onboarding" easy and Blockchain makes a lot of sense here in my opinion. However, I haven't had
time to finish this implementation, but I have left it as is so we could pair and finalise it.

### Reservation

Reservations are stored on a KV storage, which makes the implementation ineficient. Every room
has one key and all reservations are serialised on that key. To make this implementation more
efficient, we could partition it in a different way.

I started an implementation a few years ago to partition intervals on a KV storage.
It is based on a whitepaper if you are interested.

[Overlap Interval Partition Join Whitepaper](https://files.ifi.uzh.ch/boehlen/Papers/DBG14.pdf)

## Metamask

### Sign challenge

```js
const accounts = await ethereum.request({ method: 'eth_requestAccounts' })
let challenge = "foo bar" // Received from backend
res = await ethereum.request({ method: 'personal_sign', params: [challenge, accounts[0]] })
```
