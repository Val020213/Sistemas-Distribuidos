# Distributed Web Scraper — Chord DHT

A horizontally scalable web-scraping service built on a **Chord distributed hash table**. Users submit URLs through a web UI; the ring decides which node owns each URL, scrapes it, replicates the result, and serves it back — and it keeps working while nodes join, leave, or crash.

Final project for **Distributed Systems**, Computer Science, University of Havana (2024–2025).

> 🇪🇸 Spanish version: [`README.es.md`](./README.es.md) · Course reports (Spanish): [`2doInforme.md`](./2doInforme.md), [`Seguridad.md`](./Seguridad.md)

---

## What it does

- **Submit a URL** for scraping from the web UI. The request is hashed onto the ring and routed to the node that owns that key.
- **Track status** — every URL is `pending`, `success`, or `error`. Failed jobs can be retried.
- **Download** the scraped content of any successfully processed URL, including URLs requested by *other* users — the ring is a shared store, not a per-user cache.
- **Scale by adding containers.** New server nodes discover the ring on their own and take over their share of the keyspace; no configuration file lists the members.

---

## Architecture

```
┌──────────────────────────┐        ┌────────────────────────────────────────┐
│  scrapper-client-network │        │        scrapper-server-network         │
│      10.0.11.0/24        │        │             10.0.10.0/24               │
│                          │        │                                        │
│  frontend1 (Next.js)     │        │   backend1 ──┐                         │
│  frontend2 (Next.js)     │        │   backend2   ├── Chord ring (gRPC)     │
│      ...                 │        │   backend3 ──┘   each + its own MongoDB │
└───────────┬──────────────┘        └───────────────────┬────────────────────┘
            │  HTTPS / REST                             │  gRPC over mTLS
            └──────────────► router (Go) ◄──────────────┘
                    10.0.11.254 / 10.0.10.254
```

The system runs on **two isolated Docker networks** — one for clients, one for servers — bridged only by a **router** container that is dual-homed on both. Clients never address ring nodes directly.

### The ring

Each node's identifier is `SHA-1(address)` truncated to `CHORD_BITS` bits (default `8`, so a 256-slot identifier space). Every node maintains:

| Structure | Purpose |
| --- | --- |
| `Successors[]` | Successor list — survives the immediate successor dying |
| `Predecessor` | Filled in by `Notify`, verified by liveness checks |
| `Finger[m]` | Finger table, `m = CHORD_BITS` entries, for O(log n) lookups |
| `SuccessorCache[]` | Last known-good successors, used when the ring is in flux |

Background goroutines run the standard Chord maintenance loop — `Stabilize`, `FixFingersTable`, `CheckPredecessor` — plus a `replicateData` pass. `FindSuccessor` carries a hop counter and a `visited` map in the request so a partially-healed ring can't send a lookup around in circles.

### Discovery (no bootstrap list)

Server nodes announce themselves every 5 s over **UDP multicast** to `224.0.0.1:9999`. A joining node listens, picks up a live peer, and calls `FindSuccessor` against it to locate its place in the ring. The router listens on the same group to maintain its own table of reachable backends, expiring any it hasn't heard from in 15 s. **Nothing in the system is configured with a static peer list.**

### Replication & fault tolerance

Each node replicates the keys in its range `(predecessor, self]` to its successors, holding **3 replicas** (`tolerance = 3`), and drops the stale copy from the node that falls out of the window. That gives **level-2 fault tolerance**: the system keeps serving with two nodes down. When a node rejoins, `GetNodeDataFromOldSuccessor` pulls back the range it now owns and the ring re-converges.

### Transport security

Every component holds a certificate issued by a **self-signed private CA** and every hop runs over TLS 1.2+. How strictly each end *verifies* differs, and the repo is worth reading honestly on this point:

| Hop | Encrypted | Peer verification |
| --- | --- | --- |
| gRPC inbound → ring node | ✅ | ✅ `RequireAndVerifyClientCert` — the node rejects any caller without a CA-signed certificate |
| gRPC outbound ← ring node | ✅ | ❌ `InsecureSkipVerify` — the caller presents its certificate but doesn't validate the peer's |
| Client ↔ router | ✅ | ❌ `NoClientCert` — the router doesn't request client certificates |
| Router → ring | ✅ | ❌ `InsecureSkipVerify` |

