// package xxhash implements a fast non-cryptographic hashing algorithm ported
// from https://code.google.com/p/xxhash/.
package xxhash

import (
	"encoding/binary"
	"hash"
	"unsafe"
)

const (
	prime1 = 2654435761
	prime2 = 2246822519
	prime3 = 3266489917
	prime4 = 668265263
	prime5 = 374761393
)

type digest struct {
	size int
	seed uint32
	v1   uint32
	v2   uint32
	v3   uint32
	v4   uint32
	over int
	tail []byte
}

// New returns a new 32-bit xxHash hash.Hash initialized with seed.
func New(seed uint32) hash.Hash32 {
	state := new(digest)
	state.seed = seed
	state.Reset()
	return state
}

func (d *digest) BlockSize() int {
	return 4
}

func (d *digest) Reset() {
	d.size = 0
	d.v1, d.v2, d.v3, d.v4 = initialize(d.seed)
	d.over = 0
	d.tail = make([]byte, 16)
}

func (d *digest) Size() int {
	return 4
}

func (d *digest) Sum(b []byte) []byte {
	h := d.Sum32()
	r := [4]byte{}
	binary.BigEndian.PutUint32(r[:], h)
	return append(b, r[:]...)
}

func (d *digest) Sum32() uint32 {
	var h uint32

	if d.size < 16 {
		h = d.seed + prime5
	} else {
		h = (d.v1 << 1) | (d.v1 >> 31)
		h += (d.v2 << 7) | (d.v2 >> 25)
		h += (d.v3 << 12) | (d.v3 >> 20)
		h += (d.v4 << 18) | (d.v4 >> 14)
	}

	h += uint32(d.size)
	h = tail(h, d.tail, d.over)
	return fmix(h)
}

func (d *digest) Write(data []byte) (int, error) {
	i := 0
	n := len(data)
	d.size += n

	// A blocksize of 16 is required in order to Write.
	if n+d.over < 16 {
		copy(d.tail[d.over:], data)
		d.over += n
		return n, nil
	}

	// Write full blocksize of previous and current data.
	copy(d.tail[d.over:], data)
	i = 16 - d.over

	j := uintptr(unsafe.Pointer(&d.tail[0]))
	k := j + 16
	d.v1, d.v2, d.v3, d.v4 = body(d.v1, d.v2, d.v3, d.v4, j, k)
	d.over = 0

	// Write all 16 byte blocks of data.
	j = uintptr(unsafe.Pointer(&data[i]))
	k = j + uintptr(n-i)
	d.v1, d.v2, d.v3, d.v4 = body(d.v1, d.v2, d.v3, d.v4, j, k)
	i = n - (n & 15)

	// Place any remaining data in memory.
	copy(d.tail, data[i:])
	d.over = n - i

	return n, nil
}

// Checksum returns the xxHash checksum of data.
func Checksum(data []byte, seed uint32) uint32 {
	n := len(data)

	v1, v2, v3, v4 := initialize(seed)

	if n > 0 {
		i := uintptr(unsafe.Pointer(&data[0]))
		j := i + uintptr(n)
		v1, v2, v3, v4 = body(v1, v2, v3, v4, i, j)
	}

	var h uint32
	if n < 16 {
		h = seed + prime5
	} else {
		h = (v1 << 1) | (v1 >> 31)
		h += (v2 << 7) | (v2 >> 25)
		h += (v3 << 12) | (v3 >> 20)
		h += (v4 << 18) | (v4 >> 14)
	}
	h += uint32(n)

	q := data[n-(n&15):]
	h = tail(h, q, len(q))
	return fmix(h)
}

func initialize(seed uint32) (v1, v2, v3, v4 uint32) {
	v1 = seed + prime1 + prime2
	v2 = seed + prime2
	v3 = seed
	v4 = seed - prime1
	return
}

func update(v, x uint32) uint32 {
	v += x * prime2
	return ((v << 13) | (v >> 19)) * prime1
}

func body(v1, v2, v3, v4 uint32, i, n uintptr) (uint32, uint32, uint32, uint32) {
	for i <= n-16 {
		v1, i = update(v1, *(*uint32)(unsafe.Pointer(i))), i+4
		v2, i = update(v2, *(*uint32)(unsafe.Pointer(i))), i+4
		v3, i = update(v3, *(*uint32)(unsafe.Pointer(i))), i+4
		v4, i = update(v4, *(*uint32)(unsafe.Pointer(i))), i+4
	}
	return v1, v2, v3, v4
}

func tail(h uint32, p []byte, n int) uint32 {
	i := 0

	for ; i <= n-4; i += 4 {
		h += *(*uint32)(unsafe.Pointer(&p[i])) * prime3
		h = ((h << 17) | (h >> 15)) * prime4
	}

	for ; i < n; i++ {
		h += uint32(p[i]) * prime5
		h = ((h << 11) | (h >> 21)) * prime1
	}

	return h
}

func fmix(h uint32) uint32 {
	h ^= h >> 15
	h *= prime2
	h ^= h >> 13
	h *= prime3
	h ^= h >> 16
	return h
}
