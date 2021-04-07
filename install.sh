#!/bin/bash

tag=/usr/local/bin/tssh
cpu_brand=$(sysctl machdep.cpu |grep brand_string)
#down=https://github.com/luanruisong/tssh/releases/download/
down=https://github.91chifun.workers.dev/https://github.com//luanruisong/tssh/releases/download/
version=$(wget -qO- -t1 -T2 "https://api.github.com/repos/luanruisong/tssh/releases/latsest" | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g')
if [ -z "$version" ]; then
    echo "can not get latest release"
else
    suffex=intel
    result=$(echo $cpu_brand | grep "Apple M1")
    if [[ "$result" != "" ]]; then
         suffex=appleSilicon
    fi
    sudo wget -O $tag $down$version/tssh-$suffex
    sudo chmod +x $tag
fi
