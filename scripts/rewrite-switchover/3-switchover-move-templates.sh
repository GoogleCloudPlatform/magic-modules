#!/bin/bash
for subitem in mmv1/templates/terraform/*; do
    for filename in $subitem/*; do
        if [[ $filename != *"*"* ]]; then
            if [[ $filename == *"/go"* ]]; then
                # echo "found go folder $filename"
                    for gofile in $filename/*; do
                        if [[ $gofile == *".tmpl"* ]]; then
                            echo "renaming $gofile to ${gofile/\/go/}"
                            mv  $gofile ${gofile/\/go/}
                        fi
                    done
            fi
        fi
    done
done


for subitem in mmv1/third_party/terraform/*; do
    for filename in $subitem/*; do
        if [[ $filename != *"*"* ]]; then
            if [[ $filename == *"/go"* ]]; then
                echo "found go folder $filename"
                    for gofile in $filename/*; do
                        if [[ $gofile == *".tmpl"* ]]; then
                            echo "renaming $gofile to ${gofile/\/go/}"
                            mv  $gofile ${gofile/\/go/}
                        fi
                    done
            fi
        fi
    done
done