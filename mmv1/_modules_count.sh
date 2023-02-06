#!/bin/bash
folder="products/$1/"
input="${folder}api.yaml"
name_line=false
line_cache=()
resource_name=""

while IFS= read -r line
do
    if [ $name_line == true ]
    then
        name_line=false
        searchstring="    name: "
        resource_name=${line#*$searchstring}
    fi
    if [ "$line" == "  - !ruby/object:Api::Resource" ]; then
        len=${#line_cache[@]}
       
        if [ $len -gt 0 ]; then
            filename="${folder}${resource_name}.yaml"
            filename="${filename//\'}"
            echo "${filename}"
            echo "--- !ruby/object:Api::Resource" > $filename
            i=1
            while [ $i -le $len ]
            do
                line_pending=${line_cache[$i]:2}
                echo "${line_pending}" >> $filename
                i=$((i+1))
            done
        fi
        line_cache=()
        name_line=true
    fi
    if [ "$line" == "objects:" ]; then
        len=${#line_cache[@]}
        if [ $len -gt 0 ]; then
            filename="${folder}product.yaml"
            echo "${filename}"
            echo ${line_cache[0]} > $filename
            i=1
            while [ $i -le $((len-1)) ]
            do
                line_pending=${line_cache[$i]}
                echo "${line_pending}" >> $filename
                i=$((i+1))
            done
        fi
        resource_start=true
        line_cache=()
    else 
        line_cache+=("$line")
    fi
done < "$input"

len=${#line_cache[@]}
if [ $len -gt 0 ]; then
    filename="${folder}${resource_name,,}.yaml"
    filename="${filename//\'}"
    echo "${filename}"
    echo "--- !ruby/object:Api::Resource" > $filename
    i=1
    while [ $i -le $((len-1)) ]
    do
        line_pending=${line_cache[$i]:2}
        echo "${line_pending}" >> $filename
        i=$((i+1))
    done
fi


new_name='_del_api.yaml'
echo $new_name 
rename=${input/api.yaml/$new_name}
echo $rename
$(mv $input $rename)

