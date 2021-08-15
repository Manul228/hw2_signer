package main

import (
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

	wgroup := &sync.WaitGroup{}

	wgroup.Add(1)
	go calculateHash(data, crc32ch, wgroup)

	wgroup.Add(1)
	go calculateHash(md5sum, crc32md5ch, wgroup)

	wgroup.Wait()

	crc32data := <-crc32ch
	crc32md5sum := <-crc32md5ch

	out <- crc32data + "~" + crc32md5sum
}

func calculateHash(data string, ch chan string, wg *sync.WaitGroup) {
	ch <- DataSignerCrc32(data)
	wg.Done()
}

func MultiHash(in, out chan interface{}) {

}

func CombineResults(in, out chan interface{}) {

}
