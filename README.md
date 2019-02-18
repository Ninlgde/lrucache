# Ninlgde LUR Cache
this is a very fast LRU cache implemented by Ninlgde in 2019

## Installation
To install LRU cache package, you need to install Go and set your Go workspace first.

1. Download and install it:

```sh
$ go get -u github.com/Ninlgde/lrucache
```

2. Import it in your code:

```go
import "github.com/Ninlgde/lrucache/go"
```

## Quick start

```go
package main

import (
	"fmt"
	"github.com/Ninlgde/lrucache/go"
)

func main() {
	// new a lru cache
	var cache = lru.NewLRUCache(10)

	// add key and value to cache
	cache.Add(10, 10)
	for i:=1; i<9; i++ {
		cache.Add(i, i)
	}

	// print all (k,v) in cache (head -> tail)
	iter := cache.Iterator(true)
	for v := range iter.C {
		fmt.Println(v)
	}

	fmt.Println()

	cache.Add("hhhh", 10)
	cache.Add(10, "hhhhh")
	fmt.Println(cache.Find(10)) // should be "hhhhh"
	fmt.Println(cache.Size())

	// print all (k,v) in cache (tail -> head)
	for v:= range cache.Iter(false) {
		fmt.Println(v)
	}
}
```

more example see the test go files

## Benchmark
Benchmark on MacBook Pro 2018

```
goos: darwin
goarch: amd64
pkg: github.com/Ninlgde/lrucache/go
BenchmarkThreadSafeLRU_Add-12           	 3000000	       333 ns/op
BenchmarkThreadSafeLRU_Size-12          	30000000	      41.8 ns/op
BenchmarkThreadSafeLRU_Find-12          	10000000	       224 ns/op
BenchmarkThreadSafeLRU_Find2-12         	10000000	       232 ns/op
BenchmarkThreadSafeLRU_Remove-12        	10000000	       227 ns/op
BenchmarkThreadSafeLRU_Iterator-12      	   10000	    144976 ns/op
BenchmarkThreadSafeLRU_Itera-12         	   10000	    121500 ns/op
BenchmarkThreadUnsafeLRU_Add-12         	 5000000	       232 ns/op
BenchmarkThreadUnsafeLRU_Size-12          2000000000	      1.49 ns/op
BenchmarkThreadUnsafeLRU_Find-12        	10000000	       167 ns/op
BenchmarkThreadUnsafeLRU_Find2-12       	10000000	       173 ns/op
BenchmarkThreadUnsafeLRU_Remove-12      	10000000	       176 ns/op
BenchmarkThreadUnsafeLRU_Iterator-12    	   10000	    135897 ns/op
BenchmarkThreadUnsafeLRU_Iter-12        	   10000	    131905 ns/op
```