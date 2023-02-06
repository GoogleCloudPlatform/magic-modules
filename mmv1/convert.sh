#!/bin/bash

#!/bin/bash
folder="products/$1/"
input="${folder}terraform.yaml"
declare -A tf_content_dict
resource_name=""
line_cache=()
	
function get_terraform_section {

	while IFS= read -r line
	do
		if [[ "$line" = *": !ruby/object:Overrides::Terraform::ResourceOverride" ]]; then						

			len=${#line_cache[@]}
			if [ $len -gt 0 ] && [ "${resource_name}" != "" ]; then
				tf_content_dict["$resource_name"]+="${line_cache[*]}"
				# echo "${tf_content_dict[$resource_name]}"
			fi
			
			searchstring=": !ruby/object:Overrides::Terraform::ResourceOverride"			
			resource_name=${line%$searchstring}
			eval resource_name=$(tr -d ' ' <<< "$resource_name")
			line_cache=()
		else
			line_cache+=("$line")
		fi
	done < "$input"
	
	tf_content_dict["$resource_name"]+="${line_cache[*]}"
}

#get_terraform_section
copyright_cache=()

copyright_cache+=("# Copyright 2020 Google Inc.")
copyright_cache+=("# Licensed under the Apache License, Version 2.0 (the \"License\");")
copyright_cache+=("# you may not use this file except in compliance with the License.")
copyright_cache+=("# You may obtain a copy of the License at")
copyright_cache+=("#")
copyright_cache+=("#     http://www.apache.org/licenses/LICENSE-2.0")
copyright_cache+=("#")
copyright_cache+=("# Unless required by applicable law or agreed to in writing, software")
copyright_cache+=("# distributed under the License is distributed on an \"AS IS\" BASIS,")
copyright_cache+=("# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.")
copyright_cache+=("# See the License for the specific language governing permissions and")
copyright_cache+=("# limitations under the License.")

len_copyright=${#copyright_cache[@]}

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
		i=1		
		while [ $i -le $len_copyright ]
		do  
			line_copyright=${copyright_cache[$i]} 
			if [ $i -eq 1 ]; then          
				echo "${line_copyright}" > $filename  
			else    
				echo "${line_copyright}" >> $filename
			fi
			i=$((i+1))
		done
		
		echo "--- !ruby/object:Api::Resource" >> $filename
		i=1
		while [ $i -le $len ]
		do
			line_pending=${line_cache[$i]:2}                
			echo "${line_pending}" >> $filename
			i=$((i+1))
		done

		#           eval key_name=$(tr -d ' ' <<< "$resource_name")
		#           echo "${tf_content_dict[$key_name]}" >> $filename
            
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
	filename="${folder}${resource_name}.yaml"
	filename="${filename//\'}"
	echo "${filename}"
	i=1		
	while [ $i -le $len_copyright ]
	do  
		line_copyright=${copyright_cache[$i]} 
		if [ $i -eq 1 ]; then          
			echo "${line_copyright}" > $filename  
		else    
			echo "${line_copyright}" >> $filename
		fi
		i=$((i+1))
	done
		
	echo "--- !ruby/object:Api::Resource" >> $filename
	i=1
	while [ $i -le $((len-1)) ]
	do
		line_pending=${line_cache[$i]:2}
		echo "${line_pending}" >> $filename
		i=$((i+1))
	done
	
#	eval resource_name=$(tr -d ' ' <<< "$resource_name")	
#	echo "${tf_content_dict[$resource_name]}" >> $filename
	
#	echo "added ${resource_name} ~~ ${tf_content_dict[$resource_name]}"
	
#	for key in "${!tf_content_dict[@]}"; do
#		echo "${#key} - ${key}"
#		echo "${#resource_name} - ${resource_name}"
#		
#		if [ "$key" == "$resource" ]; then
#			echo "matched ${resource_name} - ${key}"
#		fi
#	done

fi


echo "api.yaml renamed" 
$(mv ${folder}"/api.yaml" "/usr/local/google/home/juliuskelly/go/src/github.com/googlerjk/temp/${1}_api.yaml")

#echo "terraform.yaml renamed" 
#$(mv ${folder}"/terraform.yaml" "/usr/local/google/home/juliuskelly/go/src/github.com/googlerjk/temp/${1}_terraform.yaml")

