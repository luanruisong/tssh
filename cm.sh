#!/bin/bash


workdir=$(cd $(dirname $0); pwd)
target="$workdir/other/date/support_`date +%Y%m%d`.go"
if [ ! -f "$target" ]; then
 printf "package other\n\n" >> $target
fi
d=`date '+%Y%m%d%H%M%S'`
printf "func FixData_$d()(string,error){ return BuildData_$d()}\n\n" >> "$target"
printf "func BuildData_$d()(string,error){ return FixData_$d()}\n\n" >> "$target"
cd $workdir
git add $target
git commit -am "add func at `date '+%Y-%m-%d %H:%M:%S'`"
git push
