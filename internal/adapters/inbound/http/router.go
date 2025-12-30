package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/jeffjlins/okra/internal/usecase"
	"github.com/jeffjlins/okra/internal/zerrors"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) *zerrors.Error[HttpError]

func Make(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			var e *zerrors.Error[HttpError]
			if !errors.As(err, &e) {
				e = zerrors.New(ErrHttp).With("status_code", http.StatusInternalServerError)
			}
			logError(e)
			writeErrorJSON(w, e)
			return
		}
	}
}

func logError(e *zerrors.Error[HttpError]) {
	useJSON := os.Getenv("LOG_FORMAT") == "json"
	if useJSON {
		logErrorJSON(e)
	} else {
		logErrorText(e)
	}
}

func logErrorJSON(e *zerrors.Error[HttpError]) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	if wrapped := e.Unwrap(); wrapped != nil {
		logger = logger.With("wrapped", wrapped.Error())
	}

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
		if hasStackTrace, ok := wrapped.(interface{ StackTrace() string }); ok {
			if stackTrace := hasStackTrace.StackTrace(); stackTrace != "" {
				logger = logger.With("wrapped_stack", stackTrace)
			}
		}
	}

	logger.Error("HTTP API Error", "code", e.CodeString(), "message", e.Error())
}

func logErrorText(e *zerrors.Error[HttpError]) {
	fmt.Fprintf(os.Stderr, "[TimestampGoesHere | ERROR | HTTPError] %q", e.Error())

	if tags := e.GetTags(); len(tags) > 0 {
		fmt.Fprintf(os.Stderr, " | tags=%v", tags)
	}

	if data := e.Data(); len(data) > 0 {
		fmt.Fprintf(os.Stderr, " | data=%v", data)
	}

	fmt.Fprintf(os.Stderr, "\n")

	if stackTrace := e.StackTrace(); stackTrace != "" {
		fmt.Fprintf(os.Stderr, "Stack trace:\n%s\n", stackTrace)
	}

	//TODO: list stack traces one at a time with the isolated code for that error and it's data map instead of putting the data map in the first message but keep the tags in the first message
	if wrapped := e.Unwrap(); wrapped != nil {
		if hasStackTrace, ok := wrapped.(interface{ StackTrace() string }); ok {
			if stackTrace := hasStackTrace.StackTrace(); stackTrace != "" {
				fmt.Fprintf(os.Stderr, "Wrapped stack trace:\n%s\n", stackTrace)
			}
		}
	}
}

type errorResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Data    map[string]any     `json:"data,omitempty"`
	Tags    []string           `json:"tags,omitempty"`
	Wrapped []wrappedErrorItem `json:"wrapped,omitempty"`
}

type wrappedErrorItem struct {
	Code string         `json:"code"`
	Data map[string]any `json:"data,omitempty"`
}

func getErrorCodeOrMessage(err error) string {
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

func writeErrorJSON(w http.ResponseWriter, e *zerrors.Error[HttpError]) {
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
		Code: getErrorCodeOrMessage(err),
	}

	if hasData, ok := err.(interface{ Data() map[string]any }); ok {
		if data := hasData.Data(); len(data) > 0 {
			item.Data = make(map[string]any, len(data))
			for k, v := range data {
				item.Data[k] = v
			}
		}
	}

	*items = append(*items, item)

	if wrapped := errors.Unwrap(err); wrapped != nil {
		buildWrappedErrorsRecursive(wrapped, items)
	}
}

func NewRouter(uomService *usecase.UomService) *http.ServeMux {
	mux := http.NewServeMux()

	var ch CustomHandler = CustomHandler{uomService}
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /uom", Make(ch.createUomHandler))
	mux.HandleFunc("GET /uom/{id}", getUomByIDHandler(uomService))
	mux.HandleFunc("GET /uom", getAllUomsHandler(uomService))
	mux.HandleFunc("DELETE /uom/{id}", deleteUomHandler(uomService))
	mux.HandleFunc("PUT /uom/{id}", updateUomHandler(uomService))

	return mux
}
