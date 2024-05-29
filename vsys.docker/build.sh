#!/bin/bash

goProjects=("../vsys.dbhelper" "../vsys.rest")

for i in "${goProjects[@]}"
do
   echo --------------------$i------------------ | tr [a-z] [A-Z]
   cd "$i"
   
   go clean
   go mod tidy
   CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -buildvcs=false -o ${i##*/} .

done