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
import "github.com/Ninlgde/lrucache/go
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