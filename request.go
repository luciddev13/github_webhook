package github_webhook

import (
	"log"
	"net/http"
)

func newRequestHandler(handlerFunc HandlerFunc, secretEnvVarName string, maxPayloadSize int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawPayload, err := loadRawPayload(r, maxPayloadSize)
		if err != nil {
			log.Printf("failed to load payload, reason: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = validateSignature(secretEnvVarName, rawPayload, r)
		if err != nil {
			log.Printf("failed signature validation, reason: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		payload, err := parsePayload(rawPayload)
		if err != nil {
			log.Printf("failed to parse payload, reason: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = handlerFunc(payload); err != nil {
			log.Printf("failed to handle request, reason: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
