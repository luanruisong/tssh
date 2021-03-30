#!/bin/bash

tag=/usr/local/bin/tssh
down=https://github.com/luanruisong/tssh/releases/download/
v=$0
suffex=$1

sudo wget -O $tag $down$v/tssh-$suffex

sudo chmod +x $tag