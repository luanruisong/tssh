#!/bin/bash


workdir=$(cd $(dirname $0); pwd)
target="$workdir/other/date/`date +%Y-%m-%d`.go"
if [ ! -d "target"]; then
 echo "package other\n\n" >> "$target"
fi
echo "func FixData_`date '+%Y%m%d%H%M%S'`()(string,error){}\n\n" >> "$target"
cd $workdir
git add $target
git commit -am "add func at `date '+%Y-%m-%d %H:%M:%S'`"
git push
