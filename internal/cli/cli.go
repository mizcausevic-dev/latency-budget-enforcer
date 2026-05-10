package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mizcausevic-dev/latency-budget-enforcer/internal/engine"
)

func Run() error {
	mode := flag.String("mode", "server", "server or cli")
	input := flag.String("input", engine.SampleBudgetPath, "path to budget JSON")
	format := flag.String("format", "text", "text or json")
	flag.Parse()

	if *mode == "server" {
		return nil
	}

	request, err := engine.LoadRequest(*input)
	if err != nil {
		return err
	}

	report := engine.Evaluate(request)
	if *format == "json" {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(report)
	}

	fmt.Println("Latency Budget Enforcer")
	fmt.Println("=======================")
	fmt.Printf("System: %s\n", report.SystemName)
	fmt.Printf("Observations analyzed: %d\n", report.ObservationsAnalyzed)
	fmt.Printf("Overall score: %d\n\n", report.OverallScore)

	for _, finding := range report.Findings {
		fmt.Printf("[%s] %s (score %d)\n", toUpper(finding.Severity), finding.Code, finding.Score)
		fmt.Printf("Summary: %s\n", finding.Summary)
		fmt.Println("Evidence:")
		for _, item := range finding.Evidence {
			fmt.Printf("  - %s\n", item)
		}
		fmt.Printf("Next action: %s\n\n", finding.RecommendedNextAction)
	}

	return nil
}

func toUpper(value string) string {
	switch value {
	case "critical":
		return "CRITICAL"
	case "high":
		return "HIGH"
	case "moderate":
		return "MODERATE"
	default:
		return value
	}
}
