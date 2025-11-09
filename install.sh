#!/bin/bash

echo "Building the Binary"
go build -o agai .

echo "Installing ..."
sudo mv agai /usr/bin/agai

echo "instalation Done"
agai -h