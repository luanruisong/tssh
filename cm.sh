#!/bin/bash


workdir=$(cd $(dirname $0); pwd)
target="$workdir/other/date/support_`date +%Y%m%d%`.go"
if [ ! -f "$target" ]; then
 echo "package other\n\n" >> $target
fi
echo "func FixData_`date '+%Y%m%d%H%M%S'`()(string,error){}\n\n" >> "$target"
cd $workdir
git add $target
git commit -am "add func at `date '+%Y-%m-%d %H:%M:%S'`"
git push
