# build gossh and passtool 
# support architecture: amd64 386 
# suporrt os: linux and windows
# authur: andes 
# email: email.tata@qq.com

#!/bin/bash

workhome=$(cd $(dirname $0) && pwd)
binpath=${workhome}/bin

oss=(linux windows)
arches=(amd64 386)
target=(gossh passtool)

for arch in ${arches[@]};do
	for os in ${oss[@]};do
		for t in ${target[@]};do
			cmd="CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build"
			echo "${cmd} ${t}.go"
			eval ${cmd} ${t}.go
			
			if [[ ! -d ${binpath}/${arch}/${os} ]];then
					mkdir -p ${binpath}/${arch}/${os}
			fi

			if [[ "${os}" == "windows"  ]];then
					mv ${t}.exe  ${binpath}/${arch}/${os}
			else
					mv ${t}  ${binpath}/${arch}/${os}
			fi
		done
	done
done

