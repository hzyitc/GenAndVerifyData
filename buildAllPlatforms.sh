#!/bin/bash

OUTPUT="bin/"
PROGRAM="GenAndVerifyData"
LDFLAGS="-s -w"

mkdir -p "${OUTPUT}"
rm -f "${OUTPUT}/${PROGRAM}-"*

platforms=(
	linux/386
	linux/amd64
	linux/arm
	linux/arm64
	linux/mips/softfloat
	linux/mips64
	linux/mips64le
	linux/mipsle/softfloat
	windows/386
	windows/amd64
	windows/arm
)
# platforms=($(go tool dist list))

for i in "${platforms[@]}"; do
	os="$(echo "${i}" | awk -F/ '{print $1}')"
	arch="$(echo "${i}" | awk -F/ '{print $2}')"
	mips="$(echo "${i}" | awk -F/ '{print $3}')"

	[ "${os}" == "windows" ] && ext="exe"

	filename="${OUTPUT}/${PROGRAM}-${os}-${arch}${ext:+.$ext}"
	echo "build ${filename} for ${i}"
	CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" GOMIPS="${mips}" \
		go build -trimpath -ldflags "${LDFLAGS}" -o "${filename}"
done
