package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithJSON padroniza a resposta JSON com status HTTP específico
func RespondWithJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Erro ao codificar resposta JSON: %v", err)
		}
	}
}

// RespondWithError padroniza mensagens de erro com status HTTP específico
func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, map[string]string{"error": message})
}

// RespondCreated padroniza resposta para recursos criados (201)
func RespondCreated(w http.ResponseWriter, data any) {
	RespondWithJSON(w, http.StatusCreated, data)
}

// RespondOK padroniza resposta de sucesso (200)
func RespondOK(w http.ResponseWriter, data any) {
	RespondWithJSON(w, http.StatusOK, data)
}

// RespondNoContent padroniza resposta sem conteúdo (204)
func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
