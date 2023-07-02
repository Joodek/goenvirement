package goenv

import (
	"log"
)

const MAX_ALLOWED_CALLS = 1000

var stack = make(map[string]int)

// to avoid infinit loops , in case one key reads the same key multiple times
// which is not a normal in env files
// so we record each call and panic if it reaches the max allowed calls for each key
func throttle(key1 string, key2 string) {

	calls := stack[key1+key2]

	if calls > MAX_ALLOWED_CALLS {
		log.Panic("recursion detected : trying to read " + key2 + " by " + key1)
	}

	stack[key1+key2]++
}
