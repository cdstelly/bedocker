package rpcshared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type BulkExtractor struct {
	NumberRequests int
	RequestHistory []float64
}

type Args struct {
	Data   []byte
	DataID string
}

func (t *BulkExtractor) Extract(args *Args, reply *string) error {
	pathToTool := "/usr/local/bin/bulk_extractor"
	fmt.Println("Path to BulkExtractor: ", pathToTool)

	bulk_Input_Directory := "/tmp/bulk_in/" + args.DataID + "/"
	bulk_Output_Directory := "/tmp/bulk_out/" + args.DataID + "/"

	// Bulk Extractor works by inspecting data in a directory (or a raw disk dump), setup the passed in data as such:
	os.MkdirAll(bulk_Input_Directory, 0644)
	filepath := bulk_Input_Directory + "data.dat"
	fmt.Println("Path to write given data to: ", filepath)
	err := ioutil.WriteFile(filepath, args.Data, 0644)
	if err != nil {
		log.Println("Error writing input to directory.", err)
	}

	//Setup the shell command to launch Bulk Extractor
	opts := []string{"-m", "1", "-b", "banner.txt", "-R", bulk_Input_Directory, "-o", bulk_Output_Directory}
	//Should look like the following: /usr/local/bin/bulk_extractor -m 1 -R /tmp/bulk_in/ -o /tmp/bulk_out/
	cmd := exec.Command(pathToTool, opts...)

	//Capture STDOUT
	var out bytes.Buffer
	cmd.Stdout = &out

	// Setup a timer
	startTime := time.Now()

	// Actually run the command:
	err = cmd.Run()
	fmt.Println("[-] Output: ", out.String())

	//Post process the BE output
	jsonMapping := make(map[string]string)
	files, err := ioutil.ReadDir(bulk_Output_Directory)
	for _, f := range files {
		if f.Size() > 0 {
			fmt.Println("File ", f.Name(), " is not null")
			filedata, ferr := ioutil.ReadFile(bulk_Output_Directory + f.Name())
			if ferr != nil {
				fmt.Println("was not able to open file at: ", f.Name(), ferr)
			}
			// put data in json, append to running list of stuff
			jsonKey := strings.Split(f.Name(), ".")[0]
			jsonValue := "bulkextract-" + string(filedata)
			jsonMapping[jsonKey] = jsonValue
		}
	}

	// Dump everything into JSON in preperation for Elasticsearch upload
	jsonString, err := json.Marshal(jsonMapping)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(jsonString))

	//We want to return the JSON in addition to STDOUT
	*reply = out.String()
	if err != nil {
		log.Println(err)
	}

	// Save execution time to history
	executionTime := time.Since(startTime).Seconds() //use seconds as opposed to nanoseconds, returns float64 which is required with stats package
	t.NumberRequests += 1
	t.RequestHistory = append(t.RequestHistory, executionTime)

	// If all goes well, remove temp directories
	remerr := os.RemoveAll(bulk_Input_Directory)
	remerr = os.RemoveAll(bulk_Output_Directory)
	if remerr != nil {
		fmt.Println("Error cleaning up temporary directories: ", remerr)
	}
	return nil
}

func (t *BulkExtractor) GetHistory(args *Args, reply *[]float64) error {
	*reply = t.RequestHistory
	return nil
}
