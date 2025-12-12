# `ZeroNetworkLogger`: A Resilient Asynchronous Network Logger

This document outlines the design, implementation, and operational scenarios for the `ZeroNetworkLogger`, a robust logging solution designed for high availability and performance.

## 1. Design Philosophy

The primary goal of the `ZeroNetworkLogger` is to **decouple application performance from the availability of a remote logging server**. Standard network logging is often a blocking operation; if the remote server is slow, unresponsive, or down, the application can hang or crash.

This implementation avoids that problem by adopting a **non-blocking, asynchronous, buffered** approach. It ensures that the application can continue to log messages at full speed, regardless of the network connection's state.

## 2. Core Components

The logger is split into two main components:

1.  **`ZeroNetworkLogger`**: The public-facing struct that implements the standard `Logger` interface. Its primary role is to provide the familiar logging methods (`Info`, `Error`, etc.) and to configure the underlying writer.

2.  **`AsyncWriter`**: An `io.Writer` implementation that contains the core logic. It is responsible for:
    *   **Buffering**: Storing log messages in an in-memory channel (`logBuffer`) when they arrive.
    *   **Asynchronous Processing**: Running a dedicated background goroutine (`run()`) to manage the network connection and send logs from the buffer.
    *   **Automatic Reconnection**: Detecting connection failures and continuously attempting to re-establish a connection without blocking the application.
    *   **Fallback Logging**: Immediately writing logs to a console logger to ensure visibility even when the network is down.

## 3. How It Works: Scenarios

The system is designed to handle several key scenarios gracefully.

### Scenario 1: Normal Operation (Server is Up)

1.  When a log is generated (e.g., `logger.Info(...)`), the call is passed to `zerolog`, which invokes the `AsyncWriter.Write()` method.
2.  `Write()` immediately places a copy of the log message into the `logBuffer` channel. This is a very fast, non-blocking operation.
3.  The background `run()` goroutine, which is already connected, is continuously reading from `logBuffer`. It picks up the new message and writes it to the active network connection.

**Outcome**: Logs are sent to the remote server almost instantly with minimal impact on the application's performance.

### Scenario 2: Initial Startup (Server is Down)

1.  `NewZeroNetworkLogger` is called, which creates an `AsyncWriter` and starts the `run()` goroutine.
2.  The `run()` goroutine immediately detects that the connection (`w.conn`) is `nil` and calls `connectWithRetry()`.
3.  `connectWithRetry()` attempts to connect using `net.DialTimeout`. It fails but, instead of crashing, it simply waits for a `defaultRetryDelay` before looping to try again. This retry loop continues in the background.
4.  Meanwhile, the application starts and begins logging. Calls to `Write()` succeed instantly:
    *   The log is printed to the console via the `fallbackWriter`.
    *   The log is queued in the `logBuffer`.
5.  Eventually, the remote log server comes online. The `connectWithRetry()` call succeeds, establishing a connection.
6.  The `run()` goroutine moves on to `processBuffer()`, where it finds the queued logs and begins flushing them to the now-active connection.

**Outcome**: The application starts up without delay, logs are visible locally, and once the connection is established, all buffered logs are sent automatically.

### Scenario 3: Server Goes Down During Operation

1.  The system is running normally, and the `run()` goroutine is processing logs from the buffer.
2.  The remote server crashes. The next `w.conn.Write()` call returns a network error (e.g., "broken pipe").
3.  This error is detected inside `processBuffer()`:
    *   A failure message is printed to the console.
    *   The failed log message is safely put back at the front of the queue.
    *   The faulty connection is closed, and `w.conn` is set to `nil`.
    *   `processBuffer()` returns `false`, signaling a connection failure.
4.  The main `run()` loop sees the `false` return value, and, because `w.conn` is now `nil`, it loops back to the connection phase, calling `connectWithRetry()`.
5.  The system is now effectively back in **Scenario 2**, buffering new logs and attempting to reconnect in the background.

**Outcome**: The application continues to run unaffected. Logs are preserved in the buffer and visible on the console until the remote server is available again.

### Scenario 4: Graceful Shutdown

1.  The application calls the main `logger.Close()` method.
2.  This triggers `writer.Close()`, which closes the `stopChan`.
3.  The `run()` goroutine, which is listening on this channel, breaks out of its loop.
4.  Before exiting, it calls `drainBuffer()` to send any remaining logs from the `logBuffer` over the active connection, ensuring no data is lost.
5.  The goroutine exits, the `wg.Wait()` call unblocks, and the network connection is cleanly closed.

**Outcome**: The application shuts down cleanly, having flushed all pending log messages.

## 4. Note on Timestamp Accuracy

A common concern with asynchronous logging is whether timestamps reflect the time of the event or the time the log was sent.

This implementation **ensures timestamps are always accurate**.

The `zerolog` library captures the current time and serializes the entire log message (including the timestamp) into a JSON byte slice *at the moment the logging method is called*. This byte slice is then treated as an immutable package that gets placed in the buffer.

Even if the network connection is down for minutes, the worker goroutine will eventually send the original, unchanged byte slice. The remote log server will therefore always receive the timestamp from when the event actually occurred.

## 5. Architecture Diagram

Below is a simplified illustration of how `ZeroNetworkLogger` handles logs asynchronously:

```
+------------------+          +--------------------+
|  Application     |          | AsyncWriter        |
|  Goroutines      |          | (Buffered Writer)  |
|                  |          |                    |
| logger.Info(...) |  ----->  | logBuffer channel  |
| logger.Error(...)|          |  (FIFO queue)      |
+------------------+          +--------------------+
                                       |
                                       | background goroutine
                                       v
                              +--------------------+
                              | Network Connection |
                              |  (TCP or VSOCK)    |
                              +--------------------+
                                       |
                                       | writes log messages
                                       v
                             +---------------------+
                             | Remote Log Server   |
                             +---------------------+

Fallback Path:
If the network is unavailable or buffer is full:
+------------------+
| fallbackWriter    |
|  (console, file)  |
+------------------+
```

**Flow Explanation:**

1. Application calls `logger.X(...)`.
2. `AsyncWriter.Write()` copies the log message into `logBuffer`.
3. Background goroutine sends messages from `logBuffer` to the remote server.
4. If the network is down or buffer overflows, messages are written immediately to the `fallbackWriter`.

**Miscellanea**

Because of asynchronous writing and buffering/flushing, logs may not always be written to an output file in chronological order. Here is a useful command to sort a log file by timestamp:


```
$ jq -s 'sort_by(.time)[]'   bin/log_server_test.log
```