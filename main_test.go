package main

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	hash := string_hash([]byte("hello"))
	if hash != 1335831723 % HASH_MOD {
		t.Error("hash error")
	}
}

func TestIO(t *testing.T) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	for i := 0; i < 10; i++ {
		h := r.Intn(HASH_MOD)
		f := fmt.Sprintf("HASH_%d", h)
		cnt := r.Intn(1e4)
		ch1 := make(chan SI, cnt)
		ok := make(chan bool)
		go write(ch1, ok)
		for cnt > 0 {
			ch1 <- SI{[]byte(f), h}
			cnt--
		}
		close(ch1)
		if <-ok {
			ch2 := make(chan []byte)
			go read(f, ch2)
			for {
				v, ok := <- ch2
				if !ok { break }
				if string(v) != f {
					t.Error("io error")
				}
			}
		}
	}
	for i := 0; i < HASH_MOD; i++ {
		f := fmt.Sprintf("HASH_%d", i)
		os.Remove(f)
	}
}
