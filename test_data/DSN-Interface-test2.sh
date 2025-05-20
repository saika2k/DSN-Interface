#!/bin/bash

./DSN-Interface BFTUpload -file file6 -total 4
sleep 600
#Note: the CID of file1: bafykbzacebylwxfbdyazodwqqln5modpp6w32jzbn7p2pzh6qdy4b6zpcxyjc
#Note: we need to wait the storage miner to handle the deal

./DSN-Interface BFTRetrieve -file file6
sleep 10
./lotus client import file6
./lotus client import recover_file6
#Note: the CID of file1_v1 should be the same as file1


########################################################################################
#below is the content of the reference file of file1
#bafykbzacebylwxfbdyazodwqqln5modpp6w32jzbn7p2pzh6qdy4b6zpcxyjc -1
#bafykbzaceckdlpbjc4bf4fhyqb3m32p3uywebirqt6fewjweulyw4heuxuibc 1
#bafykbzacedo2l4l4nu33dtnfnxqdi5gojpwgmmnu7il5zm2he4gipqvbdn2go 2
#bafykbzaceb6djch625f4plidvoq2ahgh6mk4ixf2uurfmbwmfg3ssro5io4eo 3
#bafykbzacedo3rbh6dtntzkftouwli5y3l4lwsjy7qpnyvdyrt7joowha75wpg 2
########################################################################################

