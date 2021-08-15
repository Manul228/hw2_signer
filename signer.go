package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

const TH = 6

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	wg := new(sync.WaitGroup)

	for _, jobFunc := range jobs {
		wg.Add(1)
		out := make(chan interface{})
		go worker(jobFunc, in, out, wg)
		in = out
	}
	wg.Wait()
}

func worker(jobFunc job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	go jobFunc(in, out)
}

func SingleHash(in, out chan interface{}) {
	wg := new(sync.WaitGroup)

	for data := range in {
		md5sum := DataSignerMd5(data.(string))

		wg.Add(1)
		go workerSingleHash(data.(string), md5sum, out, wg)
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
	wg := new(sync.WaitGroup)

	for data := range in {
		wg.Add(1)
		workerMultiHash(data.(string), out, wg)
	}
}

func workerMultiHash(data string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	wgroup := new(sync.WaitGroup)

	hashes := make([]string, TH)

	for i := 0; i < TH; i++ {
		wgroup.Add(1)
		go calculateMultiHash(data, i, hashes, wgroup)
	}

	wg.Wait()

	out <- strings.Join(hashes, "")
}

func calculateMultiHash(data string, th int, hashes []string, wg *sync.WaitGroup) {
	hashes[th] = strconv.Itoa(th) + DataSignerCrc32(data)
	wg.Done()
}

func CombineResults(in, out chan interface{}) {

	var hashes []string

	for data := range in {
		hashes = append(hashes, data.(string))
	}

	sort.Strings(hashes)

	out <- strings.Join(hashes, "_")
}
