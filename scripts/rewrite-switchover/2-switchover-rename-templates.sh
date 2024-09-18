#!/bin/bash
for subitem in mmv1/templates/terraform/*; do
    for filename in $subitem/*; do
        if [[ $filename != *"*"* ]]; then
            if [[ $filename == *".erb"* ]]; then
                echo "renaming $filename to ${filename/.erb/.tmpl}"
                mv $filename ${filename/.erb/.tmpl}
            fi
        fi
    done
done

for subitem in mmv1/templates/terraform/iam/*; do
    for filename in $subitem/*; do
        if [[ $filename != *"*"* ]]; then
            if [[ $filename == *".erb"* ]]; then
                echo "renaming $filename to ${filename/.erb/.tmpl}"
                mv $filename ${filename/.erb/.tmpl}
            fi
        fi
    done
done

for subitem in mmv1/third_party/terraform/*; do
    for filename in $subitem/*; do
        if [[ $filename != *"*"* ]]; then
            if [[ $filename == *".erb"* ]]; then
                echo "renaming $filename to ${filename/.erb/.tmpl}"
                mv $filename ${filename/.erb/.tmpl}
            fi
        fi
    done
done

for subitem in mmv1/third_party/terraform/services/*; do
    for filename in $subitem/*; do
        if [[ $filename != *"*"* ]]; then
            if [[ $filename == *".erb"* ]]; then
                echo "renaming $filename to ${filename/.erb/.tmpl}"
                mv $filename ${filename/.erb/.tmpl}
            fi
        fi
    done
done