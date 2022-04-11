package handlers

import (
	"backend/models"
	"encoding/json"
	"net/http"
)

func (m *Repository) StatusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := models.AppStatus{
		Status:      "Available",
		Environment: m.App.Config.Env,
		Version:     m.App.Config.Version,
	}

	js, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		m.App.Logger.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
