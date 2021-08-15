package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

const TH = 6

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})

	wg := sync.WaitGroup{}

	for _, jobFunc := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		go worker(jobFunc, in, out, &wg)
		in = out
	}
	wg.Wait()
}

func worker(jobFunc job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	jobFunc(in, out)
}

func SingleHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}

	for data := range in {
		dataString := fmt.Sprintf("%v", data)
		md5sum := DataSignerMd5(dataString)

		wg.Add(1)
		go workerSingleHash(dataString, md5sum, out, &wg)
	}
	wg.Wait()
}

func workerSingleHash(data string, md5sum string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	crc32ch := make(chan string)
	crc32md5ch := make(chan string)

	go calculateHash(data, crc32ch)
	go calculateHash(md5sum, crc32md5ch)

	crc32data := <-crc32ch
	crc32md5sum := <-crc32md5ch

	out <- crc32data + "~" + crc32md5sum
}

func calculateHash(data string, ch chan string) {
	ch <- DataSignerCrc32(data)
}

func MultiHash(in, out chan interface{}) {
	wg := sync.WaitGroup{}

	for data := range in {
		dataString := fmt.Sprintf("%v", data)

		wg.Add(1)
		go workerMultiHash(dataString, out, &wg)
	}
	wg.Wait()
}

func workerMultiHash(data string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	wgroup := sync.WaitGroup{}

	hashes := make([]string, TH)

	for i := 0; i < TH; i++ {
		dataPrefix := fmt.Sprintf("%v%v", i, data)
		wgroup.Add(1)
		go calculateMultiHash(dataPrefix, i, hashes, &wgroup)
	}

	wgroup.Wait()

	out <- strings.Join(hashes, "")
}

func calculateMultiHash(data string, th int, hashes []string, wgroup *sync.WaitGroup) {
	defer wgroup.Done()
	hashes[th] = DataSignerCrc32(data)
}

func CombineResults(in, out chan interface{}) {
	var hashes []string

	for data := range in {
		dataString := fmt.Sprintf("%v", data)
		hashes = append(hashes, dataString)
	}

	sort.Strings(hashes)

	out <- strings.Join(hashes, "_")
}
