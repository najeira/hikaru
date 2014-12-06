package hikaru

import (
	"net/url"
	"strconv"
	"errors"
)

var (
	ErrKeyNotExist = errors.New("not exist")
)

type Values struct {
	u *url.URL
	v url.Values
}

func NewValues(u *url.URL) *Values {
	return &Values{u: u}
}

func (v *Values) Values() url.Values {
	if v.v == nil {
		v.v = v.u.Query()
	}
	return v.v
}

// Has returns whether the request has the given key in the route values and 
// the query.
func (v Values) Has(key string) bool {
	_, ok := v.Values()[key]
	return ok
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key,
// Get returns the failover string.
// To access multiple values, use the map directly.
func (v Values) String(key string, failover string) string {
	ret, err := v.TryString(key)
	if err != nil {
		return failover
	}
	return ret
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key,
// Get returns the ErrKeyNotExist.
func (v Values) TryString(key string) (string, error) {
	ss, ok := v.Values()[key]
	if ok && ss != nil && len(ss) > 0 {
		return ss[0]
	}
	return "", ErrKeyNotExist
}

func (v Values) Int(key string, failover int64) int64 {
	ret, err := v.TryInt(key)
	if err != nil {
		return failover
	}
	return ret
}

func (v Values) TryInt(key string) (int64, error) {
	s, err := v.TryString(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(s, 10, 64)
}

func (v Values) Float(key string, failover float64) float64 {
	ret, err := v.TryFloat(key)
	if err != nil {
		return failover
	}
	return ret
}

func (v Values) TryFloat(key string) (float64, error) {
	s, err := v.TryString(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(s, 64)
}

func (v Values) Bool(key string, failover bool) bool {
	ret, err := v.TryBool(key)
	if err != nil {
		return failover
	}
	return ret
}

func (v Values) TryBool(key string) (bool, error) {
	s, err := v.TryString(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseBool(s)
}

// Set sets the key to value. It replaces any existing values.
func (v Values) Set(key, value string) {
	v[key] = []string{value}
}

// Add adds the key to value. It appends to any existing values associated 
// with key.
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
