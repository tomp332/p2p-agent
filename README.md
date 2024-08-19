# P2P Agent

![Go Version](https://img.shields.io/badge/go-1.22-blue)
![GitHub Actions](https://github.com/tomp332/p2p-agent/actions/workflows/ci.yaml/badge.svg)
![Codecov Master](https://codecov.io/gh/tomp332/p2p-agent/branch/master/graph/badge.svg)

## Overview

P2P Agent is a peer-to-peer (P2P) network agent built in Go. This project provides a robust framework for creating decentralized, distributed systems. It leverages the power of Go to deliver high performance and scalability.

## Features

- **Decentralized Networking**: Fully decentralized peer-to-peer communication.
- **Scalable Architecture**: Supports scaling across multiple nodes.
- **Secure Communication**: Implements secure channels for data transmission.
- **Extensible**: Easily extendable for custom P2P applications.

## Getting Started

To get started with P2P Agent, follow these steps:

### Prerequisites

- Go 1.22 or later
- Git

### Installation

Clone the repository:

```bash
git clone https://github.com/tomp332/p2p-agent.git
cd p2p-agent
```

### Configuration Example

```yaml
server:
  host: localhost
  port: 8080
nodes:
  file-node:
      auth:
        username: admin
        password: admin
      bootstrap_peer_addrs:
        - host: localhost
          port: 9999
          username: admin
          password: 1
        - host: localhost
          port: 9991
          username: admin
          password: 1
logger_mode: dev
logger_level: trace

```