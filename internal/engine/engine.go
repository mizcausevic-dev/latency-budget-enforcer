package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/mizcausevic-dev/latency-budget-enforcer/internal/model"
)

const SampleBudgetPath = "data/sample-budget.json"

func LoadSampleRequest() (model.BudgetRequest, error) {
	return LoadRequest(SampleBudgetPath)
}

func LoadRequest(path string) (model.BudgetRequest, error) {
	resolvedPath, err := resolvePath(path)
	if err != nil {
		return model.BudgetRequest{}, err
	}

	payload, err := os.ReadFile(resolvedPath)
	if err != nil {
		return model.BudgetRequest{}, err
	}

	var request model.BudgetRequest
	if err := json.Unmarshal(payload, &request); err != nil {
		return model.BudgetRequest{}, err
	}

	return request, nil
}

func resolvePath(path string) (string, error) {
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("could not resolve sample path %q", path)
	}

	candidate := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", path))
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}

	return "", fmt.Errorf("could not resolve sample path %q", path)
}

func Evaluate(request model.BudgetRequest) model.Report {
	findings := make([]model.Finding, 0, 5)

	if finding := sustainedBreach(request); finding != nil {
		findings = append(findings, *finding)
	}
	if finding := p95Drift(request); finding != nil {
		findings = append(findings, *finding)
	}
	if finding := p99Breach(request); finding != nil {
		findings = append(findings, *finding)
	}
	if finding := dependencyDrag(request); finding != nil {
		findings = append(findings, *finding)
	}
	if finding := errorPressureCoupling(request); finding != nil {
		findings = append(findings, *finding)
	}

	overall := 18
	for _, finding := range findings {
		if finding.Score > overall {
			overall = finding.Score
		}
	}

	sort.Slice(findings, func(i, j int) bool {
		return findings[i].Score > findings[j].Score
	})

	return model.Report{
		SystemName:           request.SystemName,
		ServiceOwner:         request.ServiceOwner,
		ObservationsAnalyzed: len(request.Observations),
		OverallScore:         overall,
		Findings:             findings,
	}
}

func sustainedBreach(request model.BudgetRequest) *model.Finding {
	evidence := make([]string, 0)
	count := 0
	for _, observation := range request.Observations {
		if observation.P95MS > request.P95BudgetMS {
			count++
			evidence = append(evidence, fmt.Sprintf("%s %s p95=%dms budget=%dms", observation.Region, observation.Path, observation.P95MS, request.P95BudgetMS))
		}
	}
	if count < 4 {
		return nil
	}

	return &model.Finding{
		Code:     "sustained_breach",
		Severity: "critical",
		Score:    92,
		Summary:  "The primary service path is breaching its latency budget often enough to require coordinated action.",
		Evidence: evidence,
		RecommendedNextAction: "Stabilize the dominant path before more traffic is routed through it and confirm whether a rollback, cache change, or dependency cut is needed.",
	}
}

func p95Drift(request model.BudgetRequest) *model.Finding {
	evidence := make([]string, 0)
	for _, observation := range request.Observations {
		drift := observation.P95MS - request.P95BudgetMS
		if drift >= 40 {
			evidence = append(evidence, fmt.Sprintf("%s %s drift=%dms", observation.Region, observation.Path, drift))
		}
	}
	if len(evidence) == 0 {
		return nil
	}

	return &model.Finding{
		Code:     "p95_drift",
		Severity: "high",
		Score:    81,
		Summary:  "Normal-path latency is drifting far enough above the budget to erode experience under ordinary load.",
		Evidence: evidence,
		RecommendedNextAction: "Profile the critical handlers and remove latency overhead on the p95 path before the breach becomes the new baseline.",
	}
}

func p99Breach(request model.BudgetRequest) *model.Finding {
	evidence := make([]string, 0)
	for _, observation := range request.Observations {
		if observation.P99MS > request.P99BudgetMS {
			evidence = append(evidence, fmt.Sprintf("%s %s p99=%dms budget=%dms", observation.Region, observation.Path, observation.P99MS, request.P99BudgetMS))
		}
	}
	if len(evidence) < 3 {
		return nil
	}

	return &model.Finding{
		Code:     "p99_breach",
		Severity: "high",
		Score:    84,
		Summary:  "Tail-latency behavior is beyond the agreed guardrail and is starting to shape system trust.",
		Evidence: evidence,
		RecommendedNextAction: "Contain the p99 spikes by reviewing queueing, cold-start behavior, and concurrency pressure on the worst path.",
	}
}

func dependencyDrag(request model.BudgetRequest) *model.Finding {
	evidence := make([]string, 0)
	for _, observation := range request.Observations {
		if observation.DependencyMS >= 150 {
			evidence = append(evidence, fmt.Sprintf("%s %s dependency=%dms", observation.Region, observation.Path, observation.DependencyMS))
		}
	}
	if len(evidence) == 0 {
		return nil
	}

	return &model.Finding{
		Code:     "dependency_drag",
		Severity: "moderate",
		Score:    67,
		Summary:  "A downstream layer is contributing enough latency to distort the overall service path.",
		Evidence: evidence,
		RecommendedNextAction: "Isolate the slow dependency and confirm whether caching, fan-out reduction, or regional routing changes are needed.",
	}
}

func errorPressureCoupling(request model.BudgetRequest) *model.Finding {
	evidence := make([]string, 0)
	for _, observation := range request.Observations {
		if observation.ErrorRate > request.MaxErrorRate && observation.P95MS > request.P95BudgetMS {
			evidence = append(evidence, fmt.Sprintf("%s %s error_rate=%.1f%% p95=%dms", observation.Region, observation.Path, observation.ErrorRate, observation.P95MS))
		}
	}
	if len(evidence) == 0 {
		return nil
	}

	return &model.Finding{
		Code:     "error_pressure_coupling",
		Severity: "high",
		Score:    86,
		Summary:  "Latency and elevated error pressure are rising together on the same path, increasing user-visible failure risk.",
		Evidence: evidence,
		RecommendedNextAction: "Treat the path as an incident lane and coordinate rollback, dependency review, and traffic shaping together.",
	}
}
