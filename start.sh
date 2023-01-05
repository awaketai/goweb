cd build

make build

cd ..
# start http server
# ADDRESS=:8888 ./gw app start 
# ./gw app start --address=:8082 --daemon=false
./gw app start
# generate swagger file
# ./gw swagger  gen
# ./gw provider new