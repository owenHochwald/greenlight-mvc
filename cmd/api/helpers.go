package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"owenHochwald.greenlight/internal/validator"
)

func (app *application) readJSON(c *gin.Context, dst interface{}) error {
	// max allowable byte size of input
	maxBytes := 1_048_576
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(maxBytes))

	buf, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	// Refill body for Gin's binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))

	// Strict JSON decoding
	dec := json.NewDecoder(bytes.NewBuffer(buf))
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	return c.ShouldBindJSON(dst)
}

func (app *application) parseID(c *gin.Context) (int64, error) {
	movieId := c.Param("parseID")
	id, err := strconv.ParseInt(movieId, 10, 64)

	if err != nil {
		return 0, nil
	}

	return id, nil
}

func (app *application) readCSV(qs url.Values, key string, defaultValues []string) []string {
	s := qs.Get(key)

	if s == "" {
		return defaultValues
	}

	return strings.Split(s, ",")
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	i := qs.Get(key)

	if i == "" {
		return defaultValue
	}

	parsedInt, err := strconv.Atoi(i)

	if err != nil {
		v.AddError(key, "must be a valid integer value")
		return defaultValue
	}

	return parsedInt
}
