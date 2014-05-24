package xxhash

import (
	"encoding/binary"
	"reflect"
	"testing"
)

var golden32 = []struct {
	sum  []byte
	text []byte
}{
	{[]byte{0x0b, 0x2c, 0xb7, 0x92}, []byte("")},
	{[]byte{0xf5, 0x14, 0x70, 0x6f}, []byte("a")},
	{[]byte{0x2d, 0xee, 0x2a, 0x72}, []byte("ab")},
	{[]byte{0xaa, 0x3d, 0xa8, 0xff}, []byte("abc")},
	{[]byte{0xbe, 0x6e, 0x29, 0x9f}, []byte("abcd")},
	{[]byte{0xde, 0x13, 0xc4, 0x0a}, []byte("abcde")},
	{[]byte{0xda, 0x32, 0xca, 0xfb}, []byte("abcdef")},
	{[]byte{0xa7, 0x6d, 0x52, 0xdb}, []byte("abcdefg")},
	{[]byte{0xb6, 0x46, 0xe7, 0xe3}, []byte("abcdefgh")},
	{[]byte{0x1f, 0x73, 0x83, 0x97}, []byte("abcdefghi")},
	{[]byte{0x5b, 0xb3, 0xb6, 0xb8}, []byte("abcdefghij")},
	{[]byte{0xdb, 0x49, 0xea, 0x1c}, []byte("Discard medicine more than two years old.")},
	{[]byte{0xbb, 0x78, 0x79, 0x88}, []byte("He who has a shady past knows that nice guys finish last.")},
	{[]byte{0xc1, 0x59, 0x17, 0xea}, []byte("I wouldn't marry him with a ten foot pole.")},
	{[]byte{0xfe, 0x67, 0x63, 0xa3}, []byte("Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave")},
	{[]byte{0x1a, 0x88, 0xd2, 0x2b}, []byte("The days of the digital watch are numbered.  -Tom Stoppard")},
	{[]byte{0x4a, 0x9c, 0x59, 0xc8}, []byte("Nepal premier won't resign.")},
	{[]byte{0xb6, 0x69, 0xec, 0xe9}, []byte("For every action there is an equal and opposite government program.")},
	{[]byte{0x21, 0xa2, 0x38, 0x2b}, []byte("His money is twice tainted: 'taint yours and 'taint mine.")},
	{[]byte{0x7b, 0x35, 0xab, 0xee}, []byte("There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977")},
	{[]byte{0xc2, 0xce, 0x6d, 0xf9}, []byte("It's a tiny change to the code and not completely disgusting. - Bob Manchek")},
	{[]byte{0x3e, 0xa3, 0xac, 0x72}, []byte("size:  a.out:  bad magic")},
	{[]byte{0x79, 0xd0, 0x91, 0xd1}, []byte("The major problem is with sendmail.  -Mark Horton")},
	{[]byte{0x3e, 0xf0, 0xe8, 0xd8}, []byte("Give me a rock, paper and scissors and I will move the world.  CCFestoon")},
	{[]byte{0xc2, 0x59, 0x2a, 0x29}, []byte("If the enemy is within range, then so are you.")},
	{[]byte{0x02, 0xe2, 0xd4, 0xab}, []byte("It's well we cannot hear the screams/That we create in others' dreams.")},
	{[]byte{0x72, 0x0f, 0xf9, 0x43}, []byte("You remind me of a TV show, but that's all right: I watch it anyway.")},
	{[]byte{0xf8, 0x94, 0xd4, 0x3a}, []byte("C is as portable as Stonehedge!!")},
	{[]byte{0x60, 0xec, 0xec, 0x42}, []byte("Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley")},
	{[]byte{0xd4, 0x66, 0xb5, 0xb8}, []byte("The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule")},
	{[]byte{0x44, 0xd0, 0xd5, 0xbe}, []byte("How can you write a big system without C++?  -Paul Glick")},
}

func TestHash(t *testing.T) {
	h := New(1)
	f := func(b []byte) []byte {
		h.Reset()
		h.Write(b)
		return h.Sum(nil)
	}
	doTest(t, f)
}

func TestChecksum(t *testing.T) {
	f := func(b []byte) []byte {
		r := [4]byte{}
		binary.BigEndian.PutUint32(r[:], Checksum(b, 1))
		return r[:]
	}
	doTest(t, f)
}

func doTest(t *testing.T, f func([]byte) []byte) {
	for i, input := range golden32 {
		expect := input.sum
		result := f(input.text)
		if !reflect.DeepEqual(expect, result) {
			t.Fatalf("%d | Expected %x got %x", i, expect, result)
		}
	}
}

func BenchmarkHashing(b *testing.B) {
	h := New(1)
	f := func(data []byte) {
		h.Reset()
		h.Write(data)
		h.Sum(nil)
	}
	doBenchmark(b, f)
}

func BenchmarkChecksum(b *testing.B) {
	f := func(data []byte) {
		Checksum(data, 1)
	}
	doBenchmark(b, f)
}

func doBenchmark(b *testing.B, f func([]byte)) {
	b.SetBytes(1024)
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(data)
	}
}
