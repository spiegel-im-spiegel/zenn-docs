#!/bin/bash
y=$(date '+%Y')
m=$(date '+%m')
d=$(date '+%d')
if [ $# -ne 0 ]; then
  npx zenn new:article --slug "$y$m$d-$1"
else
  npx zenn new:article
fi
