package main

import (
	"hash/crc32"
	"sync"
)

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

	crc32ch <- DataSignerCrc32(data)
	crc32md5ch <- DataSignerCrc32(md5sum)

	crc32data := <-crc32ch
	crc32md5sum := <-crc32md5ch

	out <- crc32data + "~" + crc32md5sum
}

func MultiHash(in, out chan interface{}) {

}

func CombineResults(in, out chan interface{}) {

}
