package main

import (
	"os"
	"time"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"rpcshared"
	// Log results to our webserver
	"rpclogger"
)

var (
	MyName     string
	BrokerHost string
	MyType     string
)

func init() {
	MyName = Generate(2, "-")
	BrokerHost = os.Getenv("BROKERHOST")
	if len(BrokerHost) == 0 {
		BrokerHost = "trex1:5050"
	}
	MyType = "BulkExtractor"
}

// Send the whole request history periodically
// TODO: Decay the RequestHistory buffer. This struct will eventually get huge..
func PeriodicUpdate(myRPCInstance *rpcshared.BulkExtractor) {
	for {
		time.Sleep(time.Millisecond * 5000)
		rpclogger.SubmitReport(BrokerHost, MyName, MyType, myRPCInstance.RequestHistory)
	}
}

func startServer() {
	fmt.Println("started server")
	be := new(rpcshared.BulkExtractor)
	rpc.Register(be)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":5555")
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	go http.Serve(l, nil)
	// go PeriodicUpdate(be)
}

//Start the server, listen forever.
func main() {
    numThreads := os.Getenv("BE_THREADS")
    if len(numThreads)==0{
        numThreads="1"
    }
    fmt.Println("Number of threads to provide BE based on env:", numThreads)

	startServer()
	meta := make(chan int)
	x := <-meta /// wait for a while, and listen
	fmt.Println(x)
}
