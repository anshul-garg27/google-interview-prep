package helpers

import (
	"bytes"
	"coffee/core/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/spf13/viper"
)

func ToInt64(i int64) *int64 {
	return &i
}

func ToFloat64(i float64) *float64 {
	return &i
}

func ToString(s string) *string {
	return &s
}

func ToBool(b bool) *bool {
	return &b
}

func ParseToInt64(s string) *int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}

	return &i
}
func ParseToStringFromInt(i int) string {
	return strconv.Itoa(i)
}

func GetDaySuffix(day int) string {
	suffix := "th"
	switch day {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}
	return fmt.Sprintf("%d%s", day, suffix)
}

func EncodeEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // Invalid email format, return original string
	}

	username := parts[0]
	domain := parts[1]

	encodedUsername := strings.Repeat("*", len(username))
	encodedEmail := fmt.Sprintf("%s@%s", encodedUsername, domain)

	return encodedEmail
}

func EncodePhone(phone string) string {
	encodedPhone := "+91 " + "**********"
	return encodedPhone
}

// StringToTimestamp converts a Unix timestamp string to a time.Time value.
// The timestampString parameter should be a string representation of the Unix timestamp.
// If the conversion is successful, it returns the corresponding time.Time value.
// If an error occurs during conversion, it returns a zero time.Time value and the error.

func StringToTimestamp(timestampString string) (time.Time, error) {
	timestamp, err := strconv.ParseInt(timestampString, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	timeValue := time.Unix(timestamp, 0)
	return timeValue, nil
}

type ResponseWriterWithContent struct {
	http.ResponseWriter
	ResponseBodyCopy *bytes.Buffer
	InputBodyCopy    []byte
}

func (rw *ResponseWriterWithContent) Write(b []byte) (int, error) {
	// Write the response body
	rw.ResponseBodyCopy.Write(b)
	return rw.ResponseWriter.Write(b)
}

func GetSentryEvent(rw *ResponseWriterWithContent, request *http.Request, requestBody []byte) {
	message := domain.ErrorLog{
		Headers: make(map[string][]string),
	}
	message.RequestUrl = request.URL.String()
	var jsonResponse map[string]interface{}
	if !strings.Contains(message.RequestUrl, "/heartbeat/") && !strings.Contains(message.RequestUrl, "/metrics") {
		if err := json.Unmarshal(rw.ResponseBodyCopy.Bytes(), &jsonResponse); err == nil {
			if statusObj, ok := jsonResponse["status"].(map[string]interface{}); ok {
				statusCode, ok := statusObj["status"].(string)
				if ok && (statusCode == "ERROR" || statusCode == "500") {
					message.Type = statusCode
					errorMessage, ok := statusObj["message"].(string)
					if ok {
						message.Message = errorMessage
					}
					// Read the request body
					if len(requestBody) > 0 {
						message.RequestBody = string(requestBody)
					}
					message.RequestType = request.Method
					for name, values := range request.Header {
						message.Headers[name] = values
					}
					event := makeEvent(message)
					sentry.CaptureEvent(event)
				}
			}
		} else {
			message.Message = rw.ResponseBodyCopy.String()

			if len(requestBody) > 0 {
				message.RequestBody = string(requestBody)
			}
			message.RequestType = request.Method
			for name, values := range request.Header {
				message.Headers[name] = values
			}
			message.Type = "ERRROR"
			event := makeEvent(message)
			sentry.CaptureEvent(event)
		}
	}
}
func makeEvent(message domain.ErrorLog) *sentry.Event {
	event := sentry.NewEvent()
	event.Level = sentry.LevelError
	event.Environment = viper.GetString("ENV")
	event.Message = message.Message

	// Set tags
	event.Tags["type"] = message.Type
	event.Extra["request_type"] = message.RequestType
	event.Extra["request_url"] = message.RequestUrl
	event.Extra["request_body"] = message.RequestBody
	event.Extra["headers"] = message.Headers
	return event
}

func GetRequestBody(r *http.Request) []byte {
	body, err := io.ReadAll(r.Body)
	if err == nil {
		if len(body) > 0 {
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	}
	return body
}
