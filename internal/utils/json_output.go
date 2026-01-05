package utils

import (
	"github.com/sanjay-subramanya/drift/internal/core/model"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func WriteJSON(path, base string, findings []model.Finding) error {
	type JSONDetail struct {
		Message string `json:"Message"`
	}

	type JSONDrift struct {
		Drift struct {
			Base    string `json:"Base"`
			Summary string `json:"Summary"`
		} `json:"Drift"`
		Breakdown map[string]JSONDetail `json:"Breakdown,omitempty"`
	}

	var out []JSONDrift

	for _, f := range findings {
		j := JSONDrift{
			Breakdown: make(map[string]JSONDetail),
		}

		j.Drift.Base = base

		lines := strings.Split(f.Drift.Summary, "\n")
		if len(lines) > 0 {
			j.Drift.Summary = strings.TrimSuffix(lines[0], ";")
		}

		for _, line := range lines[1:] {
			line = strings.TrimSpace(line)

			switch {
			case strings.HasPrefix(line, "[CRITICAL]"):
				j.Breakdown["CRITICAL"] = JSONDetail{
					Message: strings.TrimSpace(strings.TrimPrefix(line, "[CRITICAL]")),
				}
			case strings.HasPrefix(line, "[HIGH]"):
				j.Breakdown["HIGH"] = JSONDetail{
					Message: strings.TrimSpace(strings.TrimPrefix(line, "[HIGH]")),
				}
			case strings.HasPrefix(line, "[LOW]"):
				j.Breakdown["LOW"] = JSONDetail{
					Message: strings.TrimSpace(strings.TrimPrefix(line, "[LOW]")),
				}
			}
		}

		out = append(out, j)
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	fmt.Println("json output recorded to", path)
	return os.WriteFile(path, data, 0644)
}