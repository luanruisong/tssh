#!/bin/bash


workdir=$(cd $(dirname $0); pwd)
target="$workdir/other/date/`date +%Y-%m-%d`.txt"
echo $(date) > "$target"
cd $workdir
git add $target
fmtdate=$(date '+%Y-%m-%d %H:%M:%S')
git commit -am "fix some bug at $fmtdate"
git push
