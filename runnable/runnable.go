package runnable

import (
	"log"
	"math/big"
	"time"
)

func factorial(n *big.Int) (result *big.Int) {
	defer timeTrack(time.Now(), "factorial")
	//in seconds to sleep.
	time.Sleep(8 * time.Second) //
	return n
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
