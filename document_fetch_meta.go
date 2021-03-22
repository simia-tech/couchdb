package couchdb

// // DocumentFetchMetaRequest defines the document fetch request.
// type DocumentFetchMetaRequest struct {
// 	document *DocumentRef
// 	id       string

// 	ctx      context.Context
// 	revision string
// }

// // WithContext adds a context to the request.
// func (dfr *DocumentFetchMetaRequest) WithContext(ctx context.Context) *DocumentFetchMetaRequest {
// 	dfr.ctx = ctx
// 	return dfr
// }

// // WithRevision adds a revision to the request.
// func (dfr *DocumentFetchMetaRequest) WithRevision(revision string) *DocumentFetchMetaRequest {
// 	dfr.revision = revision
// 	return dfr
// }

// // Do performs the request.
// func (dfr *DocumentFetchMetaRequest) Do() (*DocumentFetchMetaResponse, error) {
// 	request, err := dfr.document.database.client.requestFor(dfr.ctx, nil, http.MethodGet, dfr.document.database.name, dfr.id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	request.Header.Set("Content-Type", MimeTypeJSON)
// 	request.Header.Set("Accept", MimeTypeJSON)
// 	if dfr.revision != "" {
// 		request.Header.Set("If-None-Match", dfr.revision)
// 	}

// 	response, err := dfr.document.database.client.do(request)
// 	if err != nil {
// 		return nil, errx.Annotatef(err, "request [%s] [%s]", request.Method, request.URL)
// 	}
// 	defer response.Body.Close()

// 	if err = evaluateResponseStatus(response.StatusCode, "",
// 		http.StatusOK,
// 		http.StatusNotModified,
// 		http.StatusUnauthorized, http.StatusNotFound); err != nil {
// 		return nil, err
// 	}

// 	r := &DocumentFetchMetaResponse{}
// 	r.Revision = strings.Trim(response.Header.Get("ETag"), `"`)
// 	r.ContentLength, err = strconv.ParseUint(response.Header.Get("Content-Length"), 10, 64)
// 	if err != nil {
// 		return nil, errx.Annotatef(err, "parse uint64 [%s]", response.Header.Get("Content-Length"))
// 	}
// 	return r, nil
// }

// // DocumentFetchMetaResponse defines the document create response.
// type DocumentFetchMetaResponse struct {
// 	ContentLength uint64
// 	Revision      string
// }
