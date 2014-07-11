package hikaru

import (
	"fmt"
	"net/url"
	"strconv"
)

type Value struct {
	v interface{}
}

func (v Value) Empty() bool {
	return v.v == nil
}

func (v Value) String() (string, error) {
	switch v2 := v.v.(type) {
	case string:
		return v2, nil
	}
	return "", fmt.Errorf("invalid type %T", v)
}

func (v Value) Int() (int64, error) {
	switch v2 := v.v.(type) {
	case int:
	case int8:
	case int16:
	case int32:
		return int64(v2), nil
	case int64:
		return v2, nil
	case string:
		return strconv.ParseInt(v2, 10, 64)
	}
	return 0, fmt.Errorf("invalid type %T", v)
}

func (v Value) Float() (float64, error) {
	switch v2 := v.v.(type) {
	case float32:
		return float64(v2), nil
	case float64:
		return v2, nil
	case string:
		return strconv.ParseFloat(v2, 64)
	}
	return 0, fmt.Errorf("invalid type %T", v)
}

func (v Value) Bool() (bool, error) {
	switch v2 := v.v.(type) {
	case bool:
		return v2, nil
	case string:
		return strconv.ParseBool(v2)
	}
	return false, fmt.Errorf("invalid type %T", v)
}

type Values url.Values

// Returns whether the request has the given key
// in route values and query.
func (v Values) Has(key string) bool {
	_, ok := v[key]
	return ok
}

// Returns the first value associated with the given key
// from route values and query.
// If there are no values associated with the key, returns "".
// To access multiple values of a key, use Vals.
func (v Values) Value(key string) Value {
	r := Value{}
	ss, ok := v[key]
	if ok && ss != nil && len(ss) > 0 {
		r.v = ss[0]
	}
	return r
}

func (v Values) String(key string) (string, error) {
	ss, ok := v[key]
	if ok && ss != nil && len(ss) > 0 {
		return ss[0], nil
	}
	return "", fmt.Errorf("%s not found", key)
}

func (v Values) StringE(key string, failover string) string {
	s, err := v.String(key)
	if err != nil {
		return failover
	}
	return s
}

func (v Values) Int(key string) (int64, error) {
	val := v.Value(key)
	if val.Empty() {
		return 0, fmt.Errorf("%s not found", key)
	}
	return val.Int()
}

func (v Values) IntE(key string, failover int64) int64 {
	s, err := v.Int(key)
	if err != nil {
		return failover
	}
	return s
}

func (v Values) Float(key string) (float64, error) {
	val := v.Value(key)
	if val.Empty() {
		return 0, fmt.Errorf("%s not found", key)
	}
	return val.Float()
}

func (v Values) FloatE(key string, failover float64) float64 {
	s, err := v.Float(key)
	if err != nil {
		return failover
	}
	return s
}

func (v Values) Bool(key string) (bool, error) {
	val := v.Value(key)
	if val.Empty() {
		return false, fmt.Errorf("%s not found", key)
	}
	return val.Bool()
}

func (v Values) BoolE(key string, failover bool) bool {
	s, err := v.Bool(key)
	if err != nil {
		return failover
	}
	return s
}

// Returns the list of values associated with the given key
// from route values and query.
// If there are no values associated with the key, returns empty slice.
func (v Values) Values(key string) []Value {
	ss, ok := v[key]
	if !ok || ss == nil || len(ss) <= 0 {
		return nil
	}
	vs := make([]Value, len(ss))
	for i := range vs {
		vs[i] = Value{ss[i]}
	}
	return vs
}

func (v Values) Strings(key string) []string {
	ss, ok := v[key]
	if !ok || ss == nil || len(ss) <= 0 {
		return nil
	}
	return ss
}

func (v Values) Ints(key string) ([]int64, error) {
	ss, ok := v[key]
	if !ok || ss == nil || len(ss) <= 0 {
		return nil, fmt.Errorf("%s not found", key)
	}
	vs := make([]int64, len(ss))
	for i := range vs {
		val := Value{ss[i]}
		conv, err := val.Int()
		if err != nil {
			return nil, err
		}
		vs[i] = conv
	}
	return vs, nil
}

// Set sets the key to value. It replaces any existing
// values.
func (v Values) Set(key, value string) {
	v[key] = []string{value}
}

// Add adds the key to value. It appends to any existing
// values associated with key.
func (v Values) Add(key, value string) {
	v[key] = append(v[key], value)
}

// Del deletes the values associated with key.
func (v Values) Del(key string) {
	delete(v, key)
}

func (v Values) Update(v2 Values) {
	for key, ss := range v2 {
		if ss != nil && len(ss) > 0 {
			for _, s := range ss {
				v.Add(key, s)
			}
		}
	}
}
