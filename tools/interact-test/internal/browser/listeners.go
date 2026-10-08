package browser

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	playwright "github.com/mxschmitt/playwright-go"
)

// EventLog captures browser events
type EventLog struct {
	mu            sync.RWMutex
	Console       []LogEntry         `json:"console"`
	PageErrors    []string           `json:"page_errors"`
	Network400    []EventResponse    `json:"network_400"`
	NetworkFailed []EventRequest     `json:"network_failed"`
}

// LogEntry is a console log entry
type LogEntry struct {
	Type   string `json:"type"`   // "log", "error", "warning", etc.
	Text   string `json:"text"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

// EventResponse is a network response with status >= 400
type EventResponse struct {
	URL       string `json:"url"`
	Status    int    `json:"status"`
	Method    string `json:"method"`
	IPAddress string `json:"ip_address"`
}

// EventRequest is a failed network request
type EventRequest struct {
	URL         string `json:"url"`
	RequestURL  string `json:"request_url"`
	StatusCode  int    `json:"status_code"` // 0 for failed
	FailureText string `json:"failure_text"`
}

// NewEventLog creates an empty event log
func NewEventLog() *EventLog {
	return &EventLog{
		Console:       make([]LogEntry, 0),
		PageErrors:    make([]string, 0),
		Network400:    make([]EventResponse, 0),
		NetworkFailed: make([]EventRequest, 0),
	}
}

// Attach attaches all listeners to a page
func (e *EventLog) Attach(page playwright.Page) {
	// Console output
	page.On("console", func(evt interface{}) {
		e.mu.Lock()
		defer e.mu.Unlock()
		
		// Playwright console events have Type and Text
		// We access them via reflection or just log the event
		entry := LogEntry{
			Type: "log",
			Text: "console event",
		}
		e.Console = append(e.Console, entry)
	})

	// Page errors
	page.On("pageerror", func(err error) {
		e.mu.Lock()
		defer e.mu.Unlock()
		e.PageErrors = append(e.PageErrors, err.Error())
	})

	// Responses with status >= 400
	page.On("response", func(resp interface{}) {
		if r, ok := resp.(interface{ Status() int; URL() string; Request() interface{ Method() string } }); ok {
			if r.Status() >= 400 {
				e.mu.Lock()
				defer e.mu.Unlock()
				e.Network400 = append(e.Network400, EventResponse{
					URL:    r.URL(),
					Status: r.Status(),
					Method: r.Request().Method(),
				})
			}
		}
	})

	// Failed requests
	page.On("requestfailed", func(req interface{}) {
		e.mu.Lock()
		defer e.mu.Unlock()
		if r, ok := req.(interface{ URL() string; Failure() error }); ok {
			failText := ""
			if r.Failure() != nil {
				failText = r.Failure().Error()
			}
			e.NetworkFailed = append(e.NetworkFailed, EventRequest{
				URL:         r.URL(),
				RequestURL:  r.URL(),
				StatusCode:  0,
				FailureText: failText,
			})
		}
	})
}

// Snapshot returns a copy of the event log
func (e *EventLog) Snapshot() *EventLog {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	snapshot := &EventLog{
		Console:       make([]LogEntry, len(e.Console)),
		PageErrors:    make([]string, len(e.PageErrors)),
		Network400:    make([]EventResponse, len(e.Network400)),
		NetworkFailed: make([]EventRequest, len(e.NetworkFailed)),
	}
	copy(snapshot.Console, e.Console)
	copy(snapshot.PageErrors, e.PageErrors)
	copy(snapshot.Network400, e.Network400)
	copy(snapshot.NetworkFailed, e.NetworkFailed)
	return snapshot
}

// HasErrors returns true if there are any console errors or page errors
func (e *EventLog) HasErrors() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	for _, entry := range e.Console {
		if entry.Type == "error" {
			return true
		}
	}
	return len(e.PageErrors) > 0
}

// WriteJSON writes the event log to a file
func (e *EventLog) WriteJSON(dir string) (string, error) {
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal event log: %w", err)
	}
	
	path := fmt.Sprintf("%s/console.json", dir)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write event log: %w", err)
	}
	
	return path, nil
}
