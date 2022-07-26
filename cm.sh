#!/bin/bash


workdir=$(cd $(dirname $0); pwd)
target="$workdir/other/date/`date +%Y-%m-%d`.go"
echo "func FixData_`date '+%Y%m%d%H%M%S'`()(string,error){}" >> "$target"
cd $workdir
git add $target
git commit -am "add func at `date '+%Y-%m-%d %H:%M:%S'`"
git push
