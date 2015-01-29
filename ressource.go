package wsi

import (
	"net/http"
)

// QueryFunc makes the sql query and returns a Scanner. If an error is returned, QueryFunc must write
// to the reponsewriter (set the status code etc). If no error is returned QueryFunc must not write
// to the response write. specific headers are the exception and may be set.
type QueryFunc func(QueryOptions, http.ResponseWriter, *http.Request) (Scanner, error)

// ExecFunc makes the sql exec and writes to the response writer. If must return an error, if
// some happened, so that the error may be passed to the general error handler
type ExecFunc func(map[string]interface{}, http.ResponseWriter, *http.Request) error

type Ressource func() ColumnsMapper

func (r Ressource) Exec(e ExecFunc) Exec {
	return Exec{mapperFn: r, fn: e, dec: JSONDecoder}
}

func (r Ressource) Query(q QueryFunc) Query {
	return Query{encFn: NewJSONStreamer, mapperFn: r, fn: q}
}
