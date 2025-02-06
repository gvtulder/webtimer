#!/bin/bash

rm -rf public
mkdir -p public

cd server/
go build .
cd ..

cp html/style.css public/

timestamp=$(date +%s)
sed -E 's/(style.css|main.js|icon.png)/\1?'$timestamp'/' < html/index.html > public/index.html
sed --in-place -E 's/..\/dist/./' public/index.html
sed --in-place -E "s/webtimer\.startApp\(.+\)/webtimer.startApp('ws')/" public/index.html

cp dist/main.js public/

rsync -avz public/ 10.10.10.11:/var/www/html/webtimer/
rsync -avz server/webtimer 10.10.10.11:/var/www/html/webtimer/webtimer
