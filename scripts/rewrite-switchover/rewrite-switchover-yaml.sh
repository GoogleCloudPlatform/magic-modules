#!/bin/bash
for service in mmv1/products/*; do
    echo "$service"
    for filename in $service/*; do
        # echo "$filename"
        if [[ $filename != *"go_"* ]]; then
            echo "removing $filename"
            rm -rf $filename
        fi
    done

    for filename in $service/*; do
        if [[ $filename == *"go_"* ]]; then
            echo "renaming $filename to ${filename/go_/}"
            mv $filename "${filename/go_/}"
        fi
    done
done