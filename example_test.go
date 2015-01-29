// +build go1.1

package wsi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-on/wsi"
)

type Person struct {
	Id   int
	Name string
	Age  int `json:",omitempty"`
}

// maps the columns to the fields of a new Person
// must be a pointer method
func (p *Person) MapColumns(colToField map[string]interface{}) {
	colToField["id"] = &p.Id
	colToField["name"] = &p.Name
	colToField["age"] = &p.Age
}

// and example error handler
func printErr(r *http.Request, err error) {
	fmt.Printf("Error in route GET %s: %s\n", r.URL.Path, err.Error())
}

// we need to return an error here, even if we handle the response writing, so that the general
// error handler may be called
func createPerson(m map[string]interface{}, w http.ResponseWriter) error {
	// we fake a created response here
	res := map[string]interface{}{"Id": 400, "Name": m["name"]}
	w.WriteHeader(http.StatusCreated)
	wsi.ServeJSON(res, w)
	return nil
}

// findPersons defines the search sql.
// it must handle edge case, like limit = 0 or max limits, however limit and offset will never be < 0
func findPersons(opts wsi.QueryOptions) (wsi.Scanner, error) {
	// here we use a fake scanner to simulate database content
	return wsi.NewTestQuery(testData...), nil
	/*
		For real database queries you might want to do something like this:

		if len(opts.OrderBy) == 0 {
			opts.OrderBy = append(opts.OrderBy, "id asc")
		}

		// handle max limit
		limit := opts.Limit
		if limit == 0 || limit > 30 {
			limit = 30
		}

		return wsi.DBQuery(DB, "SELECT id,name from person ORDER BY $1 LIMIT $2 OFFSET $3", strings.Join(opts.OrderBy, ","), limit, opts.Offset)
	*/
}

var personRessource wsi.Ressource = func() wsi.ColumnsMapper { return &Person{} }

// create a http.Handler based on createPerson that load persons as json
var addHandler = personRessource.Exec(createPerson).SetErrorCallback(printErr)

// create a http.Handler based on findPersons that writes the resulting persons as json
var findHandler = personRessource.Query(findPersons).SetErrorCallback(printErr)

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
	// [{"Id":12,"Name":"Adrian"}
	// ,{"Id":24,"Name":"George"}
	// ]
	// -----
	// {"Id":400,"Name":"Peter"}
	//
}

var testData = []map[string]wsi.Setter{
	map[string]wsi.Setter{"id": wsi.SetInt(12), "name": wsi.SetString("Adrian")},
	map[string]wsi.Setter{"id": wsi.SetInt(24), "name": wsi.SetString("George")},
}
