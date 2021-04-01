#!/bin/bash

# Change to the directory with our code that we plan to work from 
# NOTE: DO NOT INCLUDE THIS LINE IF YOU ARE USING GO MODULES!
#cd "$GOPATH/src/lenslocked.com"

echo "==== Releasing lenslocked.com ===="
echo " Deleting the local binary if it exists (so it isn't uploaded)..." 
rm lenslocked.com
echo " Done!"

echo " Deleting existing code..."
ssh root@46.101.31.238 "rm -rf /root/go/src/lenslocked.com" 
echo " Code deleted successfully!"

echo " Uploading code..."
# The \ at the end of the line tells bash that our 
# command isn't done and wraps to the next line. 
rsync -avr --exclude '.git/*' --exclude 'tmp/*' \
  --exclude 'images/*' ./ \
  root@46.101.31.238:/root/go/src/lenslocked.com/
echo "  Code uploaded successfully!"

echo " Building the code on remote server..." 
ssh root@46.101.31.238 'export GOPATH=/root/go; \
  cd $GOPATH/src/lenslocked.com; \
  /usr/bin/go build -o /root/app/server *.go'
echo "  Code built successfully!"

echo "  Moving assets..."
ssh root@46.101.31.238 "cd /root/app; \
cp -R /root/go/src/lenslocked.com/assets ." 
echo " Assets moved successfully!"

echo "  Moving views..."
ssh root@46.101.31.238 "cd /root/app; \
cp -R /root/go/src/lenslocked.com/views ." 
echo " Views moved successfully!"

echo " Moving Caddyfile..."
ssh root@46.101.31.238 "cp /root/go/src/lenslocked.com/Caddyfile /etc/caddy/Caddyfile"
echo " Caddyfile moved successfully!"

echo " Restarting the server..."
ssh root@46.101.31.238 "sudo service lenslocked.com restart" 
echo " Server restarted successfully!"

echo " Restarting Caddy server..."
ssh root@46.101.31.238 "sudo service caddy restart" 
echo " Caddy restarted successfully!"

echo "==== Done releasing lenslocked.com ===="