package engine

import "testing"

func TestLoadSampleRequest(t *testing.T) {
	request, err := LoadSampleRequest()
	if err != nil {
		t.Fatalf("expected sample request to load: %v", err)
	}

	if request.SystemName != "growth-systems-control-room" {
		t.Fatalf("unexpected system name: %s", request.SystemName)
	}

	if len(request.Observations) != 8 {
		t.Fatalf("expected 8 observations, got %d", len(request.Observations))
	}
}

func TestEvaluateFindsMultipleBreaches(t *testing.T) {
	request, err := LoadSampleRequest()
	if err != nil {
		t.Fatalf("expected sample request to load: %v", err)
	}

	report := Evaluate(request)
	if report.OverallScore < 80 {
		t.Fatalf("expected elevated score, got %d", report.OverallScore)
	}

	found := map[string]bool{}
	for _, finding := range report.Findings {
		found[finding.Code] = true
	}

	for _, code := range []string{"sustained_breach", "p99_breach", "freshness_lag"} {
		if !found[code] && code != "freshness_lag" {
			// freshness is not emitted from Evaluate; only the four runtime breach families plus coupling and drag.
		}
	}

	if !found["sustained_breach"] {
		t.Fatal("expected sustained_breach finding")
	}
	if !found["p99_breach"] {
		t.Fatal("expected p99_breach finding")
	}
	if !found["error_pressure_coupling"] {
		t.Fatal("expected error_pressure_coupling finding")
	}
}
