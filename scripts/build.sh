#!/bin/bash
set -e

echo -e "\033[36mmesa build tool\033[0m"
echo -e "select build type (1-5)"
echo "1. stable"
echo "2. beta"
echo "3. nightly"
echo "4. debug"
echo "5. dev"
echo
read -p "Enter your choice (1-5): " choice

case $choice in
    1)
        export MESA_BUILD_TYPE=release
        echo -e "\033[32mBuilding Stable...\033[0m"
        ;;
    2)
        export MESA_BUILD_TYPE=beta
        echo -e "\033[33mBuilding Beta...\033[0m"
        ;;
    3)
        export MESA_BUILD_TYPE=nightly
        echo -e "\033[35mBuilding Nightly...\033[0m"
        ;;
    4)
        export MESA_BUILD_TYPE=debug
        echo -e "\033[31mBuilding Debug...\033[0m"
        ;;
    5)
        export MESA_BUILD_TYPE=internal
        echo -e "\033[34mBuilding Dev Build...\033[0m"
        ;;
    *)
        echo -e "\033[31mInvalid choice. Exiting.\033[0m"
        exit 1
        ;;
esac

echo

echo "MESA_BUILD_TYPE = $MESA_BUILD_TYPE"

distDir="dist"
mkdir -p "$distDir"

declare -a targets=(
    "linux amd64 mesa-linux-amd64"
    "linux arm64 mesa-linux-arm64"
    "linux arm mesa-linux-armv7 GOARM=7"
    "windows amd64 mesa-windows-amd64.exe"
    "windows arm64 mesa-windows-arm64.exe"
)

for target in "${targets[@]}"; do
    IFS=' ' read -r GOOS GOARCH OUTFILE EXTRA <<< "$target"
    echo -e "\033[36mBuilding $OUTFILE...\033[0m"
    if [[ $EXTRA == GOARM=* ]]; then
        GOARM_VAL=${EXTRA#GOARM=}
        env GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM_VAL go build -o "$distDir/$OUTFILE" ./cmd/mesa/main.go
    else
        env GOOS=$GOOS GOARCH=$GOARCH go build -o "$distDir/$OUTFILE" ./cmd/mesa/main.go
    fi
    if [[ $? -ne 0 ]]; then
        echo -e "\033[31mBuild failed for $OUTFILE\033[0m"
        exit 1
    fi
done

echo

echo -e "\033[36mSHA256 checksums:\033[0m"
for f in "$distDir"/*; do
    sha256sum "$f"
done

unset MESA_BUILD_TYPE 