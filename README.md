# Ninlgde LUR Cache
This is a very fast LRU cache implemented by Ninlgde in 2019

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

More examples see the test go files

## Benchmark
Benchmark on MacBook Pro 2018

```
goos: darwin
goarch: amd64
pkg: github.com/Ninlgde/lrucache/go
BenchmarkThreadSafeLRU_Add-12            	 1000000	      1289 ns/op
BenchmarkThreadSafeLRU_Add2-12           	 1000000	      1309 ns/op
BenchmarkThreadSafeLRU_Add3-12           	 3000000	       551 ns/op
BenchmarkThreadSafeLRU_Size-12           	30000000	        42.0 ns/op
BenchmarkThreadSafeLRU_Size2-12          	30000000	        45.0 ns/op
BenchmarkThreadSafeLRU_Find-12           	10000000	       227 ns/op
BenchmarkThreadSafeLRU_Find2-12          	10000000	       240 ns/op
BenchmarkThreadSafeLRU_Find3-12          	20000000	        82.1 ns/op
BenchmarkThreadSafeLRU_Remove-12         	10000000	       219 ns/op
BenchmarkThreadSafeLRU_Remove2-12        	10000000	       230 ns/op
BenchmarkThreadSafeLRU_Remove3-12        	20000000	        79.0 ns/op
BenchmarkThreadSafeLRU_Iterator-12       	   10000	    140919 ns/op
BenchmarkThreadSafeLRU_Iterator2-12      	 1000000	      2189 ns/op
BenchmarkThreadSafeLRU_Itera-12          	   10000	    133793 ns/op
BenchmarkThreadSafeLRU_Itera2-12         	 1000000	      1734 ns/op
BenchmarkThreadUnsafeLRU_Add-12          	 2000000	       965 ns/op
BenchmarkThreadUnsafeLRU_Add2-12         	 1000000	      1013 ns/op
BenchmarkThreadUnsafeLRU_Add3-12         	 3000000	       443 ns/op
BenchmarkThreadUnsafeLRU_Size-12         	2000000000	         1.50 ns/op
BenchmarkThreadUnsafeLRU_Size2-12        	2000000000	         1.49 ns/op
BenchmarkThreadUnsafeLRU_Find-12         	20000000	       178 ns/op
BenchmarkThreadUnsafeLRU_Find2-12        	10000000	       177 ns/op
BenchmarkThreadUnsafeLRU_Find3-12        	30000000	        34.7 ns/op
BenchmarkThreadUnsafeLRU_Remove-12       	20000000	       177 ns/op
BenchmarkThreadUnsafeLRU_Remove2-12      	10000000	       178 ns/op
BenchmarkThreadUnsafeLRU_Remove3-12      	50000000	        34.4 ns/op
BenchmarkThreadUnsafeLRU_Iterator-12     	   10000	    149954 ns/op
BenchmarkThreadUnsafeLRU_Iterator2-12    	 1000000	      2343 ns/op
BenchmarkThreadUnsafeLRU_Iter-12         	   10000	    152086 ns/op
BenchmarkThreadUnsafeLRU_Iter2-12        	 1000000	      1803 ns/op
```
