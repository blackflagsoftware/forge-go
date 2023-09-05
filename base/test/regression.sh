#! /bin/bash
# TODO: set up your env vars here

# you can also add a setup script to import data to work with

pwd=$(pwd)
compile_and_run() {
	cd ../cmd/rest && go build
	./rest > /dev/null &
	echo "$!"
}

pid=$(compile_and_run)
sleep 1

cd $pwd
file_path=./your-file-here.json # see README.md for other directory/multiple file options
regression -testPathFile=$file_path -environment=dev -printDebug=true

kill $pid

# make sure you delete any data you imported and started with