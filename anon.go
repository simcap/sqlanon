package main

import (
	"math/rand"
	"strings"
	"unicode"
)

type anonymizer interface {
	anonymize(f sqlField) sqlField
}

type noopAnonymizer struct{}

func (*noopAnonymizer) anonymize(f sqlField) sqlField { return f }

type stringScrambler struct{}

func (*stringScrambler) anonymize(f sqlField) sqlField {
	switch v := f.value().(type) {
	case string:
		return &field{n: f.name(), v: scrambleString(v)}
	default:
		return f
	}
}

var (
	letters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	digits  = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
)

func scrambleString(value string) string {
	var out strings.Builder

	rand.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})
	rand.Shuffle(len(digits), func(i, j int) {
		digits[i], digits[j] = digits[j], digits[i]
	})

	for i, a := range value {
		if unicode.IsLetter(a) {
			out.WriteByte(letters[i%len(letters)])
		} else if unicode.IsDigit(a) {
			out.WriteByte(digits[i%len(digits)])
		} else {
			out.WriteRune(a)
		}
	}
	return out.String()
}
