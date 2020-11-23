package main

import (
	"errors"
	"fmt"
	"github.com/tarantool/go-tarantool"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)


var conn = &tarantool.Connection{}
var err = errors.New("some err")

func insertLines() {

	var mutex sync.Mutex

	var linesCount uint64 = 1000

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100000; j++{

				mutex.Lock()

				atomic.AddUint64(&linesCount, 1)

				_, err = conn.Insert("tester", []interface{}{ linesCount, linesCount, "staticLink123321123321123312", linesCount % 2 == 0, 300, 10,10, linesCount})
				if err != nil {
					log.Fatalf("conn.Insert err:", err)
				}
				_, err = conn.Insert("target", []interface{}{ linesCount, linesCount, linesCount%5, linesCount , 5, linesCount,10})
				if err != nil {
					log.Fatalf("conn.Insert err:", err)
				}
				mutex.Unlock()

			}
		}()
	}
	for {
		time.Sleep(time.Second * 2)
		linesCheck := atomic.LoadUint64(&linesCount)
		log.Println("lines count:", linesCheck)
		if linesCheck > 1000000 { // быстрее остановится  из-за памяти
			break
		}
	}
	linesFinal := atomic.LoadUint64(&linesCount)
	fmt.Println("linesCount:", linesFinal)

}

func selectLines() {

	var linesCount uint64 = 0


	maxValue := 1000000 // paste len from tarantool

	timeSecondTest := uint64(40)

	for i := 0; i < 30; i++ {
		go func() {
			for j := 0; j < 1000000; j++{

				targetId := rand.Intn(maxValue - 40) + 40

				//// use index
				//resp, err := conn.Select("tester", "secondary", 0, 1, tarantool.IterEq, []interface{}{targetId})
				//if err != nil {
				//	log.Println("conn.Select err:", err)
				//	log.Fatalf("conn.Select reps:", resp)
				//}

				// index and full read
				resp, err := conn.Call(`myget`, []interface{}{strconv.Itoa(targetId), targetId})
				if err != nil {
					log.Println("conn.Select err:", err)
					log.Fatalf("conn.Select reps:", resp)
				}
				atomic.AddUint64(&linesCount, 1)

			}
		}()
	}

	time.Sleep(time.Duration(timeSecondTest) * time.Second)
	linesFinal := atomic.LoadUint64(&linesCount)

	fmt.Println("linesCount:", linesFinal)
	fmt.Println("RPS:", linesFinal/timeSecondTest)

}



func main() {


	conn, err = tarantool.Connect("0.0.0.0:3301", tarantool.Opts{
		User: "admin",
		Pass: "admin",
	})
	if err != nil {
		log.Fatalf("Connection refused:", err)
	}
	defer conn.Close()


	//insertLines()

	selectLines()



}
