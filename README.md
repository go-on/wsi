# wsi
integrates web requests with database/sql (golang)

[![GoDoc](https://godoc.org/github.com/go-on/wsi?status.png)](http://godoc.org/github.com/go-on/wsi)

## Example

```go

import (
    "github.com/go-on/wsi"
    "net/http"
)

type Person struct {
    Id   int
    Name string
    Age  int `json:",omitempty"`
}

// maps the columns to the fields of a new Person; must be a pointer method
func (p *Person) MapColumns(colToField map[string]interface{}) {
    colToField["id"] = &p.Id
    colToField["name"] = &p.Name
    colToField["age"] = &p.Age
}

var newPerson wsi.Ressource = func() wsi.ColumnsMapper { return &Person{} }

func findPersons(opts wsi.QueryOptions) (wsi.Scanner, error) {
    return wsi.DBQuery(
        DB, 
        "SELECT id,name from person ORDER BY $1 LIMIT $2 OFFSET $3", 
        strings.Join(opts.OrderBy, ","), 
        opts.Limit, 
        opts.Offset,
    )
}


func main() {
    // you might set an error handler as well, see the api docs
    servePersons := newPerson.Query(findPersons)

    http.Handle("/person/", servePersons)
    http.ListenAndServe(":8080",nil)    
}

```