package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"time"
)

const TOP_K = 100
const HASH_MOD = 131

// 存放 url
var urls chan []byte
// 存放对应哈希
var hashs [HASH_MOD]chan []byte
// url 文件名
var url_file string

// 初始化
func init() {
	// 创建临时目录
	os.Mkdir("./tmp", 0755)
	// 全局变量初始化
	urls = make(chan []byte, 1e6)
	for i := 0; i < HASH_MOD; i++ {
		hashs[i] = make(chan []byte, 1e4)
	}
}

// 根据哈希值获取文件名
func file_name(x int) string {
	return fmt.Sprintf("./tmp/%d", x)
}

// 字符串哈希
func string_hash(str []byte) uint32 {
	h := fnv.New32a()
	h.Write(str)
	return h.Sum32() % HASH_MOD
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

// 根据 url 哈希分组，并写入对应 chan
func make_group() {
	for {
		v, ok := <- urls
		if !ok { break }
		hashs[string_hash(v)] <- v
	}
	for i := 0; i < HASH_MOD; i++ {
		close(hashs[i])
	}
}

// 将分好组的 url 写入对应的小文件中
func wait_write() {
	var err error
	var files [HASH_MOD]*os.File
	finish := make(chan bool, HASH_MOD)
	for i := 0; i < HASH_MOD; i++ {
		fn := file_name(i)
		files[i], err = os.Create(fn)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		go func(x int) {
			for {
				v, ok := <- hashs[x]
				if !ok { break }
				files[x].Write(append(v, '\n'))
			}
			finish <- true
		}(i)
	}
	defer close(finish)
	for i := 0; i < HASH_MOD; i++ {
		<- finish
	}
}

// 从小文件中读取分好组的 url 并进行结果统计
func solve() {
	cnt := 0
	h := &PSIHeap{}
	heap.Init(h)
	for i := 0; i < HASH_MOD; i++ {
		urls := make(chan []byte, 1e6)
		mp := make(map[string]int)
		file := file_name(i)
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
	}
	for i := 0; i < h.Len(); i++ {
		fmt.Printf("%s %d\n", (*h)[i].s, (*h)[i].i)
	}
	os.RemoveAll("../tmp")
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

	go read(*file, urls)
	go make_group()
	wait_write()

	solve()
}
