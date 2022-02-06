package waffleiron_test

import (
	"regexp"
	"testing"

	wi "github.com/seiyab/waffleiron"
)

func TestRegexp(t *testing.T) {
	t.Run("[a-z,A-Z]+", func(t *testing.T) {
		for _, ptn := range []string{"[a-z,A-Z]+", "^[a-z,A-Z]+"} {
			t.Run(ptn, func(t *testing.T) {
				p := wi.Regexp(regexp.MustCompile(ptn))

				t.Run("ok", func(t *testing.T) {
					testCases := []string{
						"abc",
						"XYZ",
						"T",
						"WaffleIron",
					}
					for _, tt := range testCases {
						t.Run(tt, func(t *testing.T) {
							v, err := wi.Parse(tt, p)
							if err != nil {
								t.Fatal(err)
							}
							if v != tt {
								t.Errorf("v = %q, want %q", v, tt)
							}
						})
					}
				})

				t.Run("error", func(t *testing.T) {
					testCases := []string{
						"",
						"1abc",
						"abc1",
						"waffle iron",
					}
					for _, tt := range testCases {
						t.Run(tt, func(t *testing.T) {
							_, err := wi.Parse(tt, p)
							if err == nil {
								t.Error("expected error")
							}
						})
					}
				})
			})
		}
	})
}

func TestRegexpStr(t *testing.T) {
	t.Run("[0-9]+", func(t *testing.T) {
		p := wi.RegexpStr("[0-9]+")

		t.Run("ok", func(t *testing.T) {
			testCases := []string{
				"0",
				"123",
			}
			for _, tt := range testCases {
				t.Run(tt, func(t *testing.T) {
					v, err := wi.Parse(tt, p)
					if err != nil {
						t.Fatal(err)
					}
					if v != tt {
						t.Errorf("v = %q, want %q", v, tt)
					}
				})
			}
		})

		t.Run("error", func(t *testing.T) {
			testCases := []string{
				"",
				"123a",
				"x456",
				"waffle iron",
				"123,456",
			}
			for _, tt := range testCases {
				t.Run(tt, func(t *testing.T) {
					_, err := wi.Parse(tt, p)
					if err == nil {
						t.Error("expected error")
					}
				})
			}
		})
	})
}
