package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stevenwijaya/finance-tracker/pkg/log"
)

// data yang ingin di-mask
var sensitiveFields = []string{
	"password",
	"email",
	"token",
	"authorization",
	"refresh_token",
}

// helper masking JSON
func maskSensitive(data map[string]interface{}) map[string]interface{} {
	for key, value := range data {
		lKey := strings.ToLower(key)

		// Mask exact match
		for _, field := range sensitiveFields {
			if lKey == field {
				data[key] = "***MASKED***"
			}
		}

		// Nested
		if nested, ok := value.(map[string]interface{}); ok {
			data[key] = maskSensitive(nested)
		}

		// Array
		if arr, ok := value.([]interface{}); ok {
			for i, v := range arr {
				if obj, ok := v.(map[string]interface{}); ok {
					arr[i] = maskSensitive(obj)
				}
			}
			data[key] = arr
		}
	}
	return data
}

// untuk menangkap response body
type responseRecorder struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		// ============= REQUEST LOGGING =============
		var reqBody map[string]interface{}
		if c.Request.Body != nil {
			raw, _ := io.ReadAll(c.Request.Body)
			if len(raw) > 0 {
				_ = json.Unmarshal(raw, &reqBody)
				reqBody = maskSensitive(reqBody)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))
		}

		// log request seperti style kamu sekarang
		log.Info("[HTTP]",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"request", reqBody,
		)

		// =============== RESPONSE CAPTURE ===============
		recorder := &responseRecorder{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = recorder

		// Process request
		c.Next()

		// calculate time
		latency := time.Since(start)

		// decode & masking response body
		var resBody map[string]interface{}
		if recorder.body.Len() > 0 {
			_ = json.Unmarshal(recorder.body.Bytes(), &resBody)
			resBody = maskSensitive(resBody)
		}

		// log response tetap clean
		jsonStr, _ := json.Marshal(resBody)

		log.Info("[HTTP][RESPONSE]",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", recorder.Status(),
			"latency", latency.String(),
			"response", string(jsonStr),
		)
	}
}
