#!/usr/bin/env bash


# 编译Vue.js项目
#cd html && rm -rf dist && npm run build && cd ..
# 编译嵌入式资源
cp -rf html/dist/* cmd/matrix/html && cd cmd/matrix && cp -rf json/* html && rm -rf bindata.go && go-bindata ./html/... && go build -tags release