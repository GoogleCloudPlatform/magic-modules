#!/bin/bash

folder="/products/*"
input="${folder}"
name_line=false
line_cache=()
resource_name=""


#cd products/
#find . -maxdepth 1 -mindepth 1 -type d -printf '%f\n'

mapfile -d $'\0' array < <(find products/. -maxdepth 1 -mindepth 1 -type d -printf '%f\n')

IFS=$'\n' sorted=($(sort <<<"${array[*]}")); unset IFS
printf "%s\n" "${sorted[@]}"

#sorted=( "containeranalysis" "containerattached" )

for i in "${sorted[@]}"
do
	if [[ ${i:0:1} == "e" || ${i:0:1} == "f" ]]; then
		echo "Calling sh convert.sh ${i}"
		sh convert.sh "$i"
		line_cache+="bundle exec compiler -p \"products/${i}\" -v \"beta\" -e terraform -o \"$GOPATH/src/github.com/googlerjk/terraform-provider-google-beta\"\n"
		
	fi
#	
	#bundle exec compiler -p "products/bigquery" -v "beta" -e terraform -o "$GOPATH/src/github.com/googlerjk/terraform-provider-google-beta"
	# or do whatever with individual element of the array
done

for i in "${line_cache[@]}"
do
	echo -e $i
done
