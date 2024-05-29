#!/bin/bash

# Function to display the usage of the script
function usage() {
   echo "Usage: $0 [SOURCE_DIR] [force] [--skip-version-inc]"
   echo "SOURCE_DIR         Optional. The source directory to use."
   echo "force              Optional. Pass 'force' to force the build."
   echo "--skip-version-inc Optional. Flag to skip version increment."
   exit 1
}

# Check if the first argument is -h or --help to display the usage and exit
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
   usage
fi

# Define the data directory path for database
DATA_DIR="./data"

# Check if the data directory exists
if [ ! -d "$DATA_DIR" ]; then
    echo "Data directory does not exist. Creating..."
    mkdir -p "$DATA_DIR"
else
    echo "Data directory already exists. Using existing database..."
fi

./build.sh

containers=$(docker container ls --format '{{.Names}}' | grep 'vsys.')

#stop all containers
docker container stop $containers
docker container wait $containers
docker container rm $containers

docker image rm $(docker image ls | grep -E '^vsys\.' | awk '{print $3}')
docker image rm $(docker images --format "{{.ID}}"  -f "dangling=true")|tr -d '\n'

docker-compose --env-file env -f vsys-compose.yml -p vsys-rest build --no-cache 
docker-compose --env-file env -f vsys-compose.yml -p vsys-rest up &

docker container prune -f
docker image rm $(docker images --format "{{.ID}}"  -f "dangling=true")|tr -d '\n'
