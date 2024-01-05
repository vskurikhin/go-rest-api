#!/bin/zsh

for d in $(ls)
do
  if [[ $d == issue* ]]; then
    cd $d
    echo "Update deps in ${d}..."
    go mod tidy
    go get -t -u ./...
    cd ..
  fi
done
