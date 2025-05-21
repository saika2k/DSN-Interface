#!/bin/bash

rm -rf devgen.car
rm -rf localnet.json
rm -rf ~/.genesis-sectors
rm -rf ~/.lotus
rm -rf ~/.lotusminer
rm -rf database
rm -rf reference

mkdir reference
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


#sleep 10

#./lotus-miner net listen
#./lotus net listen

#./main bls.txt






