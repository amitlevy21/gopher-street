#! /usr/bin/env bash

scripts_dir=$(dirname "$0")
scripts_dir=$(cd $scripts_dir && pwd)
project_dir=$(dirname "$scripts_dir")

docker run -p 27017:27017 -v /tmp/gopher-street-db:/data/db mongo
