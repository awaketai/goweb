cd build

make build

cd ..
# start http server
./gw app start 
# generate swagger file
# ./gw swagger  gen
# ./gw provider new