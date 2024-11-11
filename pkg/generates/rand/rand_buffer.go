package rand

import (
	"crypto/rand"
)

const (
	gBUFFER_SIZE = 10000
)

var (
	bufferChan = make(chan []byte, gBUFFER_SIZE)
)

func init() {
	go asyncProducingRandomBufferBytesLoop()
}

func asyncProducingRandomBufferBytesLoop() {
	var step int
	for {
		buffer := make([]byte, 1024)
		if n, err := rand.Read(buffer); err != nil {
			panic(err)
		} else {
			for _, step = range []int{4} {
				for i := 0; i <= n-4; i += step {
					bufferChan <- buffer[i : i+4]
				}
			}
		}
	}
}
