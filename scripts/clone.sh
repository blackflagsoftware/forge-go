#!/bin/bash

# arg[1] = repo/full_name/project name
if [ "$#" -ne 1 ]; then
    echo "Need new package folder structure (i.e: github.com/<name>/<project>";
    exit 1;
fi

if [ "$GOPATH" == "" ]; then
    echo "GOPATH is not set"
    exit 1;
fi

path=$1
full_path=$GOPATH/src/$path

# get the project name and make needed upper and env names
name=$(basename $1)
upper_name=$(echo "$name" | tr '[:lower:]' '[:upper:]')
upper_name_env="${upper_name//-/_}"

# check for directory
if [ -e $full_path ]; then
    echo "Dirctory already exists, exiting!"
    exit 1
fi

# copy files
mkdir -p $full_path

# move into the package dir
cd ../pkg
tar -xvzf forge-go-base.tar.gz -C $full_path

# go into new directory and change all the references
cd $full_path

# remove publish script
rm publish.sh

encoded=$(echo $path | sed 's;/;\\/;g')
encoded=$(echo $encoded | sed 's;\.;\\.;g')

# rename proto file
mv ./pkg/proto/proto.proto ./pkg/proto/$name.proto

# replace config
sed -i '' "s/forge-go-base/$name/g" ./config/config.go
sed -i '' "s/FORGE_GO_BASE/$upper_name_env/g" ./config/config.go

# replace rest/main.go
sed -i '' "s/FORGE_GO_BASE/$upper_name/g" ./cmd/rest/main.go

# mac
find . -type f -print0 -exec sed -i '' "s/github\.com\/blackflagsoftware\/forge-go-base/$encoded/g" {} +
# linux or gnu sed
# find . -type f -print0 -exec sed -i "s/github\.com\/blackflagsoftware\/forge-go-base/$encoded/g" {} +
echo -e "\nProject '$name' cloned..."