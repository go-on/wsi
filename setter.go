package wsi

// Setter is a helper to be used inside of fake scanner functions for unit testing.
type Setter interface {
	// Set set the target to the objects value.
	// Target must be a pointer of the objects type, otherwise an error will be returned
	Set(target interface{}) error
}

func TestQuery(d ...map[string]Setter) func(targets map[string]interface{}) (stop bool, err error) {
	var counter int
	return func(targets map[string]interface{}) (stop bool, err error) {
		if len(d) == 0 {
			return true, nil
		}
		err = SetMap(d[counter], targets)
		if err != nil {
			return true, err
		}
		counter++
		return counter >= len(d), nil
	}
}

func NewTestQuery(d ...map[string]Setter) Scanner {
	return NewTestScanner(TestQuery(d...))
}

// SetMap sets each value in the target by the setter of the same key inside src.
// If is assumed that type inside the target is a pointer to the underlying type of the corresponding setter
// If this conditions are not met, an error will be returned
// SetMap is to be used inside a function for a fake scanner
func SetMap(src map[string]Setter, target map[string]interface{}) error {
	for col, setter := range src {
		t, ok := target[col]
		if !ok {
			continue
		}
		err := setter.Set(t)
		if err != nil {
			return &SetMapError{err.(*SetError), col}
		}
	}
	return nil
}

type SetInt int

// SetInt sets an int side the target. The target must be *int.
func (i SetInt) Set(target interface{}) error {
	ptr, ok := target.(*int)
	if !ok {
		return &SetError{target, "*int"}
	}
	*ptr = int(i)
	return nil
}

type SetFloat64 float64

// SetFloat64 sets a float side the target. The target must be *float64.
func (f SetFloat64) Set(target interface{}) error {
	ptr, ok := target.(*float64)
	if !ok {
		return &SetError{target, "*float64"}
	}
	*ptr = float64(f)
	return nil
}

type SetString string

// SetString sets a string inside the target. The target must be *string.
func (s SetString) Set(target interface{}) error {
	ptr, ok := target.(*string)
	if !ok {
		return &SetError{target, "*string"}
	}
	*ptr = string(s)
	return nil
}

type SetBool bool

// SetBool sets a bool inside the target. The target must be *bool.
func (b SetBool) Set(target interface{}) error {
	ptr, ok := target.(*bool)
	if !ok {
		return &SetError{target, "*bool"}
	}
	*ptr = bool(b)
	return nil
}

type SetError struct {
	Target interface{}
	MustBe string
}

func (s *SetError) Error() string {
	return "wrong type of target, must be '" + s.MustBe + "'"
}

type SetMapError struct {
	*SetError
	Column string
}

func (s *SetMapError) Error() string {
	return "wrong type of target for column '" + s.Column + "', must be '" + s.SetError.MustBe + "'"
}
