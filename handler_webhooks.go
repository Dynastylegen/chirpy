package main

import (
	"encoding/json"
	"net/http"

	"github.com/Dynastylegen/chirpy/internal/auth"
	"github.com/google/uuid"
)

type WebHook struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		http.Error(w, "apiKey not Found", http.StatusUnauthorized)
		return
	}
	if apiKey != cfg.polkaKey {
		http.Error(w, "apiKey not Found", http.StatusUnauthorized)
		return
	}

	var params WebHook
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	err = cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		http.Error(w, "User not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
