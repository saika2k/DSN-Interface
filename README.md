# DSN-Interface


## Discription
DSN-Interface is an auxiliary program used to provide automated solutions for file operations of clients of the two decentralized storage networks FileDAG and BFT-DSN.
DSN-Interface currently supports automated execution of file upload, update, fork, merge and retrieve operations of FileDAG as well as file upload and retrieval operations of BFT-DSN.

## Environment Requirement

go 1.18

## Install FileDAG and BFT-DSN

Users need to first download FileDAG and BFT-DSN. These two projects can be clones on Github.

Clone project FileDAG:
```sh
https://github.com/FileDAG/lotus.git
cd lotus
git checkout filedag
```
This project is the implementation of the paper "FileDAG: A Multi-Version Decentralized Storage Network Built on DAG-Based Blockchain" (IEEE TC 2023), you can get this paper on http://ieeexplore.ieee.org/abstract/document/10159425

Clone project BFT-DSN:
```sh
https://github.com/FileDAG/lotus.git
cd lotus
git checkout bft-dsn
```
This project is the implementation of the paper "BFT-DSN: A Byzantine Fault-Tolerant Decentralized Storage Network" (IEEE TC 2024), you can get this paper on https://ieeexplore.ieee.org/abstract/document/10436433

Users can follow the README.md of these two projects to install these two projects.

## Install DSN-Interface

DSN-Interface depends on bsdiff, bspatch, zfec and zunfec. Users can install them with the following command
```sh
sudo apt update && sudo apt install -y bsdiff python3-pip && pip3 install zfec
```

Note: when installing zfec, you may see this warning:
```sh
WARNING: The scripts zfec and zunfec are installed in '/home/yourname/.local/bin' which is not on PATH.
```
When this warning appear, you need to add the zfec installation path to the environment variable PATH by using the following command, otherwise an error will occur when DSN-Interface is running because it cannot find the zfec and zunfec commands.
```sh
echo 'export PATH=$PATH:~/.local/bin' >> ~/.bashrc
source ~/.bashrc
```
Finally, we can clone DSN-Interface project and compile it:
```sh
git clone https://github.com/saika2k/DSN-Interface.git
cd DSN-Interface
go install
```
The above command will generate an executable file DSN-Interface in the GOBIN path.

## Usage

Rmember to copy DSN-Interface to the path of FileDAG or BFT-DSN.

To know how to use DSN-Interface just run the following command for help:
```bash
./DSN-Interface
```

In addition, we provide some auxiliary test scripts and test files in the test-data folder of DSN-Interface.
The script test-run.sh is used to start FileDAG and BFT-DSN locally and build a local test network. This network supports uploading files not exceeding 8MB.

The script DSN-Interface-test1.sh is used to test the correctness of operaters in FileDAG. We use the file1, file2, file3, file4 and file5 to run the test. 

1) This test first upload and retrieve file1, the first version of the file. 
2) Then we update file1, the new version is file2, and retrieve file2. 
3) After that we update file2, the new version is file3, and retrieve file3. 
4) After that we test the fork command which we based on file2 and update the file with file4, and retrieve file4. 
5) Finally, we test the merge command which we based on file2, file3 and file4 and update these file with file5, and retrieve file5. 

The results of this test are shown in the comments in this file.

The script DSN-Interface-test2.sh is used to test the correctness of operaters in BFT-DSN. We use the file6 to run the test. 

1) This test first upload file6, we set the client to upload 4 EC shares.
2) Then we retrieve file6, where the client retrieve each share of the file and try to recover the file.

The results of this test are shown in the comments in this file.

Users can use these scripts to become familiar with the various commands of DSN-Interface and use their own data to play with DSN-Interface.


