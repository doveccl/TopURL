package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"runtime"
	"time"
)

const TOP_K = 100
const HASH_MOD = 131

// pair of url & hash
type SI struct {
	str []byte
	hash int
}

// 字符串哈希
func string_hash(str []byte) int {
	h := fnv.New32a()
	h.Write(str)
	return int(h.Sum32() % HASH_MOD)
}

// 从指定文件中逐行读取 url，并写入 chan
func read(file string, c chan []byte) {
	defer close(c)
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		c <- line
	}
}

// 根据 url 哈希分组，并写入 chan
func make_group(u chan []byte, h chan SI) {
	cores := runtime.NumCPU()
	f := make(chan int, cores)
	worker := func() {
		for {
			v, ok := <- u
			if !ok { break }
			h <- SI{v, string_hash(v)}
		}
		f <- 1
	}
	for i := 0; i < cores; i++ {
		go worker()
	}
	for i := 0; i < cores; i++ {
		<- f
	}
	close(h)
}

// 将分好组的 url 写入对应的小文件中
func write(h chan SI, f chan bool) {
	var err error
	var files [HASH_MOD]*os.File
	for i := 0; i < HASH_MOD; i++ {
		fn := fmt.Sprintf("HASH_%d", i)
		files[i], err = os.Create(fn)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}
	for {
		v, ok := <- h
		if !ok { break }
		files[v.hash].Write(v.str)
		files[v.hash].Write([]byte("\n"))
	}
	for i := 0; i < HASH_MOD; i++ {
		files[i].Close()
	}
	f <- true
}

// 从小文件中读取分好组的 url 并进行结果统计
func solve() {
	cnt := 0
	h := &PSIHeap{}
	heap.Init(h)
	for i := 0; i < HASH_MOD; i++ {
		urls := make(chan []byte, 1e6)
		mp := make(map[string]int)
		file := fmt.Sprintf("HASH_%d", i)
		go read(file, urls)
		for {
			v, ok := <- urls
			s := string(v)
			if !ok { break }
			mp[s] = mp[s] + 1
		}
		for k, v := range mp {
			heap.Push(h, PSI{k, v})
			cnt++
			for cnt > TOP_K {
				heap.Pop(h)
				cnt--
			}
		}
		os.Remove(file)
	}
	for i := 0; i < h.Len(); i++ {
		fmt.Printf("%s %d\n", (*h)[i].s, (*h)[i].i)
	}
}

// 计算运行时间
func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

func main() {
	defer elapsed("program")()

	file := flag.String("file", "", "big url file path")
	flag.Parse()

	urls := make(chan []byte, 1e6)
	hashs := make(chan SI, 1e6)
	finish := make(chan bool)
	go read(*file, urls)
	go make_group(urls, hashs)
	go write(hashs, finish)

	if <- finish { solve() }
}
