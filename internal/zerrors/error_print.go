package zerrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func (e *Error[T]) SlogErrorJSON(w io.Writer) {
	logger := slog.New(slog.NewJSONHandler(w, nil))

	if tags := e.GetTags(); len(tags) > 0 {
		logger = logger.With("tags", tags)
	}

	if data := e.Data(); len(data) > 0 {
		logger = logger.With("data", data)
	}

	if stackTrace := e.StackTrace(); stackTrace != "" {
		logger = logger.With("stack", stackTrace)
	}

	if wrapped := e.Unwrap(); wrapped != nil {
		wrappedErrors := buildWrappedErrors(wrapped)
		if len(wrappedErrors) > 0 {
			logger = logger.With("wrapped", wrappedErrors)
		}
	}

	logger.Error("HTTP API Error", "code", e.CodeString(), "message", e.Error())
}

func (e *Error[T]) LogErrorText(w io.Writer) {
	timestamp := time.Now().Format(time.RFC3339)
	errorType := reflect.TypeOf(e.Code()).Name()
	fmt.Fprintf(w, "[%s | ERROR | %s] %q", timestamp, errorType, e.Error())

	if tags := e.GetTags(); len(tags) > 0 {
		fmt.Fprintf(w, " | tags=%v", tags)
	}

	fmt.Fprintf(w, "\n")

	wrappedItems := buildWrappedErrors(e)
	for seq, item := range wrappedItems {
		prefix := "Caused By: "
		if seq == 0 {
			prefix = "Error: "
		}
		fmt.Fprintf(w, "%s%s\n", prefix, item.Code)

		if len(item.Data) > 0 {
			fmt.Fprintf(w, "  Data:\n")
			for k, v := range item.Data {
				fmt.Fprintf(w, "    %s=%v\n", k, v)
			}
		}
		if item.StackTrace != "" {
			fmt.Fprintf(w, "  Stack trace:%s\n", item.StackTrace)
		}
	}
}

func (e *Error[T]) WriteErrorHttpJSON(w http.ResponseWriter) {
	var code int
	if c, ok := e.Get("status_code"); ok {
		code = c.(int)
	} else {
		code = http.StatusInternalServerError
	}

	response := errorResponse{
		Code:    e.CodeString(),
		Message: e.Error(),
		Tags:    e.GetTags(),
	}

	if len(e.GetTags()) == 0 {
		response.Tags = nil
	}

	data := e.Data()
	delete(data, "status_code")
	if len(data) > 0 {
		response.Data = data
	}

	if wrapped := e.Unwrap(); wrapped != nil {
		response.Wrapped = buildWrappedErrors(wrapped)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

type errorResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Data    map[string]any     `json:"data,omitempty"`
	Tags    []string           `json:"tags,omitempty"`
	Wrapped []wrappedErrorItem `json:"wrapped,omitempty"`
}

type wrappedErrorItem struct {
	Code       string         `json:"code"`
	Data       map[string]any `json:"data,omitempty"`
	StackTrace string         `json:"stack_trace,omitempty"`
}

func buildWrappedErrors(err error) []wrappedErrorItem {
	var items []wrappedErrorItem
	buildWrappedErrorsRecursive(err, &items)
	return items
}

func buildWrappedErrorsRecursive(err error, items *[]wrappedErrorItem) {
	if err == nil {
		return
	}

	item := wrappedErrorItem{
		Code: getErrorCodeOrOwnMessage(err),
	}

	if hasData, ok := err.(interface{ Data() map[string]any }); ok {
		if data := hasData.Data(); len(data) > 0 {
			item.Data = make(map[string]any, len(data))
			for k, v := range data {
				item.Data[k] = v
			}
		}
	}

	if hasStackTrace, ok := err.(interface{ StackTrace() string }); ok {
		if stackTrace := hasStackTrace.StackTrace(); stackTrace != "" {
			item.StackTrace = stackTrace
		}
	}

	*items = append(*items, item)

	if wrapped := errors.Unwrap(err); wrapped != nil {
		buildWrappedErrorsRecursive(wrapped, items)
	}
}

func getErrorCodeOrOwnMessage(err error) string {
	if hasCode, ok := err.(interface{ CodeString() string }); ok {
		return hasCode.CodeString()
	}
	return getErrorOwnMessage(err)
}

func getErrorOwnMessage(err error) string {
	if err == nil {
		return ""
	}

	msg := err.Error()
	wrapped := errors.Unwrap(err)

	if wrapped != nil {
		wrappedMsg := wrapped.Error()
		if len(wrappedMsg) < len(msg) && strings.HasSuffix(msg, wrappedMsg) {
			prefix := msg[:len(msg)-len(wrappedMsg)]
			prefix = strings.TrimSuffix(prefix, ": ")
			prefix = strings.TrimSuffix(prefix, ":")
			prefix = strings.TrimSpace(prefix)
			if len(prefix) > 0 {
				return prefix
			}
		}
	}

	return msg
}
