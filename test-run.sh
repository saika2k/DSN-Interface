#!/bin/bash

rm -rf devgen.car
rm -rf localnet.json
rm -rf ~/.genesis-sectors
rm -rf ~/.lotus
rm -rf ~/.lotusminer
rm -rf DB
rm -rf SST
rm -rf client_upload
rm -rf client_download
rm -rf database

mkdir DB
mkdir SST
mkdir client_upload
mkdir client_download
touch acc
touch leaf
./lotus fetch-params 8MiB
./lotus-seed pre-seal --sector-size 8MiB --num-sectors 2
./lotus-seed genesis new localnet.json
./lotus-seed genesis add-miner localnet.json ~/.genesis-sectors/pre-seal-t01000.json
nohup ./lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false > lotus.log 2>&1 &

ps -ef | grep lotus

echo "sleep 30s" && sleep 30s
./lotus wallet import --as-default ~/.genesis-sectors/pre-seal-t01000.key
./lotus-miner init --genesis-miner --actor=t01000 --sector-size=8MiB --pre-sealed-sectors=~/.genesis-sectors --pre-sealed-metadata=~/.genesis-sectors/pre-seal-t01000.json --nosync
#tmux new-session -s "lotus-miner" -d "./lotus-miner run --nosync"
nohup ./lotus-miner run --nosync > miner.log 2>&1 &

dd if=/dev/urandom of="database" bs=8M count=300

./lotus client import 4.5M
./lotus client import 5.0M
./lotus client import 5.5M
./lotus client import 6.0M
./lotus client import 6.5M
./lotus client import 7.0M
./lotus client import 7.5M

#sleep 10

#./lotus-miner net listen
#./lotus net listen

#./main bls.txt






