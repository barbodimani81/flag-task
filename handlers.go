package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Post("/flags", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name         string   `json:"name"`
			Dependencies []string `json:"dependencies"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err := createFeatureFlag(db, req.Name, req.Dependencies)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	r.Post("/flags/{name}/toggle", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")

		var req struct {
			Enable bool   `json:"enable"`
			Actor  string `json:"actor"`
			Reason string `json:"reason"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		err := toggleFlag(db, name, req.Enable, req.Actor, req.Reason)
		if err != nil {
			if err.Error()[:24] == "missing active dependenc" {
				http.Error(w, err.Error(), http.StatusConflict)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/flags/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		var flag FeatureFlag
		err := db.QueryRow(`SELECT id, name, enabled, created_at, updated_at FROM feature_flags WHERE name = ?`, name).
			Scan(&flag.ID, &flag.Name, &flag.Enabled, &flag.CreatedAt, &flag.UpdatedAt)
		if err != nil {
			http.Error(w, "flag not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(flag)
	})

	r.Get("/flags/{name}/logs", func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		var flagID int
		err := db.QueryRow(`SELECT id FROM feature_flags WHERE name = ?`, name).Scan(&flagID)
		if err != nil {
			http.Error(w, "flag not found", http.StatusNotFound)
			return
		}
		rows, err := db.Query(`SELECT action, actor, reason, timestamp FROM audit_logs WHERE flag_id = ? ORDER BY timestamp DESC`, flagID)
		if err != nil {
			http.Error(w, "failed to get logs", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var logs []AuditLog
		for rows.Next() {
			var log AuditLog
			log.FlagID = flagID
			rows.Scan(&log.Action, &log.Actor, &log.Reason, &log.Timestamp)
			logs = append(logs, log)
		}
		json.NewEncoder(w).Encode(logs)
	})

	return r
}
