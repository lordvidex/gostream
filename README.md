# gostream

## Working Principle
<p align="center">
  <img src="./assets/sample.png"/>
</p>

## Motivation
How can we create several different clients that connect to only one server in a cluster but receive updates performed on
any server in the cluster?

## Architecture
gostream has two main components:
- gRPC servers that communicate with one another using redis channels to sync data 
- gRPC clients that connect to ANY of the servers and receive combined updates from all servers.

### Other components
- CLI / TUI interface for client interactions
- CLI for retrieving server status and extra information

## Creating servers and clients
### Install CLI tool

```bash
go install github.com/lordvidex/gostream/cmd/gostream@latest
```

### 
```bash
# Create a new server instance
gostream server serve --port <port> -c config.toml

# Create a new client instance with list of server instances that are client-side loadbalanced
gostream client -c config.toml
```

## Things to check out for?
- How long-lived are the connections? Client reconnections when TCP connection gets broken? Can I connect to the same stream for a month+??
- Server stream connection cleanup during disconnect.
- Proper monitoring.