So the ring genuinely refuses unauthenticated members, but full mutual authentication is **designed and documented, not yet enforced end to end** — closing the gap means dropping `InsecureSkipVerify` and switching the router to `RequireAndVerifyClientCert`. The intended design and the certificate-generation procedure are in [`Seguridad.md`](./Seguridad.md) *(Spanish)*.

There is no user authentication either — the system stores no user data, and the threat model is a coursework deployment. JWT or OAuth at the router would be the natural next step.

---

## Tech stack

| Layer | Choice |
| --- | --- |
| Ring nodes | Go 1.23, gRPC + Protocol Buffers |
| HTTP API | Gin (`/tasks`, `/task`, `/health`) |
| Router / gateway | Go, `net/http` reverse proxy, TLS on both sides |
| Per-node storage | MongoDB (one instance per ring node) |
| Discovery | UDP multicast (`224.0.0.1:9999`) |
| Client | Next.js 15 (App Router), React 19, TypeScript, MUI + MUI X Data Grid, React Hook Form + Zod, Tailwind — served over HTTPS with its own client certificate |
| Packaging | Docker, two user-defined bridge networks |
| Tests | Go `testing` + Testcontainers (ephemeral MongoDB) |

---

## Repository layout

```
src/
├─ server/                     # Chord ring node
│  ├─ internal/chord/          # ring logic, chord.proto, generated gRPC stubs
│  ├─ internal/scraper/        # worker pool that fetches queued URLs
│  ├─ internal/database/       # MongoDB service layer (+ Testcontainers tests)
│  ├─ internal/multicast/      # periodic self-announcement
│  └─ internal/server/         # Gin routes
├─ router/                     # dual-homed gateway: discovery + mTLS proxy
├─ client/                     # Next.js frontend
└─ certs/                      # private CA + component certificates
docker-utils/{sh,bat}/         # network bootstrap and node-scaling scripts
```

---

## Running it

**Requirements:** Docker and Docker Compose.

```bash
git clone https://github.com/Val020213/Sistemas-Distribuidos
cd Sistemas-Distribuidos
```

**1. Create the networks and start the router**

```bash
chmod +x docker-utils/sh/*.sh
sudo bash docker-utils/sh/init-network.sh      # Linux
docker-utils\bat\init-networks.bat             # Windows
```

**2. Add ring nodes** — run once per node you want (each gets its own MongoDB):

```bash
sudo bash docker-utils/sh/add-server.sh
```

**3. Add clients** — run once per frontend instance:

```bash
sudo bash docker-utils/sh/add-client.sh
```

Clients are published on `3001`, `3002`, … — e.g. <http://localhost:3001>. From the landing page choose **Iniciar Scrapper**, add a URL, and download the result once the status flips to `success`.

**Killing nodes** — pass the instance numbers to remove. This is also how you exercise fault tolerance: drop up to two nodes and the ring keeps serving.

```bash
sudo bash docker-utils/sh/clear-servers.sh 2 3
```

### Configuration

| Variable | Default | Meaning |
| --- | --- | --- |
| `CHORD_BITS` | `8` | Width of the identifier space (`2^m` slots) |
| `IP_ADDRESS` | — | Address the node hashes into its ring ID |
| `RPC_PORT` | — | gRPC listen port |
| `NEXT_PUBLIC_API_URL` | `http://10.0.11.254:8080` | Router address, as seen by clients |

The replication factor is currently the compile-time constant `tolerance = 3` in `src/server/internal/chord/chord.go`.

### Tests

```bash
cd src/server && go test ./...
```

Database and route tests spin up a throwaway MongoDB via Testcontainers, so Docker must be running.

---

## Authors

- **Osvaldo R. Moreno Prieto** — [@Val020213](https://github.com/Val020213)
- **Daniel Toledo Martínez** — [@Phann020126](https://github.com/Phann020126)
