#!/bin/bash
if [ $# -ne 0 ]; then
  npx zenn new:article --slug $1
else
  npx zenn new:article
fi
