# Singularity 

## Overview

Singularity is a message broker that enables asynchronous communication between applications using a publish/subscribe pattern. Built with Go, it offers high performance, low latency, and a simple yet powerful protocol for message routing.

## Features

- 🚀 **High Performance**: Efficient message routing and delivery
- 🔄 **Pub/Sub Pattern**: Support for multiple publishers and subscribers
- 📦 **Multiple Queues**: Create and manage multiple message queues
- 🔌 **Custom Protocol**: Lightweight binary protocol for efficient communication

## Protocol Specification

Singularity uses a custom binary protocol designed for efficiency:

```
Packet Structure:
      1 Byte        2 Bytes          2 Bytes
  +---------------+----------------+----------------+
  |packet type    |remainingLenght |metadata length |
  +---------------|----------------|----------------|
  | Metadata      | message length |  message       |
  +---------------+----------------+----------------+
  0->(2^16)-3 Byte     2 Bytes      0->(2^16)-3 Byte
```

## Quick Start

### Installation

```bash
go install github.com/Oussama/singularity
```


