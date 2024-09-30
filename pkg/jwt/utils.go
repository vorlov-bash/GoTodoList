package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
)

func toBase64(data map[string]any) string {
	dataJson, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	return base64.URLEncoding.EncodeToString(dataJson)
}

func generateHeader(alg string, typ string) string {
	header := map[string]any{
		"alg": alg,
		"typ": typ,
	}
	return toBase64(header)
}

func generatePayload(sub string, exp int64) string {
	payload := map[string]any{
		"sub": sub,
		"exp": exp,
	}
	return toBase64(payload)
}

func generateSignature(header string, payload string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(header + "." + payload))
	return string(h.Sum(nil))
}

func getPayload(payload string) (map[string]interface{}, error) {
	decodedPayload, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	var payloadMap map[string]interface{}
	err = json.Unmarshal(decodedPayload, &payloadMap)
	if err != nil {
		return nil, err
	}
	return payloadMap, nil
}
