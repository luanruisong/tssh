#!/bin/bash

tag=/usr/local/bin/tssh
down=https://github.com/luanruisong/tssh/releases/download/
v=$1
suffex=$2

sudo wget -O $tag $down$v/tssh-$2

sudo chmod +x $tag