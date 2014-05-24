# xxHash
Go package for a fast non-cryptographic hash function. Referenced from
[here](https://code.google.com/p/xxhash/). It's significantly faster than the
standard library's FNV hash, mostly because of the `unsafe` package. It
satisfies the hash.Hash32 interface but a quicker checksum function is provided
as well.

### Installation
~~~
go get github.com/asadmshah/xxhash
~~~

### Usage
Using it through the hash.Hash interface:
~~~
var seed uint32 = 1
h := xxhash.New(seed)
h.Write([]byte("a"))
h.Sum(nil) // Returns []byte{0x0b, 0x2c, 0xb7, 0x92}
~~~
Or to quickly hash something:
~~~
var seed uint32 = 1
xxhash.Checksum([]byte("a"), seed) // Returns []byte{0x0b, 0x2c, 0xb7, 0x92}
~~~

### Benchmarks
Benchmarks are run using 1024 bytes.
~~~
BenchmarkHashing	 5000000	       369 ns/op	2774.14 MB/s	      24 B/op	       2 allocs/op
BenchmarkChecksum	10000000	       183 ns/op	5586.22 MB/s	       0 B/op	       0 allocs/op
BenchmarkFNVHash	 1000000	      1353 ns/op	 756.68 MB/s	       8 B/op	       1 allocs/op
~~~
