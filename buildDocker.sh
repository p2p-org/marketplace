#!/usr/bin/env bash

docker_img_mp_name=dwh_marketplace

cur_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

cd $cur_path

docker build -t $docker_img_mp_name .

