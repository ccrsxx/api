# Rest API in Go

My first Rest API using Go. I'll use this Repo to learn and experiment with Go best practices for building RESTful APIs.

## Directory Structure

```md
├── cmd/
│ └── api/
│ └── main.go # Entry point (Equivalent to src/index.ts)
│
├── internal/ # Your application logic (Private code)
│ ├── config/ # (src/core/config)
│ │ └── env.go # Loads ENV variables
│ │
│ ├── middleware/ # (src/core/middlewares + loaders)
│ │ ├── cors.go # CORS middleware
│ │ ├── logger.go # Request logging
│ │ ├── rate_limit.go # Rate limiting logic
│ │ └── recovery.go # Global Error Handler (prevents crashes)
│ │
│ ├── modules/ # (src/modules) - Feature based
│ │ ├── auth/
│ │ │ ├── handler.go # (auth.controller.ts)
│ │ │ ├── service.go # (auth.service.ts)
│ │ │ └── route.go # (auth.route.ts)
│ │ ├── home/
│ │ ├── jellyfin/
│ │ ├── og/ # The Proxy handler we discussed
│ │ ├── spotify/
│ │ ├── sse/
│ │ └── tools/
│ │
│ ├── server/ # Wiring everything together
│ │ └── routes.go # Registers all module routes to the Mux
│ │
│ └── utils/ # (src/core/utils)
│ ├── response.go # JSON response helpers (success/error)
│ ├── validator.go # Validation logic (Zod replacement)
│ └── http.go # HTTP client helpers
│
├── go.mod # Package manager (package.json)
└── go.sum # Lock file (package-lock.json)
```
