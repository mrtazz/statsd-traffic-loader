// statsd traffic loader

package main

import (
	"flag"
	"fmt"
	"time"
	"os"
	"net"
	"math/rand"
	"github.com/jmcvetta/randutil"
)

const (
  COUNTER = iota
  TIMER
  SET
  GAUGE
  NOTYPE
)

const GOROUTINE_NUM = 100

type randomDataMaker struct {
  src rand.Source
}


func main() {
	// note, that variables are pointers
	packet_rate := flag.Int("s", 30000, "packets per second to send")
	port        := flag.Int("p", 8125, "port to send to")

	flag.Usage = usage

	flag.Parse()

  hostname := flag.Arg(0)

  if (hostname == "") {
    flag.Usage()
  }

  rand.Seed(time.Now().Unix())

	fmt.Printf("Sending %d packets/s to %s on port %d.\n",
	*packet_rate, hostname, *port)
  ticker := time.NewTicker(time.Second)
  sendPackets(ticker.C, hostname, *port, *packet_rate)

}

func usage() {
  fmt.Fprintf(os.Stderr, "usage: statsd-traffic-loader [-sp] hostname\n")
  flag.PrintDefaults()
  os.Exit(1)
}

func sendPackets(timer <-chan time.Time, hostname string, port int, count int) {
	connectionString := fmt.Sprintf("%s:%d", hostname, port)
  per_goroutine := count % GOROUTINE_NUM
	conn, _ := net.Dial("udp", connectionString)
  randomSrc := randomDataMaker{rand.NewSource(1028890720402726901)}
  for {
    select {
    case <- timer:
      // do stuff
      for x := 0; x < GOROUTINE_NUM; x++ {
        go sendStatsdPacket(randomSrc, per_goroutine, conn)
      }
    }
  }
}

func sendStatsdPacket(random randomDataMaker, packets int, conn net.Conn) {
  for x := 0; x < packets; x++ {
    metricname, _ := randutil.AlphaString(32)
    value := rand.Int() & 0xff
    var update_string string
    thetype := rand.Intn(NOTYPE)
    switch thetype {
    case TIMER:
      update_string = fmt.Sprintf("%s:%d|ms", metricname, value)
    case COUNTER:
      update_string = fmt.Sprintf("%s:%d|ms", metricname, value)
    case SET:
      update_string = fmt.Sprintf("%s:%d|ms", metricname, value)
    case GAUGE:
      update_string = fmt.Sprintf("%s:%d|ms", metricname, value)
    }
    fmt.Fprintf(conn, update_string)
  }
}

func (r *randomDataMaker) Read(p []byte) (n int, err error) {
  for i := range p {
    p[i] = byte(r.src.Int63() & 0xff)
  }
  return len(p), nil
}

