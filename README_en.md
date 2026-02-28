# chLogger 🚀 (Channel Logger)

**chLogger** is a high-performance, asynchronous logger for Go, built on channels and the worker-pool pattern.

Designed to minimize I/O latency in the main application thread, it is optimized for high-throughput data ingestion. Logs are sent to a non-blocking channel and processed by dedicated minion goroutines. It requires no heavy external dependencies or complex structs — everything operates on standard Go types.

## ✨ Key Features

* **Async Out-of-the-Box**: Logs are written to a buffered `ChInLog chan [4]string` channel, ensuring your business logic is never blocked.
* **Pub/Sub Architecture & Extensibility**: Easily broadcast logs to any custom subscriber (via the `Broadcast` channel) without altering the logger's core.
* **Real-time WebSocket Streaming**: Includes a built-in `wsLoggerPlugin` that spins up a web server with a UI to stream logs directly to your browser.
* **HTTP RPC Input**: The logger can act as a standalone service, accepting JSON messages via HTTP to log events from other microservices.
* **Flexible Routing & Filtering**: Filter outputs by function name and module. Automatically groups logs and errors into separate files.

## 🛠 Data Contract

Interacting with the logger is extremely simple. You don't need to import complex types; just send an array of 4 strings (`[4]string`) into the channel:

* `[0]` **Func** — Function name (or "nil"). Used for filtering (`ConsolFilterFn`).
* `[1]` **Unit** — Package/Object (module) name. Used for filtering (`ConsolFilterUn`).
* `[2]` **Text** — The actual log message / payload.
* `[3]` **Status** — Message status. An error is marked as `"1"` or `"ERROR"`, while normal execution is `""` or `"0"`.

## ⚠️ Important: Graceful Shutdown

Because logs are queued in a channel buffer, a **hard crash or panic** may result in lost logs. 

To guarantee that all messages are written to their respective destinations, you must use the graceful shutdown method (`Stop()`), which waits for the channels to drain before exiting.

## 📦 Quick Start

### Basic Usage

```go
package main

import (
	"fmt"
	"[github.com/xela07ax/chLogger](https://github.com/xela07ax/chLogger)"
	"time"
)

func main() {
	// 1. Configure the logger
	logEr := chLogger.NewChLoger(&chLogger.Config{
		ConsolFilterFn: map[string]int{"Front Http Server": 0},
		ConsolFilterUn: map[string]int{"Pooling": 1},
		Mode:           0,
		Dir:            "x-loger", // Directory for log files
	})
	
	// 2. Start the daemon (minions)
	logEr.RunLogerDaemon()

	// 3. Send logs: [Func, Unit, Text, Status]
	logEr.ChInLog <- [4]string{"Welcome", "nil", "System successfully started", "0"}
	logEr.ChInLog <- [4]string{"Database", "Connection", "Connection timeout", "1"} // 1 = Error

	// Simulate work
	time.Sleep(1 * time.Second)
	
	// 4. Graceful shutdown (wait for all logs to be written)
	logEr.Stop()
}
```

### 🔌 Extensibility: Adding Custom Subscribers
You can easily add new handlers to the logger (e.g., to send data to Kafka, Telegram, or ELK). Just pass your own channel to the Broadcast field in the configuration:

```go
// Create a channel for your custom subscriber
myCustomSubscriber := make(chan []byte, 100)

// Connect it to the logger config
cfg := &chLogger.Config{
    Dir:       "x-loger",
    Broadcast: myCustomSubscriber, 
}

// Start your worker/handler
go func(in <-chan []byte) {
    for msg := range in {
        // Any custom byte processing logic goes here
        fmt.Printf("Subscriber intercepted log: %s\n", string(msg))
    }
}(myCustomSubscriber)
```

### 🌐 Built-in WebSocket Plugin
As a ready-to-use example of a subscriber, the project includes wsLoggerPlugin. It intercepts the Broadcast channel and broadcasts logs to all connected clients via WebSockets.

```go
hub := wsLoggerPlugin.NewWsLogger()
go hub.Run()

http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    hub.ServeWs(w, r)
})
```
