package waffleiron

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type FuncParser[T any] func(r *reader) (T, error)

func (p FuncParser[T]) Parse(r *reader) (T, error) {
	return p(r)
}

// Rune returns Parser that consumes a rune and return it if the current rune is same as rn
func Rune(rn rune) Parser[rune] {
	return runeParser{rn}
}

type runeParser struct {
	rn rune
}

// Parse implements Parser interface
func (p runeParser) Parse(r *reader) (rune, error) {
	ch, _, err := r.readRune()
	if err != nil {
		return 0, errors.Wrapf(err, "at %s", r.pos)
	}
	if ch != p.rn {
		return 0, errors.Errorf("expected %q, found %q at %s", p.rn, ch, r.pos)
	}
	return ch, nil
}

// String returns Parser that consumes a string and return it if remaining string starts with str
func String(str string) Parser[string] {
	return stringParser{str}
}

type stringParser struct {
	str string
}

// Parse implements Parser interface
func (p stringParser) Parse(r *reader) (string, error) {
	overrun := int64(len(p.str)) > int64(len(r.str))-r.idx
	if overrun || !strings.HasPrefix(r.remainingString(), p.str) {
		return "", errors.Errorf("expected %q, but not found at %s", p.str, r.pos)
	}
	s, err := r.consumeBytes(len(p.str))
	if err != nil || s != p.str {
		panic(fmt.Sprintf(
			"waffleiron.String(%q) consumed wrong bytes. it might be bug. please submit an issue.",
			p.str,
		))
	}
	return p.str, nil
}

// Regexp returns Parser that consume a string and return it if remaining string matches re
func Regexp(re *regexp.Regexp) Parser[string] {
	if !strings.HasPrefix(re.String(), "^") {
		return regexpParser{
			re: regexp.MustCompile("^" + re.String()),
		}
	}
	return regexpParser{re}
}

type regexpParser struct {
	re *regexp.Regexp
}

// Parse implements Parser interface
func (p regexpParser) Parse(r *reader) (string, error) {
	str := r.remainingString()
	loc := p.re.FindStringIndex(str)
	if len(loc) == 0 {
		return "", errors.Errorf("expected to match %q at %s", p.re, r.pos)
	}
	if loc[0] != 0 {
		panic("regex matched on loc[0] != 0. it might be bug. please submit an issue.")
	}
	s, err := r.consumeBytes(loc[1])
	if err != nil {
		panic(fmt.Sprintf(
			"waffleiron.Regexp(%s) consumed wrong bytes. it might be bug. please submit an issue.",
			p.re,
		))
	}
	return s, nil
}

// Regexp compiles str as a regexp and returns a Regexp parser
func RegexpStr(str string) Parser[string] {
	return Regexp(regexp.MustCompile(str))
}

var intParser Parser[int]

// Int returns a Parser that parses int
func Int() Parser[int] {
	if intParser == nil {
		intRegexp := regexp.MustCompile("^[+\\-]?[0-9]+")
		intParser = Map(
			Regexp(intRegexp),
			func(s string) int {
				i, err := strconv.Atoi(s)
				if err != nil {
					panic("failed to convert into int. this seems to be a bug. please submit an issue.")
				}
				return i
			},
		)
	}
	return intParser
}

// Pure returns a Parser that returns value without consuming a Reader
func Pure[T any](value T) Parser[T] {
	return pureParser[T]{value}
}

type pureParser[T any] struct {
	value T
}

// Parse implements Parser interface
func (p pureParser[T]) Parse(r *reader) (T, error) {
	return p.value, nil
}
