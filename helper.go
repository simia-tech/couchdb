package couchdb

import (
	"log"
	"net/http"

	"github.com/simia-tech/errx"
)

// Some MIME types.
const (
	MimeTypeTextPlain        = "text/plain"
	MimeTypeJSON             = "application/json"
	MimeTypeMultipartRelated = "multipart/related"
	MimeTypeOctetStream      = "application/octet-stream"
)

func evaluateResponseStatus(statusCode int, reason string, validStatusCodes ...int) error {
	valid := false
	for _, validStatusCode := range validStatusCodes {
		valid = valid || (validStatusCode == statusCode)
	}
	if !valid {
		return errx.Errorf("expected status codes %v, got %3d", validStatusCodes, statusCode)
	}

	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusCreated:
		return nil
	case http.StatusAccepted:
		return nil
	case http.StatusBadRequest:
		return errx.BadRequestf(stringOr(reason, "bad request"))
	case http.StatusUnauthorized:
		return errx.Unauthorizedf(stringOr(reason, "unauthorized"))
	case http.StatusNotFound:
		return errx.NotFoundf(stringOr(reason, "not found"))
	case http.StatusPreconditionFailed:
		return errx.AlreadyExistsf(stringOr(reason, "precondition failed"))
	case http.StatusConflict:
		return errx.AlreadyExistsf(stringOr(reason, "conflict"))
	default:
		log.Printf("status code %d", statusCode)
	}
	return nil
}

func stringOr(text, defaultText string) string {
	if text == "" {
		return defaultText
	}
	return text
}
