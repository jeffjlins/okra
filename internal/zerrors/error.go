package zerrors

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/emirpasic/gods/v2/sets/hashset"
)

type Error[T ~string] struct {
	code       T
	wrappedErr error
	tags       *hashset.Set[string]
	data       map[string]any
	stack      *stack
}

// New creates a new Error instance.
func New[T ~string](code T) *Error[T] {
	return &Error[T]{
		code:       code,
		wrappedErr: nil,
		data:       map[string]any{},
		tags:       hashset.New[string](),
		stack:      captureStack(2),
	}
}

// TODO: see comm [Structured Errors in Go](https://news.ycombinator.com/item?id=44148734)
func (e *Error[T]) With(k string, v any) *Error[T] {
	e.data[k] = v
	return e
}

// WithError wraps an existing error.
func (e *Error[T]) WithError(err error) *Error[T] {
	e.wrappedErr = err

	// Propagate the tags
	if err != nil {
		if wrappedErr, ok := err.(interface{ GetTags() []string }); ok {
			if tags := wrappedErr.GetTags(); len(tags) > 0 {
				e.tags.Add(tags...)
			}
		}
	}

	return e
}

// Error implements the error interface.
func (e *Error[T]) Error() string {
	if e == nil {
		return ""
	}
	if e.wrappedErr != nil {
		return fmt.Sprintf("%s: %s", e.code, e.wrappedErr.Error())
	}
	return string(e.code)
}

// // Errorf formats and wraps an error message.
// func (e *Error[T]) Errorf(format string, a ...any) *Error[T] {
// 	e.wrappedErr = fmt.Errorf(format, a...)
// 	return e
// }

// Unwrap implements error unwrapping.
func (e *Error[T]) Unwrap() error {
	return e.wrappedErr
}

func (e *Error[T]) Tags(tags ...string) *Error[T] {
	e.tags.Add(tags...)
	return e
}

func (e *Error[T]) HasTags(tags ...string) bool {
	return e.tags.Contains(tags...)
}

func (e *Error[T]) GetTags() []string {
	if e == nil {
		return []string{}
	}
	if e.tags == nil {
		return []string{}
	}
	return e.tags.Values()
}

func (e *Error[T]) Data() map[string]any {
	if e == nil || e.data == nil {
		return make(map[string]any)
	}
	result := make(map[string]any, len(e.data))
	for k, v := range e.data {
		result[k] = v
	}
	return result
}

func (e *Error[T]) Get(key string) (any, bool) {
	val, ok := e.data[key]
	return val, ok
}

func (e *Error[T]) StackTrace() string {
	if e == nil || e.stack == nil {
		return ""
	}
	return e.stack.String()
}

//nolint:ireturn // This is fine
func (e *Error[T]) Code() T {
	return e.code
}

func (e *Error[T]) CodeString() string {
	return string(e.code)
}

func HasCode[T ~string](err error, code T) bool {
	var e *Error[T]
	if errors.As(err, &e) {
		return e.Code() == code
	}
	return false
}

func (e *Error[T]) LogValue() slog.Value {
	if e == nil {
		at := []slog.Attr{slog.String("missing", "Error is missing")}
		return slog.GroupValue(at...)
	}
	// Create base attributes
	attrs := []slog.Attr{
		slog.String("code", string(e.code)),
		slog.String("error", e.Error()),
	}

	// Add data group if there's any custom data
	if len(e.data) > 0 {
		// Convert map entries directly to key-value pairs for slog.Group
		//nolint:mnd // 2 is the pair nr
		dataArgs := make([]any, 0, len(e.data)*2)
		for k, v := range e.data {
			dataArgs = append(dataArgs, k, v)
		}
		attrs = append(attrs, slog.Group("data", dataArgs...))
	}

	if !e.tags.Empty() {
		attrs = append(attrs, slog.Any("tags", e.GetTags()))
	}

	// Handle wrapped error
	if e.wrappedErr != nil {
		if logValuer, ok := e.wrappedErr.(slog.LogValuer); ok {
			attrs = append(attrs, slog.Any("wrapped", logValuer.LogValue()))
		} else {
			attrs = append(attrs, slog.String("wrapped", e.wrappedErr.Error()))
		}
	}

	if e.stack != nil {
		attrs = append(attrs, slog.String("stack", e.stack.String()))
	}

	return slog.GroupValue(attrs...)
}

// Is implements error comparison.
func (e *Error[T]) Is(target error) bool {
	t, ok := target.(*Error[T])
	if !ok {
		return false
	}
	return e.code == t.code
}

// As implements error casting.
func (e *Error[T]) As(target any) bool {
	if targetErr, ok := target.(**Error[T]); ok {
		*targetErr = e
		return true
	}

	if e.wrappedErr != nil {
		if asErr, ok := e.wrappedErr.(interface{ As(any) bool }); ok {
			return asErr.As(target)
		}
	}

	return false
}

// As implements error casting with a callback.
func As[T ~string, V any](err error, fn func(zerr *Error[T]) V) (*V, bool) {
	var zerr *Error[T]
	if errors.As(err, &zerr) {
		val := fn(zerr)
		return &val, true
	}
	var empty *V
	return empty, false
}
