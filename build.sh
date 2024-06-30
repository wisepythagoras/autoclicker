#!/bin/bash

cd cmd/cli
go build
mv cli ../../autoclicker-cli

cd ../..

cd cmd/gui
go build
mv gui ../../autoclicker-gui
