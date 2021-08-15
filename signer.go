package main

import (
	"sync"
)

const TH = 6

func ExecutePipeline(jobs ...job) {

}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}

	for data := range in {
		md5sum := DataSignerMd5(data.(string))

		wg.Add(1)
		go workerSingleHash(wg, data.(string), md5sum, out)
	}
	wg.Wait()
}

func workerSingleHash(wg *sync.WaitGroup, data string, md5sum string, out chan interface{}) {
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

}

func CombineResults(in, out chan interface{}) {

}
