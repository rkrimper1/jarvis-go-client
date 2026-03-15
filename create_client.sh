#!/bin/sh
#    Proprietary and Confidential
#
#    This source code is the property of:
#
#     Robert Krimper (c) 2026
#
#     https://www.krimper.com
#
#    Author:         Robert Krimper, https://www.linkedin.com/in/robert-krimper
#    Modified by:    
#    Module:         create_client.sh
#    Description:    Creates a zip file of a working gRPC client for GO.
#                    This will be labeled with the current version of jarvis code.
#
#

ZIP_DIR=jarvis-go-client

# Copy all the required files and directories in order to run the Go Client
rm -rf jarvis-client-version_number-go.zip
rm -rf $ZIP_DIR
mkdir -p $ZIP_DIR

# Copy the client README.md
cp client_README.md $ZIP_DIR/README.md

# Copy essential client single files
cp go.mod $ZIP_DIR/go.mod
cp go.sum $ZIP_DIR/go.sum
cp main.go $ZIP_DIR/main.go

cp -R api $ZIP_DIR

# zip the client up
 zip -r jarvis-client-version_number-go.zip jarvis-go-client