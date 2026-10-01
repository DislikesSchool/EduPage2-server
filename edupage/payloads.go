package edupage

import (
	"bytes"
	"compress/flate"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// User-Agent used for all requests to EduPage.
// EduPage's PHP backend + Cloudflare blocks the default Go UA.
const EduPageUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

type MessagePayload struct {
	SelectedUser string
	Text         string
	Attachments  string
	Typ          string
}

type CanteenPayload struct {
	BoarderID   string            `json:"stravnikid"`
	Edupage     string            `json:"edupage"`
	BoarderUser string            `json:"stravnikUser,omitempty"`
	Date        string            `json:"mysqlDate"` // YYYY-mm-dd
	FIDS        map[string]string `json:"jids"`
	View        string            `json:"view"`
	Permission  string            `json:"pravo"`
	Action      string            `json:"akcia"`
}

func CreateMessage(receiver, text, attachments string) MessagePayload {
	return MessagePayload{
		SelectedUser: receiver,
		Text:         text,
		Attachments:  attachments,
		Typ:          "sprava",
	}
}

// pythonQuote mimics urllib.parse.quote(s, safe='/').
// It percent-encodes UTF-8 bytes, leaving [A-Za-z0-9-_.~/] unescaped.
func pythonQuote(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~' || c == '/' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

// encodeFormData mimics Python's ModuleHelper.encode_form_data /
// RequestUtil.encode_form_data: quote(key)=quote(value) joined by &.
func encodeFormData(data map[string]string) string {
	// Sort keys for determinism (not required, but stable for tests).
	// Preserve insertion order is not needed; server accepts any order.
	parts := make([]string, 0, len(data))
	for k, v := range data {
		parts = append(parts, pythonQuote(k)+"="+pythonQuote(v))
	}
	// NOTE: map iteration is random; sort to keep hashes stable in logs.
	// We don't sort the actual request semantics, but it doesn't matter.
	return strings.Join(parts, "&")
}

// rawDeflate compresses data with raw DEFLATE (wbits=-15),
// matching Python's zlib.compressobj(-1, DEFLATED, -15).
func rawDeflate(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodeRequestBody mirrors edupage-api's RequestData.encode_request_body:
//  1. urlencode dict
//  2. raw-deflate + standard base64, prefix "dz:"
//  3. eqacs = sha1 hex of eqap, eqaz = "1"
//  4. urlencode {eqap, eqacs, eqaz}
//
// This is exactly what /global/pics/js/edubarUtils.js does for eqav=1:
//
//	var gz = new Zlib.RawDeflate(encoder.encode(cs));
//	cs0 = 'dz:'+btoa(cs1);
//	obj = { eqap: cs0, eqacs: sha1(cs0), eqaz: '1' }
func EncodeRequestBody(data map[string]string) (string, error) {
	inner := encodeFormData(data)
	compressed, err := rawDeflate([]byte(inner))
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(compressed)
	eqap := "dz:" + b64
	h := sha1.Sum([]byte(eqap))
	eqacs := fmt.Sprintf("%x", h)
	outer := encodeFormData(map[string]string{
		"eqap":  eqap,
		"eqacs": eqacs,
		"eqaz":  "1",
	})
	return outer, nil
}

// CreatePayload builds url.Values {eqap, eqacs, eqaz} using the modern
// dz:+raw-deflate+sha1-hex encoding. Kept signature for compatibility.
func CreatePayload(data map[string]string) url.Values {
	body, err := EncodeRequestBody(data)
	if err != nil {
		// Fallback: should never happen; return plain form.
		values := url.Values{}
		for k, v := range data {
			values.Add(k, v)
		}
		return values
	}
	// body is already url-encoded "eqap=..&eqacs=..&eqaz=1", parse back
	// to url.Values so callers can use PostForm.
	vals, err := url.ParseQuery(body)
	if err != nil {
		values := url.Values{}
		values.Add("eqap", "")
		return values
	}
	return vals
}

// DecodeResponseBody handles EduPage's eqav envelope:
//   - "eqwd:..." -> wrong data (server asks to retry with eqav+1)
//   - "eqz:..."  -> standard base64 (UTF-8 JSON)
//   - otherwise  -> raw body as-is
func DecodeResponseBody(body []byte) ([]byte, error) {
	s := string(body)
	if strings.HasPrefix(s, "eqwd:") {
		return nil, fmt.Errorf("eqwd: wrong data (server asked to retry)")
	}
	if strings.HasPrefix(s, "eqz:") {
		enc := strings.TrimSpace(s[4:])
		// Server uses standard base64 (btoa/atob), not URL-safe.
		// Try Std first, then URL-safe, then raw.
		if decoded, err := base64.StdEncoding.DecodeString(enc); err == nil {
			return decoded, nil
		}
		if decoded, err := base64.URLEncoding.DecodeString(enc); err == nil {
			return decoded, nil
		}
		// Some old payloads have newlines; strip whitespace and retry.
		clean := strings.ReplaceAll(strings.ReplaceAll(enc, "\n", ""), "\r", "")
		if decoded, err := base64.StdEncoding.DecodeString(clean); err == nil {
			return decoded, nil
		}
		return nil, fmt.Errorf("failed to base64-decode eqz: response")
	}
	// Legacy EduPage2-server assumed first 4 bytes are a prefix like "eqz:".
	// Keep compatibility: if body looks like base64 with 4-char prefix, try it.
	if len(s) > 4 && (strings.HasPrefix(s, "eqz:") || strings.HasPrefix(s, "eqap")) {
		return body, nil
	}
	return body, nil
}

// DecodeLegacyBody is used by timeline/results/canteen fetchers that
// historically did base64.StdEncoding.Decode(body[4:]).
// It now supports both eqz: envelopes and plain JSON.
func DecodeLegacyBody(body []byte) ([]byte, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	s := strings.TrimSpace(string(body))
	// Modern $j.post with eqav=1 returns "eqz:"+base64
	if strings.HasPrefix(s, "eqz:") || strings.HasPrefix(s, "eqwd:") {
		return DecodeResponseBody([]byte(s))
	}
	// Plain JSON (new timeline/notifications endpoint returns JSON directly)
	if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
		return []byte(s), nil
	}
	// Legacy: first 4 bytes are prefix, rest is base64 (standard).
	// Try standard, then URL-safe.
	if len(body) > 4 {
		payload := strings.TrimSpace(s[4:])
		if decoded, err := base64.StdEncoding.DecodeString(payload); err == nil {
			// Trim NUL padding like old code did.
			decoded = bytes.Trim(decoded, "\x00")
			return decoded, nil
		}
		if decoded, err := base64.URLEncoding.DecodeString(payload); err == nil {
			decoded = bytes.Trim(decoded, "\x00")
			return decoded, nil
		}
	}
	// Last resort: return as-is and let JSON parser fail with context.
	return body, nil
}

// ParseRPCResponse decodes an eqz:/raw RPC response and unmarshals JSON.
func ParseRPCResponse(text string) (map[string]interface{}, error) {
	decoded, err := DecodeResponseBody([]byte(text))
	if err != nil {
		return nil, err
	}
	var out map[string]interface{}
	if err := json.Unmarshal(decoded, &out); err != nil {
		// Some RPC responses are wrapped differently; try raw text.
		if err2 := json.Unmarshal([]byte(text), &out); err2 == nil {
			return out, nil
		}
		return nil, err
	}
	return out, nil
}
