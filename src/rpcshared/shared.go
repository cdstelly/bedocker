package rpcshared

import (
	"strings"
    "fmt"
    "os/exec"
    "log"
    "bytes"
	"io/ioutil"
	"encoding/json"
)

type BulkExtractor string

type Args struct {
	Data   []byte
	DataID string
}

func (t *BulkExtractor) Extract(args *Args, reply *string) error {
    pathToTool := "/usr/local/bin/bulk_extractor"
    fmt.Println("Path to BulkExtractor: " , pathToTool)

	bulk_Input_Directory  := "/tmp/bulk_in/"
	bulk_Output_Directory := "/tmp/bulk_out/"    //make sure this does not exist beforehand per BE docs

	// Bulk Extractor works by inspecting data in a directory (or a raw disk dump), setup the passed in data as such:
	filepath := bulk_Input_Directory + args.DataID
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

    err = cmd.Run()
	fmt.Println("[-] Output: ", out.String())


	//Post process the BE output
	jsonMapping := make(map[string]string)
	files, err := ioutil.ReadDir(bulk_Output_Directory)
    for _, f := range files {
		if f.Size() > 0 {
			fmt.Println("File ", f.Name(), " is not null")
			filedata,ferr := ioutil.ReadFile(bulk_Output_Directory + f.Name())
			if ferr != nil {
				fmt.Println("was not able to open file at: ", f.Name(), ferr)
			}
			// put data in json, append to running list of stuff
			jsonKey := strings.Split(f.Name(), ".")[0]
			jsonValue := "bulkextract-" + string(filedata)
			jsonMapping[jsonKey] = jsonValue
		}
    }

	//for key, value := range jsonMapping {
	//	fmt.Println("Key:", key, "Value:", value)
	//}

	jsonString, err := json.Marshal(jsonMapping)
    if err != nil {
            log.Println(err)
	}
	fmt.Println(string(jsonString))

	*reply = out.String()
    if err != nil {
            log.Println(err)
    }
    return err
}

