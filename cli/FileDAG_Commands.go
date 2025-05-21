package cli

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type fileRef struct {
	fileCID         string
	previousVersion int
}

func (cli *CLI) Upload(file string) {
	//check whether the reference file exist to get whether the file exist in FileDAG
	dir, _ := os.Getwd()
	referencePath := filepath.Join(dir, "reference", file)
	//fmt.Println(referencePath)
	_, err := os.Stat(referencePath)

	if err != nil {
		//reference file not found, first import it into FileDAG
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

		//save the reference of the file
		WriteReference(file, fileReference)

		fmt.Println("Uploaded", file, "( version 1 ) to storage miner.")
	} else {
		//reference file found, upload it to FileDAG
		fmt.Println("the file", file, "already exist in FileDAG. If you want to update it please use the update/fork/merge command.")
	}
}

func (cli *CLI) Retrieve(file string, version int) {
	//check whether the reference file exist to get whether the file exist in FileDAG
	dir, _ := os.Getwd()
	referencePath := filepath.Join(dir, "reference", file)
	//fmt.Println(referencePath)
	_, err := os.Stat(referencePath)

	if err != nil {
		//reference file not found, the user might input a wrong filename or need to upload the file first
		fmt.Println("unable to retrieve the file", file, ". You might give a wrong filename or need to use the upload command to upload it.")
	} else {
		//open the reference file to get the CID and previousVersion of each patch of the file
		filePrinter, _ := os.Open(referencePath)
		defer filePrinter.Close()

		//scan the file line by line and store it in an slide
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
		//for i := 0; i < len(fileReference); i++ {
		//fmt.Println(fileReference[i].fileCID, fileReference[i].previousVersion)
		//}

		//Now retrieve the target patches and original file
		versionToRetrieve := version
		totalRetrieve := 0
		for {
			//we retrieve in reverse order, i.e. first retrieve the lastest patch of the file, then determine its previous version
			//based on the reference until we finally retrieve the original file
			if versionToRetrieve == -1 {
				break
			}
			cidToRetrieve := fileReference[versionToRetrieve-1].fileCID
			retrieveFileName := strconv.Itoa(totalRetrieve)
			cmd := exec.Command("./lotus", "client", "retrieve", cidToRetrieve, retrieveFileName)
			out, _ := cmd.Output()
			fmt.Println(string(out))
			versionToRetrieve = fileReference[versionToRetrieve-1].previousVersion
			totalRetrieve += 1
		}

		//Finally, put the original version and patches together
		if totalRetrieve == 1 {
			//in this case, we only retrieve the original file, just rename the file
			retrieveName := file + "_v" + strconv.Itoa(version)
			os.Rename(strconv.Itoa(0), retrieveName)
			fmt.Println("successfully retrieve the file, saved as", retrieveName)
		} else {
			tempFile := strconv.Itoa(rand.Int())
			oldFile := strconv.Itoa(totalRetrieve - 1)
			patchFile := strconv.Itoa(totalRetrieve - 2)
			for {
				cmd2 := exec.Command("bspatch", oldFile, tempFile, patchFile)
				out2, _ := cmd2.Output()
				fmt.Println(string(out2))

				//delete the temporary files generated during the process
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
			fmt.Println("successfully retrieve the file, saved as", retrieveName)
		}
	}
}

func (cli *CLI) Update(file string, version int, updateFile string) {
	//check whether the reference file exist to get whether the file exist in FileDAG
	dir, _ := os.Getwd()
	referencePath := filepath.Join(dir, "reference", file)
	//fmt.Println(referencePath)
	_, err := os.Stat(referencePath)

	if err != nil {
		//reference file not found
		fmt.Println("unable to update the file", file, "you might give a wrong filename or need to use the upload command to upload it.")
	} else {
		//reference file found, update the file
		//first, retrieve the base version of the file
		cli.Retrieve(file, version)

		//then, generate the patch
		oldVersionName := file + "_v" + strconv.Itoa(version)
		rand.Seed(time.Now().UnixNano())
		patchName := "patch" + strconv.Itoa(rand.Int())
		//fmt.Println(patchName)
		cmd := exec.Command("bsdiff", oldVersionName, updateFile, patchName)
		cmd.Output()

		//next, import the patch to get its CID and upload it
		cmd2 := exec.Command("./lotus", "client", "import", patchName)
		out, _ := cmd2.Output()

		//The output of this command follows the format Import XXXXXX, Root CID, we now spilt the string to get the CID
		parts := strings.Split(string(out), "Root ")
		if len(parts) < 2 {
			fmt.Println("Error importing file")
			return
		}
		cid := strings.TrimSpace(parts[1])
		fmt.Println(cid)
		cmd3 := exec.Command("./lotus", "client", "deal", cid, "t01000", "0.026", "518400")
		out3, _ := cmd3.Output()
		fmt.Println(string(out3))

		//finally, update the reference of the file and save this update
		filePrinter, _ := os.Open(referencePath)
		defer filePrinter.Close()

		//scan the file line by line and store it in an slide
		fileReference := []fileRef{}
		scanner := bufio.NewScanner(filePrinter)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			parts := strings.Fields(line)
			preVersion, _ := strconv.Atoi(parts[1])
			reference := fileRef{fileCID: parts[0], previousVersion: preVersion}
			fileReference = append(fileReference, reference)
		}

		//update the new referencw
		fileReference = append(fileReference, fileRef{fileCID: cid, previousVersion: version})
		versionNum := strconv.Itoa(len(fileReference))

		//write the updated reference back to the file
		WriteReference(file, fileReference)
		fmt.Println("Updated", file, "( version", versionNum, ") to storage miner.")

		//delete the temporary files generated during the process
		os.Remove(oldVersionName)
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
		filePrinter.WriteString("\n")
	}
}
