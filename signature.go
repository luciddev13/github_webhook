package github_webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Validates the signature sent by GitHub.
func validateSignature(secretEnvVarName string, rawPayload []byte, r *http.Request) error {

	secretToken := os.Getenv(secretEnvVarName)
	if secretToken == "" {
		// To secret given, skip this check
		return nil
	}

	rawSignature := r.Header.Get("X-Hub-Signature-256")
	if rawSignature == "" {
		return fmt.Errorf("signature is empty")
	}
	if !strings.HasPrefix(rawSignature, "sha1=") {
		return fmt.Errorf("malformed signature: %s", rawSignature[:10])
	}

	signature := []byte(rawSignature)

	hash := hmac.New(sha1.New, []byte(secretToken))

	if _, err := hash.Write(rawPayload); err != nil {
		return fmt.Errorf("could not calculate hash for payload: %s", err)
	}

	calculatedHash := []byte("sha1=" + hex.EncodeToString(hash.Sum(nil)))

	if 1 != subtle.ConstantTimeEq(int32(len(calculatedHash)), int32(len(signature))) {
		return fmt.Errorf("signature length mismatch")
	}

	if 1 != subtle.ConstantTimeCompare(calculatedHash, signature) {
		return fmt.Errorf("signature mismatch")
	}

	return nil
}
