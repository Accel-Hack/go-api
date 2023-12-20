package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ErrKeyNotFound struct {
	Key string
}

func (e ErrKeyNotFound) Error() string {
	return fmt.Sprintf("%q not found", e.Key)
}

type Parser[T any] func(url.Values) (T, error)

type Parse[T any] func(string, url.Values) (T, error)

func (p Parse[T]) Key(key string) Parser[T] {
	return func(vs url.Values) (T, error) {
		return p(key, vs)
	}
}

type ParseOrNil[T any] func(string, url.Values) (*T, error)

func (p ParseOrNil[T]) Key(key string) Parser[*T] {
	return func(vs url.Values) (*T, error) {
		return p(key, vs)
	}
}

func (p Parse[T]) Required() Parse[T] {
	return func(key string, vs url.Values) (T, error) {
		var z T
		if !vs.Has(key) {
			return z, ErrKeyNotFound{key}
		}
		return p(key, vs)
	}
}

func (p Parse[T]) OrNil() ParseOrNil[T] {
	return func(key string, vs url.Values) (*T, error) {
		if !vs.Has(key) {
			return nil, nil
		}
		v, err := p(key, vs)
		if err != nil {
			return nil, err
		}
		return &v, nil
	}
}

func QueryString() Parse[string] {
	return func(key string, vs url.Values) (string, error) {
		return vs.Get(key), nil
	}
}

func QueryBool() Parse[bool] {
	return func(key string, vs url.Values) (bool, error) {
		s := vs.Get(key)
		return strconv.ParseBool(s)
	}
}

func QueryInt() Parse[int] {
	return func(key string, vs url.Values) (int, error) {
		s := vs.Get(key)
		i64, err := strconv.ParseInt(s, 10, 0)
		return int(i64), err
	}
}

func QueryUUID() Parse[uuid.UUID] {
	return func(key string, vs url.Values) (uuid.UUID, error) {
		s := vs.Get(key)
		return uuid.Parse(s)
	}
}

func QueryTime() Parse[time.Time] {
	return func(key string, vs url.Values) (time.Time, error) {
		s := vs.Get(key)
		return time.ParseInLocation(time.DateOnly, s, time.Local)
	}
}
