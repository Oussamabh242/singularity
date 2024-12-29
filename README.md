# Singularity

Singularity is a high-performance message broker designed to enable asynchronous communication between applications using a publish/subscribe pattern. Built in Go, Singularity offers low latency, efficient message routing, and a simple yet powerful protocol.

## Features

- 🚀 **High Performance**: Efficient message routing and delivery
- 🔄 **Pub/Sub Pattern**: Supports multiple publishers and subscribers for flexible communication
- 📦 **Multiple Queues**: Allows creation and management of multiple message queues
- 🔌 **Custom Protocol**: Utilizes a lightweight binary protocol for faster communication

## Protocol Specification

Singularity uses a custom binary protocol designed for efficient message handling. Here"s the structure of each packet:

```
  1 Byte        2 Bytes          2 Bytes
+---------------+----------------+----------------+
| packet type   | remainingLength | metadata length |
+---------------+----------------+----------------+
| Metadata      | message length |  message       |
+---------------+----------------+----------------+
```

- **Packet Type (1 byte)**: Identifies the type of packet.
- **Remaining Length (2 bytes)**: Specifies the length of the following data.
- **Metadata Length (2 bytes)**: Length of the metadata section.
- **Metadata**: Additional details about the message.
- **Message Length (2 bytes)**: Length of the actual message content.
- **Message**: The main content sent between applications.

## Quick Start

### Installation

To install Singularity, use the following command:

```bash
go install github.com/Oussama/singularity
```

### CLI Usage

Singularity comes with a command-line interface for easy server management. The main command is `serve`, which starts the message broker server.

#### `serve` Command

This command starts a TCP server for Singularity on the specified port.

**Usage:**

```bash
singularity serve [flags]
```

**Flags:**

| Flag           | Shortcut | Default | Description                                                        |
|----------------|----------|---------|--------------------------------------------------------------------|
| `--port`        | `-p`     | 1234    | Port to listen on                                                  |
| `--subscribers` | `-s`     | 10      | Maximum number of subscribers allowed per queue                    |
| `--messages`    | `-m`     | 100     | Maximum number of messages per queue before blocking further messages |

**Example:**

```bash
singularity serve -p 8080 -s 50 -m 500
```

This starts the server on port `8080`, allowing up to `50` subscribers per queue and a maximum of `500` messages per queue.

## JavaScript Client Library

You can interact with the Singularity server through the JavaScript client library \`@9k9k/singularity-js\`. 

### Client Library Installation

To install the JavaScript client library \`@9k9k/singularity-js\`, run the following command in your project:

\`\`\`bash
npm install @9k9k/singularity-js
\`\`\`

### Example of how to use the library:

\`\`\`javascript
import { Singularity } from "@9k9k/singularity-js";

const sing = new Singularity("localhost", 1234);

// Create a queue with name "something".
sing.connect2(sing.createQueue, "something");

// publish 10 messages to the queue "something".
for (let i = 0; i < 10; i++) {
  sing.publish("hello there", {
    queue: "something"
  });
}

function job(msg) {
  console.log(msg);
}
// subscribe to the queue "something" and treat incomming messages with the "job" function .
sing.subscribe("something", job);
\`\`\`

- The code provides the ability to connect, create queues, publish messages, and subscribe to queues using the publish/subscribe model.

### Important Notes
- The library is **still under development** and may not be stable.
