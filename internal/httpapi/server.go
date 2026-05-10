package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/mizcausevic-dev/latency-budget-enforcer/internal/engine"
	"github.com/mizcausevic-dev/latency-budget-enforcer/internal/model"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/sample", sampleHandler)
	mux.HandleFunc("/api/evaluate", evaluateHandler)
	return mux
}

func rootHandler(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, model.HealthResponse{
		Status:       "ok",
		Service:      "latency-budget-enforcer",
		Docs:         "/api/sample",
		SampleBudget: engine.SampleBudgetPath,
	})
}

func healthHandler(writer http.ResponseWriter, _ *http.Request) {
	rootHandler(writer, nil)
}

func sampleHandler(writer http.ResponseWriter, _ *http.Request) {
	request, err := engine.LoadSampleRequest()
	if err != nil {
		writeJSON(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(writer, http.StatusOK, engine.Evaluate(request))
}

func evaluateHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeJSON(writer, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var payload model.BudgetRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(writer, http.StatusOK, engine.Evaluate(payload))
}

func writeJSON(writer http.ResponseWriter, statusCode int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_ = json.NewEncoder(writer).Encode(payload)
}
