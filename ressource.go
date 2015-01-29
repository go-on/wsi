package wsi

import (
	"net/http"
)

type QueryFunc func(QueryOptions) (Scanner, error)
type ExecFunc func(map[string]interface{}, http.ResponseWriter) error

type Ressource func() ColumnsMapper

func (r Ressource) Exec(e ExecFunc) Exec {
	return Exec{mapperFn: r, fn: e, dec: JSONDecoder}
}

func (r Ressource) Query(q QueryFunc) Query {
	return Query{encFn: NewJSONStreamer, mapperFn: r, fn: q}
}
