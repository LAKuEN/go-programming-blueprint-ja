#!/bin/zsh
echo "build domainfinder..."
go build -o domainfinder

mkdir lib

echo "build synonyms..."
cd ../synonyms
go build -o ../domainfinder/lib/synonyms

echo "build available..."
cd ../available
go build -o ../domainfinder/lib/available

echo "build sprinkle..."
cd ../sprinkle
go build -o ../domainfinder/lib/sprinkle

echo "build domainify..."
cd ../domainify
go build -o ../domainfinder/lib/domainify

echo "build succeed"
