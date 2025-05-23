package dap

import (
	"bufio"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

// type DAPMessage struct {
// 	Seq     int             `json:"seq"`
// 	Type    string          `json:"type"`
// 	Command string          `json:"command,omitempty"`
// 	Args    json.RawMessage `json:"arguments,omitempty"`
// }

// type DAPResponse struct {
// 	Seq        int         `json:"seq"`
// 	Type       string      `json:"type"`
// 	RequestSeq int         `json:"request_seq"`
// 	Command    string      `json:"command"`
// 	Success    bool        `json:"success"`
// 	Body       interface{} `json:"body,omitempty"`
// 	Message    string      `json:"message,omitempty"`
// }

func Serve(r io.Reader, w io.Writer) error {
	reader := bufio.NewReader(r)

	messageIdx := 0
	storage := make(map[string]interface{})

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if !strings.HasPrefix(line, "Content-Length:") {
			continue
		}

		lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return err
		}

		reader.ReadString('\n')

		payload := make([]byte, length)
		io.ReadFull(reader, payload)

		var msg DAPMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			return err
		}

		appendToFile(string(payload), "clear")
		handleDAPMessage(msg, w, &messageIdx, storage)
	}
}
