package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Mirror the admin message types from pkg/admin/admin_interface.go.
// We duplicate them here so the CLI binary has zero internal dependencies and
// can be distributed as a standalone tool. Tests import pkg/admin to guard
// against drift.
const (
	adminResponseMessage  = "response"
	adminErrorMessage     = "error"
	keyAttestationRequest = "key_attestation"
	getVersionRequest     = "get_version"
	setLogLevelRequest    = "set_log_level"
	getLogLevelRequest    = "get_log_level"
)

var validLogLevels = []string{
	"trace", "debug", "info", "warn", "error", "fatal", "panic", "disabled",
}

var validTargets = []string{"manager", "executor", "all"}

type adminMessage struct {
	Type   string          `json:"type"`
	Target string          `json:"target,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

type errorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type setLogLevelReq struct {
	Level string `json:"level"`
}

var useColors bool

func main() {
	host := flag.String("host", "localhost", "Admin server host")
	port := flag.Int("port", 4002, "Admin server port")
	flag.BoolVar(&useColors, "colors", true, "Colorize output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Interactive CLI for sending admin commands to the Vela Manager.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -host=<addr>    Admin server host (default \"localhost\")\n")
		fmt.Fprintf(os.Stderr, "  -port=<number>  Admin server port (default 4002)\n")
		fmt.Fprintf(os.Stderr, "  -colors=<bool>  Colorize output (default true)\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -host=10.0.0.5 -port=5000\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -colors=false\n", os.Args[0])
	}

	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Vela Admin CLI")
	fmt.Printf("Server: %s\n\n", addr)

	for {
		printMenu()
		choice := prompt(scanner, "Select command")
		if choice == "" {
			continue
		}

		var msg *adminMessage
		var err error

		switch choice {
		case "1":
			msg, err = buildGetVersion(scanner)
		case "2":
			msg, err = buildGetLogLevel(scanner)
		case "3":
			msg, err = buildSetLogLevel(scanner)
		case "4":
			msg = buildKeyAttestation()
		case "5", "q", "Q":
			fmt.Println("Bye.")
			return
		default:
			fmt.Println("Invalid choice, try again.")
			continue
		}

		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			continue
		}

		sendAndPrint(addr, msg)
		fmt.Println()
	}
}

func printMenu() {
	fmt.Println("Commands:")
	fmt.Println("  1) Get Version")
	fmt.Println("  2) Get Log Level")
	fmt.Println("  3) Set Log Level")
	fmt.Println("  4) Key Attestation")
	fmt.Println("  5) Quit")
}

func prompt(scanner *bufio.Scanner, label string) string {
	fmt.Printf("%s> ", label)
	if !scanner.Scan() {
		fmt.Println()
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0) // EOF
	}
	return strings.TrimSpace(scanner.Text())
}

func promptChoice(scanner *bufio.Scanner, label string, options []string) (string, error) {
	for i, opt := range options {
		fmt.Printf("  %d) %s\n", i+1, opt)
	}
	raw := prompt(scanner, label)
	// Accept either the number or the literal value.
	for i, opt := range options {
		if raw == strconv.Itoa(i+1) || strings.EqualFold(raw, opt) {
			return opt, nil
		}
	}
	return "", fmt.Errorf("invalid choice: %q", raw)
}

// --- builders ---

func buildGetVersion(scanner *bufio.Scanner) (*adminMessage, error) {
	fmt.Println("Target:")
	target, err := promptChoice(scanner, "Target", validTargets)
	if err != nil {
		return nil, err
	}
	return &adminMessage{Type: getVersionRequest, Target: target}, nil
}

func buildGetLogLevel(scanner *bufio.Scanner) (*adminMessage, error) {
	fmt.Println("Target:")
	target, err := promptChoice(scanner, "Target", validTargets)
	if err != nil {
		return nil, err
	}
	return &adminMessage{Type: getLogLevelRequest, Target: target}, nil
}

func buildSetLogLevel(scanner *bufio.Scanner) (*adminMessage, error) {
	fmt.Println("Log level:")
	level, err := promptChoice(scanner, "Level", validLogLevels)
	if err != nil {
		return nil, err
	}
	fmt.Println("Target:")
	target, err := promptChoice(scanner, "Target", validTargets)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(setLogLevelReq{Level: level})
	if err != nil {
		return nil, fmt.Errorf("marshal log level: %w", err)
	}
	raw := json.RawMessage(data)
	return &adminMessage{Type: setLogLevelRequest, Target: target, Data: raw}, nil
}

func buildKeyAttestation() *adminMessage {
	return &adminMessage{Type: keyAttestationRequest, Target: "executor"}
}

// --- send / receive ---

func sendAndPrint(addr string, msg *adminMessage) {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		fmt.Printf("\n%sConnection error: %v%s\n", color(ansiRed), err, color(ansiReset))
		return
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(30 * time.Second)); err != nil {
		fmt.Printf("\n%sDeadline error: %v%s\n", color(ansiRed), err, color(ansiReset))
		return
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("\n%sMarshal error: %v%s\n", color(ansiRed), err, color(ansiReset))
		return
	}
	// Server expects JSON + newline delimiter.
	payload = append(payload, '\n')

	if _, err := conn.Write(payload); err != nil {
		fmt.Printf("\n%sSend error: %v%s\n", color(ansiRed), err, color(ansiReset))
		return
	}

	reader := bufio.NewReader(conn)
	respBytes, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Printf("\n%sRead error: %v%s\n", color(ansiRed), err, color(ansiReset))
		return
	}

	var resp adminMessage
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		fmt.Printf("\n%sParse error: %v%s\n", color(ansiRed), err, color(ansiReset))
		return
	}

	switch resp.Type {
	case adminResponseMessage:
		printSuccess(resp.Data)
	case adminErrorMessage:
		printError(resp.Data)
	default:
		fmt.Printf("\n%sUnexpected response type: %s\nRaw: %s%s\n", color(ansiRed), resp.Type, string(resp.Data), color(ansiReset))
	}
}

// color returns the ANSI escape sequence when colors are enabled, or an empty
// string otherwise.
func color(code string) string {
	if useColors {
		return code
	}
	return ""
}

const (
	ansiGreen = "\033[32m"
	ansiRed   = "\033[31m"
	ansiReset = "\033[0m"
)

func printSuccess(data json.RawMessage) {
	// Try to pretty-print the data; fall back to raw string.
	var pretty any
	if err := json.Unmarshal(data, &pretty); err == nil {
		formatted, _ := json.MarshalIndent(pretty, "", "  ")
		fmt.Printf("\n%sOK: %s%s\n", color(ansiGreen), string(formatted), color(ansiReset))
	} else {
		fmt.Printf("\n%sOK: %s%s\n", color(ansiGreen), string(data), color(ansiReset))
	}
}

func printError(data json.RawMessage) {
	var ed errorData
	if err := json.Unmarshal(data, &ed); err == nil {
		fmt.Printf("\n%sERROR [%s]: %s%s\n", color(ansiRed), ed.Code, ed.Message, color(ansiReset))
	} else {
		fmt.Printf("\n%sERROR: %s%s\n", color(ansiRed), string(data), color(ansiReset))
	}
}
