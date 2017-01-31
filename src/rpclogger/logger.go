package rpclogger

import (
	"log"
	"net/rpc"
)

type RPCLogger string

type LogArgs struct {
	WorkerID        string
	WorkerType		string
	WorkerHistory []float64
}

func SubmitReport(targetHost string, id string, thistype string, history []float64) error {
	args := &LogArgs{WorkerID: id, WorkerType: thistype, WorkerHistory: history}

	client, err := rpc.DialHTTP("tcp", targetHost)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	var reply string
	err = client.Call("RPCLogger.Update", args, &reply)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return err
}
