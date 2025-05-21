#!/bin/bash

./DSN-Interface BFTUpload -file file6 -total 4
sleep 600
#Note: this command upload 4 shares to storage miners, the CID of each share are stored in the reference file
#Note: we need to wait the storage miner to handle the deal

./DSN-Interface BFTRetrieve -file file6
sleep 10
./lotus client import file6
./lotus client import recover_file6
#Note: the CID of recover_file6 should be the same as file6 (bafykbzacebae567lb3olmwex5halthu4oatbszinp7cmpf6jsha5sc4mvo2p2)


########################################################################################
#below is the content of the reference file of file6
#bafykbzaceazlqdr37llxzp44nvcjwo4sshfd6os7gqygsukufyfliessqgkxu -1
#bafykbzacecz6st37u37g6kabvysy2vvb6a7qrj6yljpikuljvbwfwhbrmacjm -1
#bafykbzacebu5tr6fmj5yfmakeb7dwq3deam3cwssgyuysiwwdgxo25kkb7qds -1
#bafykbzacecs3s3thza4ubyfcmo55vkgxbd4bpun56dkyuvrsyynkoo5pvz53a -1
########################################################################################

