package couchdb

// // DocumentUpdateRequest defines the document create request.
// type DocumentUpdateRequest struct {
// 	document *DocumentRef
// 	id       string
// 	doc      Document

// 	ctx      context.Context
// 	revision string
// }

// // WithContext adds a context to the request.
// func (dur *DocumentUpdateRequest) WithContext(ctx context.Context) *DocumentUpdateRequest {
// 	dur.ctx = ctx
// 	return dur
// }

// // WithRevision adds a revision to the request.
// func (dur *DocumentUpdateRequest) WithRevision(revision string) *DocumentUpdateRequest {
// 	dur.revision = revision
// 	return dur
// }

// // Do performs the request.
// func (dur *DocumentUpdateRequest) Do() (*DocumentUpdateResponse, error) {
// 	buffer := bytes.Buffer{}
// 	if err := json.NewEncoder(&buffer).Encode(dur.doc); err != nil {
// 		return nil, errx.Annotatef(err, "json encode")
// 	}

// 	request, err := dur.document.database.client.requestFor(dur.ctx, &buffer, http.MethodPut, dur.document.database.name, dur.id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	request.Header.Set("Content-Type", MimeTypeJSON)
// 	request.Header.Set("Accept", MimeTypeJSON)
// 	if dur.revision != "" {
// 		request.Header.Set("If-Match", dur.revision)
// 	}

// 	response, err := dur.document.database.client.do(request)
// 	if err != nil {
// 		return nil, errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
// 	}
// 	defer response.Body.Close()

// 	r := &DocumentUpdateResponse{}
// 	if err := json.NewDecoder(response.Body).Decode(r); err != nil {
// 		return nil, errx.Annotatef(err, "json decode")
// 	}
// 	if err := evaluateResponseStatus(response.StatusCode, "",
// 		http.StatusCreated, http.StatusAccepted,
// 		http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound, http.StatusConflict); err != nil {
// 		return r, err
// 	}
// 	return r, nil
// }

// // DocumentUpdateResponse defines the document create response.
// type DocumentUpdateResponse struct {
// 	OK       bool   `json:"ok"`
// 	ID       string `json:"id"`
// 	Revision string `json:"rev"`
// }
