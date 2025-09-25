package communication

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"time"
)

//---
//---

// MessageReaderLoop continuously reads messages from a connection, parses them, and routes them.
// It is a shared implementation for both client and server.
func MessageReaderLoop(
	ctx context.Context,
	logPrefix string,
	conn net.Conn,
	reader *bufio.Reader,
	shutdownChan chan struct{},
	routeMessage func(context.Context, *Message),
	closeConnection func(),
) {
	defer closeConnection()

	// A constant for the read timeout, long enough to handle a brief pause in data flow
	const readTimeout = 1 * time.Second

	// A constant start/end buffer size for logging log messages
	const startEndBufSize = 16

	for {
		select {
		case <-shutdownChan:
			return
		case <-ctx.Done():
			return
		default:
			// Continue reading
		}

		// Use a temporary buffer to accumulate parts of a potentially large message
		var msgBytes []byte

		for {
			// Set the read deadline for each read operation
			conn.SetReadDeadline(time.Now().Add(readTimeout))

			// Read a chunk of data from the connection
			chunk, err := reader.ReadBytes(delimiter)
			if err != nil {
				// Check if it's a timeout error
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					if len(chunk) > 0 {
						log.Printf("%s: Read operation timed out while reading!", logPrefix)
						// bytes has been consumed on the connection communication buffer, save them before resuming reading
						msgBytes = append(msgBytes, chunk...)
						log.Printf("%s: added %d bytes to msgBytes buffer, total length now is %d.",
							logPrefix, len(chunk), len(msgBytes))
					}
					// Continue the inner loop, hoping more data arrives later in case of partial read
					continue
				}

				// If it's a non-timeout error, it's a fatal problem
				if err == io.EOF {
					// TODO shall we consider it an error?
					log.Printf("%s: Connection closed gracefully.", logPrefix)
				} else {
					log.Printf("%s: Read error, closing connection: %v", logPrefix, err)
				}
				return // Exit the entire function
			}

			// Append the received chunk to the message buffer
			if len(msgBytes) > 0 {
				// we have a partial read recovered here, log it
				log.Printf("%s: adding %d bytes to msgBytes buffer, total length %d",
					logPrefix, len(chunk), len(msgBytes)+len(chunk))
			}
			msgBytes = append(msgBytes, chunk...)

			// Check if the delimiter was found at the end of the chunk
			if len(chunk) > 0 && chunk[len(chunk)-1] == delimiter {
				break // The full message has been received, exit the inner loop
			} else {
				// TODO shall we consider it an error?
				log.Printf("%s: delimiter not found in msgBytes buffer, total length %d", logPrefix, len(msgBytes))
			}
		}

		// At this point, a full message has been read into msgBytes
		if len(msgBytes) > 2*startEndBufSize {
			log.Printf("%s: message %d bytes: %x...%x", logPrefix, len(msgBytes), msgBytes[:16], msgBytes[len(msgBytes)-16:])
		} else {
			log.Printf("%s: message %d bytes: %x", logPrefix, len(msgBytes), msgBytes)
		}

		// Parse and route the complete message
		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("%s: Error parsing message: %v", logPrefix, err)
			continue
		}

		log.Printf("%s: Received message: ID=%s, Type=%v", logPrefix, msg.ID, msg.Type)
		// Route message in a separate goroutine
		go routeMessage(ctx, &msg)
	}
}
