package communication

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/horizen-pes/pkg/logger"
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
	routeMessage func(context.Context, Message),
	closeConnection func(),
	log logger.Logger,
) {
	defer closeConnection()
	log.Info("%s: Entering message reader loop!", logPrefix)

	for {
		select {
		case <-shutdownChan:
			return
		case <-ctx.Done():
			return
		default:
			// Continue reading
		}

		msgBytes, err := ReadMessageFromSocket(conn, reader, logPrefix, log)
		if err != nil {
			return // Exit the entire function
		}

		// Parse and route the complete message
		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Error("%s: Error parsing message: %v", logPrefix, err)
			continue
		}

		log.Info("%s: Received message: ID=%s, Type=%v", logPrefix, msg.ID, msg.Type)
		// Route message in a separate goroutine, making a copy of the message to avoid capturing a loop variable.
		// By creating this explicit copy, we ensure that each goroutine receives a pointer to its own message,
		// that is local to this specific iteration of the loop
		msgCopy := msg
		go routeMessage(ctx, msgCopy)
	}
}

func ReadMessageFromSocket(conn net.Conn, reader *bufio.Reader, logPrefix string, log logger.Logger) ([]byte, error) {
	// A constant for the read timeout, long enough to handle a brief pause in data flow
	const readTimeout = 1 * time.Second

	// A constant start/end buffer size for logging log messages
	const startEndBufSize = 16

	// Use a temporary buffer to accumulate parts of a potentially large message
	var msgBytes []byte

	for {
		// Set the read deadline for each read operation
		conn.SetReadDeadline(time.Now().Add(readTimeout))

		// Read a chunk of data from the connection
		chunk, err := reader.ReadBytes(MsgDelimiter)
		if err != nil {
			// Check if it's a timeout error
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				if len(chunk) > 0 {
					log.Warn("%s: Read operation timed out while reading!", logPrefix)
					// bytes has been consumed on the connection communication buffer, save them before resuming reading
					msgBytes = append(msgBytes, chunk...)
					log.Info("%s: added %d bytes to msgBytes buffer, total length now is %d.",
						logPrefix, len(chunk), len(msgBytes))
				}
				// Continue the inner loop, hoping more data arrives later in case of partial read
				continue
			}

			// If it's a non-timeout error, it's a fatal problem
			if err == io.EOF {
				// TODO shall we consider it an error?
				log.Warn("%s: Connection closed gracefully.", logPrefix)
			} else {
				log.Warn("%s: Read error, closing connection: %v", logPrefix, err)
			}
			return nil, fmt.Errorf("connection closed")
		}

		// Append the received chunk to the message buffer
		if len(msgBytes) > 0 {
			// we have a partial read recovered here, log it
			log.Warn("%s: adding %d bytes to msgBytes buffer, total length %d",
				logPrefix, len(chunk), len(msgBytes)+len(chunk))
		}
		msgBytes = append(msgBytes, chunk...)

		// Check if the delimiter was found at the end of the chunk
		if len(chunk) > 0 && chunk[len(chunk)-1] == MsgDelimiter {
			break // The full message has been received, exit the inner loop
		} else {
			// TODO shall we consider it an error?
			log.Warn("%s: delimiter not found in msgBytes buffer, total length %d", logPrefix, len(msgBytes))
		}
	}

	// At this point, a full message has been read into msgBytes
	if len(msgBytes) > 2*startEndBufSize {
		log.Debug("%s: message %d bytes: %x...%x", logPrefix, len(msgBytes), msgBytes[:16], msgBytes[len(msgBytes)-16:])
	} else {
		log.Debug("%s: message %d bytes: %x", logPrefix, len(msgBytes), msgBytes)
	}
	return msgBytes, nil
}
