package waffleiron

import (
	"github.com/pkg/errors"

	"github.com/hashicorp/go-multierror"
)

func And[T, U any](p0 Parser[T], p1 Parser[U]) Parser[Tuple2[T, U]] {
	return andParser[T, U]{p0, p1}
}

type andParser[T, U any] struct {
	p0 Parser[T]
	p1 Parser[U]
}

func (p andParser[T, U]) Parse(r *Reader) (Tuple2[T, U], error) {
	a, err := p.p0.Parse(r)
	if err != nil {
		return Tuple2[T, U]{}, err
	}
	b, err := p.p1.Parse(r)
	if err != nil {
		return Tuple2[T, U]{}, err
	}
	return NewTuple2(a, b), nil
}

func And3[T, U, V any](p0 Parser[T], p1 Parser[U], p2 Parser[V]) Parser[Tuple3[T, U, V]] {
	return and3Parser[T, U, V]{p0, p1, p2}
}

type and3Parser[T, U, V any] struct {
	p0 Parser[T]
	p1 Parser[U]
	p2 Parser[V]
}

func (p and3Parser[T, U, V]) Parse(r *Reader) (Tuple3[T, U, V], error) {
	v0, err := p.p0.Parse(r)
	if err != nil {
		return Tuple3[T, U, V]{}, err
	}
	v1, err := p.p1.Parse(r)
	if err != nil {
		return Tuple3[T, U, V]{}, err
	}
	v2, err := p.p2.Parse(r)
	if err != nil {
		return Tuple3[T, U, V]{}, err
	}
	return NewTuple3(v0, v1, v2), nil
}

func Or[T any](p0 Parser[T], ps ...Parser[T]) Parser[T] {
	parsers := make([]Parser[T], len(ps)+1)
	parsers[0] = p0
	for i, p := range ps {
		parsers[i+1] = p
	}
	return orParser[T]{ps: parsers}
}

type orParser[T any] struct {
	ps []Parser[T]
}

func (p orParser[T]) Parse(r *Reader) (T, error) {
	var totalErr error
	for _, p := range p.ps {
		var t T
		err := r.Try(func() error {
			var e error
			t, e = p.Parse(r)
			return e
		})
		if err == nil {
			return t, nil
		}

		totalErr = multierror.Append(totalErr, err)
	}
	return *new(T), totalErr
}

func Trace[T any](name string, p Parser[T]) Parser[T] {
	return traceParser[T]{name, p}
}

type traceParser[T any] struct {
	name string
	p    Parser[T]
}

func (p traceParser[T]) Parse(r *Reader) (T, error) {
	var t T
	var err error
	r.WithTrace(p.name, func() {
		t, err = p.p.Parse(r)
	})
	if err != nil {
		return t, errors.Wrapf(err, "at %q", p.name)
	}
	return t, nil
}
