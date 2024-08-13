package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/b9uu/realty/internal/data"
	"github.com/b9uu/realty/internal/validator"
)

// Reads from r into dest. only 1MB allowed
func (app *application) ReadJSON(w http.ResponseWriter, r *http.Request, dest interface{}) error {
	const max = 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(max))

	dc := json.NewDecoder(r.Body)
	dc.DisallowUnknownFields()

	err := dc.Decode(dest)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		// returning helpful message based on the error
		switch {
		// syntax error
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		// badly formatted
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		// empty body
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		// unknown field
		case strings.HasPrefix(err.Error(), "json: unkown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unkdown field")
			return fmt.Errorf("body contains unknow key %s", fieldName)
		// check if the body too large
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", max)
		default:
			return err
		}
	}
	// check if there's multiple JSON values
	err = dc.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

// writes data to w with headers
func (app *application) writeJSON(w http.ResponseWriter, status int, data data.Envelope, headers http.Header) error {
	jsn, err := json.Marshal(data)
	if err != nil {
		return err
	}
	jsn = append(jsn, '\n')
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsn)
	return nil
}

// gets key from qs
func (app *application) readString(qs url.Values, key, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

// gets key from qs
func (app *application) readInt(
	qs url.Values, key string,
	defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	q, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return q
}
