#!/bin/bash -e 
mkdir -p bin
cd bin
go build -v github.com/slspeek/gotube/main
cd ..
mongo e2etest --eval 'db.dropDatabase();'
screen -d -m bin/main -port 8484 -db e2etest $@
cd web
protractor protConf.js
cd ..
kill $(cat gotube.pid)
