package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CLI struct{}

type intslice []int

func (i *intslice) String() string {
	return fmt.Sprint(*i)
}

func (i *intslice) Set(value string) error {
	part := strings.Split(value, ",")
	for j := 0; j < len(part); j++ {
		v, err := strconv.Atoi(strings.TrimSpace(part[j]))
		if err != nil {
			return err
		}
		*i = append(*i, v)
	}
	return nil
}

func (cli *CLI) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("	upload -file filepath : upload the first version of a file to FileDAG, -file is the filepath of the upload file")
	fmt.Println("	update -file filepath -base version -new filepath: update the version of a file, -file is the filename you input in the upload command, -base is the integer version number to apply the update -new is the filepath of the new version of the file.")
	fmt.Println("	merge -file filepath -base version1,version2,... -new filepath: merge the version of a file, -file is the filename you input in the upload command, -base contains multiple integer version numbers to apply the nerge, saperated by comma -new is the filepath of the new version of the file.")
	fmt.Println("	retrieve -file filepath -version version : retrieve a specific version of a file, -file is the filename you input in the upload command, -version is the integer version number you want to retrieve.")
	fmt.Println("	BFTUpload -file filepath -total n: upload a file to BFT-DSN, -file is the filepath of the upload file, -total is the total number of erasure shares")
	fmt.Println("	BFTRetrieve -file filepath : retrieve a file from BFT-DSN, -file is the filename you input in the BFTUpload command")
}

func (cli *CLI) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.ValidateArgs()
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	retrieveCmd := flag.NewFlagSet("retrieve", flag.ExitOnError)

	//file update and file fork process are almost the same, so I merge them into one command and make client to choose which version (updateBase) to apply the update
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	mergeCmd := flag.NewFlagSet("merge", flag.ExitOnError)
	BFTUploadCmd := flag.NewFlagSet("BFTUpload", flag.ExitOnError)
	BFTRetrieveCmd := flag.NewFlagSet("BFTRetrieve", flag.ExitOnError)

	var base intslice
	var mergeFile string
	var mergeVersion string

	uploadFile := uploadCmd.String("file", "", "the file to upload")
	retrieveFile := retrieveCmd.String("file", "", "the file to retrieve")
	retrieveVersion := retrieveCmd.Int("version", 1, "the version of the file to retrieve, default version 1")
	updateFile := updateCmd.String("file", "", "the file to update")
	updateBase := updateCmd.Int("base", 1, "the base version of the file to apply the update, default version 1")
	updateVersion := updateCmd.String("new", "", "the updated file")
	mergeCmd.StringVar(&mergeFile, "file", "", "the file to merge")
	mergeCmd.Var(&base, "base", "the base versions of the file to apply merge, multiple versions saperated with comma(,)")
	mergeCmd.StringVar(&mergeVersion, "new", "", "the updated file")
	BFTUploadFile := BFTUploadCmd.String("file", "", "the file to upload")
	BFTUploadTotal := BFTUploadCmd.Int("total", 4, "the total number of shares to upload")
	BFTRetrieveFile := BFTRetrieveCmd.String("file", "", "the file to retrieve")

	switch os.Args[1] {
	case "upload":
		err := uploadCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing upload command")
			return
		}
		//fmt.Println(*uploadFile)
		cli.Upload(*uploadFile)
	case "retrieve":
		err := retrieveCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing retrieve command")
			return
		}
		//fmt.Println(*retrieveFile)
		//fmt.Println(*retrieveVersion)
		cli.Retrieve(*retrieveFile, *retrieveVersion)
	case "update":
		err := updateCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing update command")
			return
		}
		//fmt.Println(*updateFile)
		//fmt.Println(*updateBase)
		//fmt.Println(*updateVersion)
		cli.Update(*updateFile, *updateBase, *updateVersion)
	case "merge":
		err := mergeCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing merge command")
			return
		}
		//fmt.Println(base)
		//fmt.Println(mergeFile)
		//fmt.Println(base[0])
		//fmt.Println(mergeVersion)
		cli.Update(mergeFile, base[0], mergeVersion)
	case "BFTUpload":
		err := BFTUploadCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing BFTUpload command")
			return
		}
		//fmt.Println(*BFTUploadFile)
		//fmt.Println(*BFTUploadTotal)
		cli.BFTUpload(*BFTUploadFile, *BFTUploadTotal)
	case "BFTRetrieve":
		err := BFTRetrieveCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing BFTRetrieve command")
			return
		}
		//fmt.Println(*BFTRetrieveFile)
		cli.BFTRetrieve(*BFTRetrieveFile)
	default:
		fmt.Println("Invalid command")
	}
}
