#!/bin/bash

if [ $# -ne 1 ]; then
    echo "Usage: $0 <input_filename.go>"
    exit 1
fi

input_file="$1"

declare -A platforms

platforms=(
    ["darwin"]="amd64 arm64"
    ["freebsd"]="amd64 arm arm64"
    ["linux"]="amd64 arm arm64"
    ["windows"]="amd64 arm arm64"
)

output_dir="builds"
mkdir -p "$output_dir"

# Function to compile for a specific platform
compile_platform() {
    goos="$1"
    goarch="$2"
    output_name="$output_dir/$goos-$goarch"
    echo "Compiling for $goos/$goarch: $input_file"
    env GOOS="$goos" GOARCH="$goarch" go build -o "$output_name" "$input_file"
}

# Iterate over platforms and compile asynchronously
for goos in "${!platforms[@]}"; do
    for goarch in ${platforms["$goos"]}; do
        compile_platform "$goos" "$goarch" &
    done
done

# Wait for all background tasks to finish
wait

echo "Cross-compilation complete. Output files are in the '$output_dir' directory."
