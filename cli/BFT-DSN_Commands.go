package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (cli *CLI) BFTUpload(file string, total int) {
	//check whether the reference file exist to get whether the file exist in BFT-DSN
	dir, _ := os.Getwd()
	referencePath := filepath.Join(dir, "reference", file)
	//fmt.Println(referencePath)
	_, err := os.Stat(referencePath)

	if err != nil {
		//reference file not found, upload it into BFT-DSN
		//generate the parameters of EC code
		f := (total - 1) / 3
		m := total
		k := m - f

		//create the EC shares of the file
		cmd := exec.Command("zfec", "-k", strconv.Itoa(k), "-m", strconv.Itoa(m), file)
		cmd.Output()

		//import each shares and upload them one by one
		fileReference := []fileRef{}

		for i := 0; i < total; i++ {
			//first import the share
			uploadShare := file + "." + strconv.Itoa(i) + "_" + strconv.Itoa(total) + ".fec"
			//fmt.Println(uploadShare)
			cmd := exec.Command("./lotus", "client", "import", uploadShare)
			out, _ := cmd.Output()

			//The output of this command follows the format Import XXXXXX, Root CID, we now spilt the string to get the CID
			parts := strings.Split(string(out), "Root ")
			if len(parts) < 2 {
				fmt.Println("Error importing file")
				return
			}
			cid := strings.TrimSpace(parts[1])
			fmt.Println(cid)

			//then use the cid to upload the share
			cmd2 := exec.Command("./lotus", "client", "deal", cid, "t01000", "0.026", "518400")
			out2, _ := cmd2.Output()
			fmt.Println(string(out2))

			//update file reference
			fileReference = append(fileReference, fileRef{fileCID: cid, previousVersion: -1})

			//stop 5 seconds before upload next share to avoid mistakes
			time.Sleep(5 * time.Second)
		}
		//save the reference of the file
		WriteReference(file, fileReference)
		fmt.Println("Uploaded", file, "to storage miner.")
	} else {
		//reference file found, upload it to FileDAG
		fmt.Println("the file ", file, " already exist in BFT-DSN. You can retrieve it use the BFTRetrieve command.")
	}
}

func (cli *CLI) BFTRetrieve(file string) {
	//check whether the reference file exist to get whether the file exist in BFT-DSN
	dir, _ := os.Getwd()
	referencePath := filepath.Join(dir, "reference", file)
	//fmt.Println(referencePath)
	_, err := os.Stat(referencePath)

	if err != nil {
		//reference file not found, exit
		fmt.Println("the file ", file, " is not exist in BFT-DSN. You can upload it use the BFTUpload command.")
	} else {
		//reference file found, then we can use the reference to retrieve the file
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

		//Now retrieve all the shares and recover the file
		shares := []string{}
		for i := 0; i < len(fileReference); i++ {
			shareName := "share" + strconv.Itoa(i)
			shares = append(shares, shareName)
			cmd := exec.Command("./lotus", "client", "retrieve", fileReference[i].fileCID, shareName)
			out, _ := cmd.Output()
			fmt.Println(string(out))
		}
		recoverFile := "recover_" + file
		args := append([]string{"-o", recoverFile}, shares...)
		cmd2 := exec.Command("zunfec", args...)
		cmd2.Output()

		//delete the temporary files generated during the process
		for i := 0; i < len(shares); i++ {
			os.Remove(shares[i])
		}
		fmt.Println("successfully retrieve the file, saved as", recoverFile)
	}
}
