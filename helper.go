package couchdb

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	"code.posteo.de/common/errx"
)

// Some MIME types.
const (
	MimeTypeTextPlain        = "text/plain"
	MimeTypeJSON             = "application/json"
	MimeTypeMultipartRelated = "multipart/related"
	MimeTypeOctetStream      = "application/octet-stream"
)

func evaluateResponse(response *http.Response) (interface{}, error) {
	body, err := evaluateResponseBody(textproto.MIMEHeader(response.Header), response.Body)
	if err != nil {
		return nil, errx.Annotatef(err, "evaluate body")
	}

	if err := evaluateResponseStatus(response.StatusCode, body); err != nil {
		return nil, err
	}

	return body, nil
}

func evaluateResponseStatus(statusCode int, body interface{}) error {
	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusCreated:
		return nil
	case http.StatusBadRequest:
		return errx.BadRequestf(getStringFieldOr(body, "reason", "bad request"))
	case http.StatusUnauthorized:
		return errx.Unauthorizedf(getStringFieldOr(body, "reason", "unauthorized"))
	case http.StatusNotFound:
		return errx.NotFoundf(getStringFieldOr(body, "reason", "not found"))
	case http.StatusPreconditionFailed:
		return errx.AlreadyExistsf(getStringFieldOr(body, "reason", "already exists"))
	case http.StatusConflict:
		return errx.AlreadyExistsf(getStringFieldOr(body, "reason", "already exists"))
	default:
		log.Printf("status code %d", statusCode)
	}
	return nil
}

func evaluateResponseBody(header textproto.MIMEHeader, body io.Reader) (interface{}, error) {
	contentType := header.Get("Content-Type")
	mediaType, parameters, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, errx.Annotatef(err, "parse media type [%s]", contentType)
	}

	switch mt := strings.ToLower(mediaType); mt {
	case MimeTypeJSON:
		bodyMap := map[string]interface{}{}
		if err := json.NewDecoder(body).Decode(&bodyMap); err != nil {
			return nil, errx.Annotatef(err, "json decode")
		}
		return bodyMap, nil
	case MimeTypeMultipartRelated:
		bodyMap := map[string]interface{}{}
		reader := multipart.NewReader(body, parameters["boundary"])
		for part, err := reader.NextPart(); err == nil; part, err = reader.NextPart() {
			if err != nil {
				return nil, err
			}
			name := part.FileName()
			if name == "" {
				name = "document"
			}
			value, err := evaluateResponseBody(part.Header, part)
			if err != nil {
				return bodyMap, errx.Annotatef(err, "multipart [%s]", name)
			}
			bodyMap[name] = value
		}
		return bodyMap, nil
	case MimeTypeOctetStream:
		value, err := ioutil.ReadAll(body)
		if err != nil {
			return nil, errx.Annotatef(err, "read all")
		}
		return value, nil
	case MimeTypeTextPlain:
		value, err := ioutil.ReadAll(body)
		if err != nil {
			return nil, errx.Annotatef(err, "read all")
		}
		return string(value), nil
	default:
		return nil, errx.NotImplementedf("could not handle response body type [%s]", mt)
	}
}

func getStringFieldOr(i interface{}, key, defaultValue string) string {
	if m, ok := i.(map[string]interface{}); ok {
		if reason, ok := m[key].(string); ok {
			return reason
		}
	}
	return defaultValue
}
