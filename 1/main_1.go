package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const salt = "go-2024"

type inputs struct {
	numDecimal     int
	numOctal       int
	numHexadecimal int
	numFloat       float64
	name           string
	isTaskOne      bool
	numComplex     complex64
}

func (in inputs) String() string {
	var b strings.Builder

	b.WriteString(strconv.Itoa(in.numDecimal))
	b.WriteString(strconv.Itoa(in.numOctal))
	b.WriteString(strconv.Itoa(in.numHexadecimal))
	b.WriteString(strconv.FormatFloat(in.numFloat, 'f', 2, 64))
	b.WriteString(in.name)
	b.WriteString(strconv.FormatBool(in.isTaskOne))
	b.WriteString(strconv.FormatComplex(complex128(in.numComplex), 'f', 2, 64))

	return b.String()
}

func hashRunes(runes []rune, salt string) string {
	middle := len(runes) / 2

	var b strings.Builder
	b.Grow(len(string(runes)) + len(salt))
	for _, r := range runes[:middle] {
		b.WriteRune(r)
	}
	b.WriteString(salt)
	for _, r := range runes[middle:] {
		b.WriteRune(r)
	}

	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

func main() {
	in := inputs{
		numDecimal:     23,   // десятичная система
		numOctal:       0o27, // восьмеричная система, 23
		numHexadecimal: 0x17, // шестнадцатеричная система, 23
		numFloat:       1.0,
		name:           "task_1",
		isTaskOne:      true,
		numComplex:     1 + 2i,
	}

	fmt.Printf("numDecimal (%T) %v\n", in.numDecimal, in.numDecimal)
	fmt.Printf("numOctal (%T) %v\n", in.numOctal, in.numOctal)
	fmt.Printf("numHexadecimal (%T) %v\n", in.numHexadecimal, in.numHexadecimal)
	fmt.Printf("numFloat (%T) %v\n", in.numFloat, in.numFloat)
	fmt.Printf("name (%T) %v\n", in.name, in.name)
	fmt.Printf("isTaskOne (%T) %v\n", in.isTaskOne, in.isTaskOne)
	fmt.Printf("numComplex (%T) %v\n", in.numComplex, in.numComplex)

	combined := in.String()

	runes := []rune(combined)

	fmt.Printf("combined (%d runes) = %q\n", len(runes), combined)
	fmt.Printf("sha256 = %s\n", hashRunes(runes, salt))
}
