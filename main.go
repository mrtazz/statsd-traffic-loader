// statsd traffic loader

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	COUNTER = iota
	TIMER
	NOTYPE
)

const GOROUTINE_NUM = 100

type randomDataMaker struct {
	src rand.Source
}

func main() {
	// note, that variables are pointers
	packet_rate := flag.Int("s", 30000, "packets per second to send")
	port := flag.Int("p", 8125, "port to send to")
	counters_file := flag.String("c", "stats_counter_keys.txt", "file with example counter keys")
	timers_file := flag.String("t", "stats_timer_keys.txt", "file with example timer keys")

	flag.Usage = usage

	flag.Parse()

	hostname := flag.Arg(0)

	if hostname == "" {
		flag.Usage()
	}

	rand.Seed(time.Now().Unix())

	counters, err := readLines(*counters_file)
	if err != nil {
		fmt.Printf("Error reading counters file: %v\n", err)
		os.Exit(1)
	}
	timers, err := readLines(*timers_file)
	if err != nil {
		fmt.Printf("Error reading timers file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sending %d packets/s to %s on port %d.\n",
		*packet_rate, hostname, *port)
	ticker := time.NewTicker(time.Second)
	sendPackets(ticker.C, hostname, *port, *packet_rate, counters, timers)

}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: statsd-traffic-loader [-sp] hostname\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func sendPackets(timer <-chan time.Time, hostname string, port int, count int,
	counters []string, timers []string) {
	connectionString := fmt.Sprintf("%s:%d", hostname, port)
	per_goroutine := count / GOROUTINE_NUM
	fmt.Printf("Sending %d packets per go routine.\n", per_goroutine)
	conn, _ := net.Dial("udp", connectionString)
	randomSrc := randomDataMaker{rand.NewSource(1028890720402726901)}

	for {
		select {
		case <-timer:
			// do stuff
			for x := 0; x < GOROUTINE_NUM; x++ {
				go sendStatsdPacket(randomSrc, per_goroutine, counters, timers, conn)
			}
		}
	}
}

func sendStatsdPacket(random randomDataMaker, packets int, counters []string,
	timers []string, conn net.Conn) {
	for x := 0; x < packets; x++ {
		value := rand.Int() & 0xff
		var update_string string
		thetype := rand.Intn(NOTYPE)
		switch thetype {
		case TIMER:
			metricname := timers[rand.Intn(len(timers))]
			update_string = fmt.Sprintf("%s:%d|ms", metricname, value)
		case COUNTER:
			metricname := counters[rand.Intn(len(counters))]
			update_string = fmt.Sprintf("%s:%d|c", metricname, value)
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

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
