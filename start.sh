cd build

make build

cd ..

#!/bin/bash
function on_ctrl_c() {
    echo "Ctrl+C detected.Performing some work..."
    if [ -e "./gw" ];then
      echo "gw exec will be deleted."
      rm -f ./gw
      echo "gw exec deleted"
    fi
    exit 0
}

trap "on_ctrl_c" SIGINT
# start http server
# ADDRESS=:8888 ./gw app start 
# ./gw app start --address=:8082 --daemon=false
./gw app start
# generate swagger file
# ./gw swagger  gen
# ./gw provider new