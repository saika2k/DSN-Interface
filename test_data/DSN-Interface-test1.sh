#!/bin/bash

./DSN-Interface upload -file file1
sleep 360
#Note: the CID of file1: bafykbzacebylwxfbdyazodwqqln5modpp6w32jzbn7p2pzh6qdy4b6zpcxyjc
#Note: we need to wait the storage miner to handle the deal

./DSN-Interface retrieve -file file1 -version 1
sleep 10
./lotus client import file1_v1
#Note: the CID of file1_v1 should be the same as file1

rm -rf file1_v1
#Note: the update command retrieve the old version, so we delete the file to avoid mistakes

./DSN-Interface update -file file1 -base 1 -new file2
sleep 360
#Note: the CID of the uploaded patch: bafykbzaceckdlpbjc4bf4fhyqb3m32p3uywebirqt6fewjweulyw4heuxuibc (same as the CID in the second line of reference file)

./DSN-Interface retrieve -file file1 -version 2
sleep 10
./lotus client import file2
./lotus client import file1_v2
#Note: the CID of file2 and file1_v2 should be the same: bafykbzaceb7mtjkhn3jxbrm4implkoxanwmp7ijif3le2l2lwqom4zioezbsu

rm -rf file1_v2

./DSN-Interface update -file file1 -base 2 -new file3
sleep 360
#Note: the CID of the uploaded patch: bafykbzacedo2l4l4nu33dtnfnxqdi5gojpwgmmnu7il5zm2he4gipqvbdn2go (same as the CID in the third line of reference file)

./DSN-Interface retrieve -file file1 -version 3
./lotus client import file3
./lotus client import file1_v3
#Note: the CID of file3 and file1_v3 should be the same: bafykbzacecjxono6lk5unwohxmutftxdvh4wzat5apzqo4czhvvphjqytdy4i

rm -rf file1_v3

./DSN-Interface update -file file1 -base 3 -new file4 
sleep 360
#Note: the CID of the uploaded patch: bafykbzaceb6djch625f4plidvoq2ahgh6mk4ixf2uurfmbwmfg3ssro5io4eo (same as the CID in the fourth line of reference file)

./DSN-Interface retrieve -file file1 -version 4
sleep 10
./lotus client import file4
./lotus client import file1_v4
#Note: the CID of file4 and file1_v4 should be the same: bafykbzaced36jsi36wxrvehuch3qp3dinswf3umr3uxam6egjbj2gj6inpmts

rm -rf file1_v4

./DSN-Interface merge -file file1 -base 2,3,4 -new file5
sleep 360
#Note: the CID of the uploaded patch: bafykbzacedo3rbh6dtntzkftouwli5y3l4lwsjy7qpnyvdyrt7joowha75wpg (same as the CID in the fifth line of reference file)

./DSN-Interface retrieve -file file1 -version 5
sleep 10
./lotus client import file5
./lotus client import file1_v5
#Note: the CID of file5 and file1_v5 should be the same: bafykbzacedolni6k7t6yavrp2xisakchepdujd5gy4stfxvylp7djf73ywwwu


########################################################################################
#below is the content of the reference file of file1
#bafykbzacebylwxfbdyazodwqqln5modpp6w32jzbn7p2pzh6qdy4b6zpcxyjc -1
#bafykbzaceckdlpbjc4bf4fhyqb3m32p3uywebirqt6fewjweulyw4heuxuibc 1
#bafykbzacedo2l4l4nu33dtnfnxqdi5gojpwgmmnu7il5zm2he4gipqvbdn2go 2
#bafykbzaceb6djch625f4plidvoq2ahgh6mk4ixf2uurfmbwmfg3ssro5io4eo 3
#bafykbzacedo3rbh6dtntzkftouwli5y3l4lwsjy7qpnyvdyrt7joowha75wpg 2
########################################################################################

