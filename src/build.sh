#!/usr/bin/env bash

if [ -z "${VERSION}" ]; then
  VERSION=devel
fi

rm -rf build
mkdir -p build

os_archs=("darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")
for os_arch in "${os_archs[@]}"; do
  IFS="/" read -r os arch <<< "$os_arch"
  output_name="mc-quick-${os}-${arch}"
  echo "Building for $os/$arch..."
  GOOS=$os GOARCH=$arch go build -o "build/$output_name" \
    -ldflags="-X main.Version=$VERSION"
done

