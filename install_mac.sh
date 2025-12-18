#!/bin/bash

echo "Building the Binary"
go build -o agai .

echo "Creating bin folder"
sudo mkdir /usr/local/bin

echo "Installing ..."
sudo mv agai /usr/local/bin/agai

echo "instalation Done"
agai -h