package uxgates

import (
	"encoding/json"
	"fmt"
	"os"
)

// Mode defines gate behavior mode
type Mode int

const (
	ModeOff Mode = iota // Gates don't run
	ModeWarn            // Gates run but don't fail
	ModeFail            // Gates run and failures cause run failure
)

// String converts mode to string
func (m Mode) String() string {
	switch m {
	case ModeOff:
		return "off"
	case ModeWarn:
		return "warn"
	case ModeFail:
		return "fail"
	default:
		return "unknown"
	}
}

// ParseMode parses a mode string to Mode
func ParseMode(s string) Mode {
	switch s {
	case "off":
		return ModeOff
	case "warn":
		return ModeWarn
	case "fail":
		return ModeFail
	default:
		return ModeWarn
	}
}

// Finding is a single gate finding
type Finding struct {
	Gate     string                 `json:"gate"`
	Severity string                 `json:"severity"`
	Message  string                 `json:"message"`
	Evidence map[string]interface{} `json:"evidence,omitempty"`
}

// Result is a gate check result
type Result struct {
	Gate       string     `json:"gate"`
	Target     string     `json:"target"`
	Pass       bool       `json:"pass"`
	DurationMs int64      `json:"duration_ms"`
	Findings   []Finding  `json:"findings,omitempty"`
}

// Recorder tracks gate results
type Recorder struct {
	modes   Mode
	results []Result
}

// NewRecorder creates a new recorder
func NewRecorder(mode Mode) *Recorder {
	return &Recorder{
		modes:   mode,
		results: make([]Result, 0),
	}
}

// Add records a gate result
func (r *Recorder) Add(result Result) error {
	r.results = append(r.results, result)
	
	// Only fail in ModeFail
	if !result.Pass && r.modes == ModeFail {
		return fmt.Errorf("gate %s failed on %s: %d findings", result.Gate, result.Target, len(result.Findings))
	}
	
	return nil
}

// Results returns all recorded results
func (r *Recorder) Results() []Result {
	return r.results
}

// Failed returns true if any gate failed
func (r *Recorder) Failed() bool {
	for _, res := range r.results {
		if !res.Pass {
			return true
		}
	}
	return false
}

// WriteJSON writes gate results to file
func (r *Recorder) WriteJSON(dir string) (string, error) {
	data, err := json.MarshalIndent(struct {
		Mode     string    `json:"mode"`
		Total    int       `json:"total"`
		Passed   int       `json:"passed"`
		Failed   int       `json:"failed"`
		Results  []Result  `json:"results"`
		Duration int64     `json:"duration_ms"`
	}{
		Mode:  r.modes.String(),
		Total: len(r.results),
		Passed: countPasses(r.results),
		Failed: countFails(r.results),
		Results: r.results,
	}, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal gate results: %w", err)
	}
	
	path := fmt.Sprintf("%s/gates.json", dir)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write gate results: %w", err)
	}
	
	return path, nil
}

func countPasses(results []Result) int {
	count := 0
	for _, r := range results {
		if r.Pass {
			count++
		}
	}
	return count
}

func countFails(results []Result) int {
	count := 0
	for _, r := range results {
		if !r.Pass {
			count++
		}
	}
	return count
}

// Gate is a gate check function
type Gate func() []Finding

// RunGate runs a gate and returns findings
func RunGate(name string, fn Gate) []Finding {
	findings := fn()
	
	if len(findings) > 0 {
		for i := range findings {
			findings[i].Gate = name
		}
	}
	
	return findings
}

// RecordGate records gate findings as a result
func RecordGate(rec *Recorder, name, target string, findings []Finding) {
	pass := len(findings) == 0
	rec.Add(Result{
		Gate:       name,
		Target:     target,
		Pass:       pass,
		DurationMs: 0, // would need timing from caller
		Findings:   findings,
	})
}

// ExpectedMissingMarker is used to mark intentional 404s
const ExpectedMissingMarker = "__expected-404"
