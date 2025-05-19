package cli

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type CLI struct{}

type fileRef struct {
	fileCID         string
	previousVersion int
}

func (cli *CLI) Upload(file string) {
	dir, _ := os.Getwd()
	referencePath := filepath.Join(dir, "reference", file)
	fmt.Println(referencePath)
	_, err := os.Stat(referencePath)

	if err != nil {
		//reference file not found, import it into FileDAG
		cmd := exec.Command("./lotus", "client", "import", file)
		out, _ := cmd.Output()
		//The output of this command follows the format Import XXXXXX, Root CID, we now spilt the string to get the CID
		parts := strings.Split(string(out), "Root ")
		if len(parts) < 2 {
			fmt.Println("Error importing file")
			return
		}
		cid := strings.TrimSpace(parts[1])
		fmt.Println(cid)
		//use the CID of the file to upload it to the storage miner
		cmd2 := exec.Command("./lotus", "client", "deal", cid, "t01000", "0.026", "518400")
		out2, _ := cmd2.Output()
		fmt.Println(string(out2))
		//create the reference of the file to manage the state of the file
		fileReference := []fileRef{}
		fileReference = append(fileReference, fileRef{fileCID: cid, previousVersion: -1})
		WriteReference(file, fileReference)
		fmt.Println("Uploaded ", file, "(version 1) to storage miner.")
	} else {
		//reference file found, upload it to FileDAG
		fmt.Println("the file ", file, "already exist in FileDAG. If you want to update it please use the update/fork/merge command.")
	}
}

func (cli *CLI) Retrieve(file string, version int) {
	dir, _ := os.Getwd()
	referencePath := filepath.Join(dir, "reference", file)
	fmt.Println(referencePath)
	_, err := os.Stat(referencePath)

	if err != nil {
		//reference file not found, the user might input a wrong filename or need to upload the file first
		fmt.Println("unable to the file ", file, "you might give a wrong filename or need to use the upload command to upload it.")
	} else {
		//open the reference file to get the CID and previousVersion of each patch of the file
		filePrinter, _ := os.Open(referencePath)
		defer filePrinter.Close()
		//scan the file line by line and store it in an array
		fileReference := []fileRef{}
		scanner := bufio.NewScanner(filePrinter)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			parts := strings.Fields(line)
			preVersion, _ := strconv.Atoi(parts[1])
			reference := fileRef{fileCID: parts[0], previousVersion: preVersion}
			fileReference = append(fileReference, reference)
		}
		//test scan
		for i := 0; i < len(fileReference); i++ {
			fmt.Println(fileReference[i].fileCID, fileReference[i].previousVersion)
		}
		//test scan
		//Now retrieve the target versions
		versionToRetrieve := version
		totalRetrieve := 0
		for {
			if versionToRetrieve == -1 {
				break
			}
			cidToRetrieve := fileReference[versionToRetrieve-1].fileCID
			retrieveFileName := strconv.Itoa(totalRetrieve)
			cmd := exec.Command("./lotus", "client", "retrieve", cidToRetrieve, retrieveFileName)
			out, _ := cmd.Output()
			//test command exec
			fmt.Println(string(out))
			//test command exec
			versionToRetrieve = fileReference[versionToRetrieve-1].previousVersion
			totalRetrieve += 1
		}
		//Finally, put the original version and patches together
		if totalRetrieve == 1 {
			//retrieve the original file, just rename the file
			retrieveName := file + "_v" + strconv.Itoa(version)
			os.Rename(strconv.Itoa(0), retrieveName)
		} else {
			tempFile := strconv.Itoa(rand.Int())
			oldFile := strconv.Itoa(totalRetrieve - 1)
			patchFile := strconv.Itoa(totalRetrieve - 2)
			for {
				cmd2 := exec.Command("bspatch", oldFile, tempFile, patchFile)
				out2, _ := cmd2.Output()
				fmt.Println(string(out2))
				os.Remove(oldFile)
				os.Remove(patchFile)
				totalRetrieve -= 1
				if totalRetrieve == 1 {
					break
				}
				oldFile = tempFile
				tempFile = strconv.Itoa(rand.Int())
				patchFile = strconv.Itoa(totalRetrieve - 2)
			}
			retrieveName := file + "_v" + strconv.Itoa(version)
			os.Rename(tempFile, retrieveName)
		}
	}
}

func WriteReference(file string, fileReferences []fileRef) {
	dir, _ := os.Getwd()
	referencePath := filepath.Join(dir, "reference", file)
	filePrinter, _ := os.Create(referencePath)
	defer filePrinter.Close()
	for _, ref := range fileReferences {
		line := fmt.Sprintf("%s %d", ref.fileCID, ref.previousVersion)
		filePrinter.WriteString(line)
	}
}

func (cli *CLI) Run() {
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	retrieveCmd := flag.NewFlagSet("retrieve", flag.ExitOnError)

	uploadFile := uploadCmd.String("file", "", "the file to upload")
	retrieveFile := retrieveCmd.String("file", "", "the file to retrieve")
	retrieveVersion := retrieveCmd.Int("version", 1, "the version of the file to retrieve, default version 1")

	switch os.Args[1] {
	case "upload":
		err := uploadCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing upload command")
			return
		}
		fmt.Println(*uploadFile)
		cli.Upload(*uploadFile)
	case "retrieve":
		err := retrieveCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing retrieve command")
			return
		}
		fmt.Println(*retrieveFile)
		fmt.Println(*retrieveVersion)
		cli.Retrieve(*retrieveFile, *retrieveVersion)
	default:
		fmt.Println("Invalid command")
	}
}
