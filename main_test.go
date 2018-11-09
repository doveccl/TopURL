package main

import (
	"testing"
)

func TestHash(t *testing.T) {
	hash := string_hash([]byte("hello"))
	if hash != 1335831723 % HASH_MOD {
		t.Error("hash error")
	}
}
