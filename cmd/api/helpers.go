package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
