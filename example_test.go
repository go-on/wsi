// +build go1.1

package wsi_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-on/wsi"
)

type Person struct {
	ID   int
	Name string
	Age  int `json:",omitempty" sql:",omitempty"`
}

// newPerson is a function that creates a new person.
// we need it to generate the http.Handlers
func newPerson() interface{} { return &Person{} }

// findPersonsFake fakes our query, for a realistic query, see findPersons
// if any error happens, it must write to the response writer and return an error
func findPersonsFake(limit, offset int, w http.ResponseWriter, r *http.Request) (wsi.Scanner, error) {
	return wsi.NewTestQuery([]string{"ID", "Name"}, testData...), nil
}

// creates a http.Handler based on findPersonsFake that writes the resulting persons as json
// we are using the fake query here to avoid the need for a database, you may replace findPersonsFake
// with findPersons if you have a real database connection
var findHandler = wsi.Ressource{newPerson, printErr}.Query(findPersonsFake)

var DB *sql.DB

// findPersons defines the search sql.
// it must handle edge case, like limit = 0 or max limits, however limit and offset will never be < 0
func findPersons(limit, offset int, w http.ResponseWriter, r *http.Request) (wsi.Scanner, error) {
	// handle max limit
	if limit == 0 || limit > 30 {
		limit = 30
	}

	return wsi.DBQuery(
		DB,
		`SELECT "ID","Name" from person ORDER BY "ID" LIMIT $1 OFFSET $2`,
		limit,
		offset,
	)
}

// createPerson creates a person based on the values of the given ColumnsMapper
// and writes to the given responsewriter
// we need to return an error here, even if we handle the response writing, so that the general
// error handler may be called
func createPerson(m map[string]interface{}, w http.ResponseWriter, r *http.Request) error {
	// we fake a created response here
	m["ID"] = 400
	// res := map[string]interface{}{"ID": 400, "Name": m.Map("name")}
	w.WriteHeader(http.StatusCreated)
	wsi.ServeJSON(m, w)
	return nil
}

// creates a http.Handler based on createPerson that load persons as json
var addHandler = wsi.Ressource{newPerson, printErr}.Exec(createPerson)

func Example() {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/person/", nil)
	findHandler.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	fmt.Println("-----")

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/person", strings.NewReader(`{"Name":"Peter"}`))
	addHandler.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())

	// Output:
	// [{"ID":12,"Name":"Adrian"}
	// ,{"ID":24,"Name":"George"}
	// ]
	// -----
	// {"ID":400,"Name":"Peter"}
	//
}

var testData = []map[string]wsi.Setter{
	map[string]wsi.Setter{"ID": wsi.SetInt(12), "Name": wsi.SetString("Adrian")},
	map[string]wsi.Setter{"ID": wsi.SetInt(24), "Name": wsi.SetString("George")},
}

// an example error handler
func printErr(r *http.Request, err error) {
	fmt.Printf("Error in route GET %s: %T %s\n", r.URL.Path, err, err.Error())
}
