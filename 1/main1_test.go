package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
)

var _ fmt.Stringer = inputs{}

func TestInputsString(t *testing.T) {
	tests := []struct {
		name string
		in   inputs
		want string
	}{
		{
			name: "значения из задания",
			in: inputs{
				numDecimal:     23,
				numOctal:       0o27,
				numHexadecimal: 0x17,
				numFloat:       1.0,
				name:           "task_1",
				isTaskOne:      true,
				numComplex:     1 + 2i,
			},
			want: "2323231.00task_1true(1.00+2.00i)",
		},
		{
			name: "нулевые значения",
			in:   inputs{},
			want: "0000.00false(0.00+0.00i)",
		},
		{
			name: "отрицательные значения",
			in: inputs{
				numDecimal:     -1,
				numOctal:       -8,
				numHexadecimal: -10,
				numFloat:       -1.5,
				name:           "x",
				isTaskOne:      false,
				numComplex:     -1 - 2i,
			},
			want: "-1-8-10-1.50xfalse(-1.00-2.00i)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.in.String(); got != tc.want {
				t.Errorf("inputs.String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestHashRunes(t *testing.T) {
	tests := []struct {
		name  string
		runes []rune
		salt  string
		want  string
	}{
		{
			name:  "чётная длина",
			runes: []rune("test"),
			salt:  "test salt",
			want:  "db5baf4487e1d7f6563e53830459a7feb64b3665820386072c5980d70edaa63e",
		},
		{
			name:  "чётная длина, пустая соль",
			runes: []rune("abcdef"),
			salt:  "",
			want:  "bef57ec7f53a6d40beb640a780a639c83bc29ac8a9816f1fc6c5c6dcd93c4721",
		},
		{
			name:  "нечётная длина",
			runes: []rune("abc"),
			salt:  salt,
			want:  "3a3a0ad5ed8bfa9531c50a0fcd1203fa0e4a42c0ade56848e69a0553f588aee8",
		},
		{
			name:  "две руны",
			runes: []rune("ab"),
			salt:  salt,
			want:  "660c6d86f7ab65536524d169964156f3cede5328a2086a0babcae3dd0e99fda0",
		},
		{
			name:  "одна руна",
			runes: []rune("a"),
			salt:  "",
			want:  "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
		},
		{
			name:  "пустой срез",
			runes: nil,
			salt:  salt,
			want:  "66802df107aace17871a5b610ff9eb11706e13477bb24e93966ca80671c0fac6",
		},
		{
			name:  "многобайтные руны",
			runes: []rune("日本語"),
			salt:  salt,
			want:  "020100bd1e05591fabf6f7615349e63b06576bc0622f6569f00b47f304e1ee0d",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := hashRunes(tc.runes, tc.salt); got != tc.want {
				t.Errorf("hashRunes(%q, %q) = %s, want %s", string(tc.runes), tc.salt, got, tc.want)
			}
		})
	}
}

func TestHashRunes_SaltIsInside(t *testing.T) {
	got := hashRunes([]rune("abcdef"), salt)

	sum := sha256.Sum256([]byte("abc" + salt + "def"))
	want := hex.EncodeToString(sum[:])

	if got != want {
		t.Errorf("соль не в середине:\ngot  %s\nwant %s", got, want)
	}
}

func TestHashRunes_SplitsByRunes(t *testing.T) {
	runes := []rune("日本語")

	got := hashRunes(runes, salt)

	sum := sha256.Sum256([]byte("日" + salt + "本語"))
	want := hex.EncodeToString(sum[:])

	if got != want {
		t.Fatalf("hashRunes(日本語) = %s, want %s", got, want)
	}

	b := []byte(string(runes))
	middle := len(b) / 2
	byteSum := sha256.Sum256([]byte(string(b[:middle]) + salt + string(b[middle:])))
	if got == hex.EncodeToString(byteSum[:]) {
		t.Error("соль вставлена по байтовой границе, а не по границе рун")
	}
}
