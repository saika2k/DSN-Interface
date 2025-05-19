package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"flag"
	"strings"
)

type CLI struct {}

func (cli *CLI) Upload(file string) {
	dir, _ := os.Getwd()
	referencePath := filepath.Join(dir, "reference",file)
	fmt.Println(referencePath)
    _, err := os.Stat(referencePath) 

	if err != nil {
		//reference file not found, import it into FileDAG
		cmd := exec.Command("./lotus", "client", "import", file)
		out, _ := cmd.Output()
		parts := strings.Split(string(out), "Root ")
		if len(parts) < 2 { 
			fmt.Println("Error importing file")
			return
		}
		cid := strings.TrimSpace(parts[1])
		fmt.Println(cid)
		cmd2 := exec.Command("./lotus", "client", "deal", cid, "t01000", "0.026", "518400")
		out2, _ := cmd2.Output()
		fmt.Println(string(out2))
	} else {
		//reference file found, upload it to FileDAG
		fmt.Println("the file ", file, "already exist in FileDAG. If you want to update it please use the update/fork/merge command." )
	}
}

func (cli *CLI) Run() {
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)

	uploadFile := uploadCmd.String("file", "", "the file to upload")

	switch os.Args[1] {
	case "upload":
		err := uploadCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing upload command")
			return
		}
		fmt.Println(*uploadFile)
		cli.Upload(*uploadFile)
	default:
		fmt.Println("Invalid command")
	}
}