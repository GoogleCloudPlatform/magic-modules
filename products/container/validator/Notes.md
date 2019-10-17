This directory contains the configuration required to build
terraform-google-conversion (a.k.a. terraform-mapper). Since terraform provider
is mostly handwritten, there are no terraform.yaml file that could be used to
generate the code for terraform-mapper.

Instead, a simplified version of terraform.yaml is provided here which can be
used as a drop-in replacement to generate code for terraform-mapper.

To generate code for terraform mapper, run the following in a new branch:

```
cp products/container/validator/*.yaml products/container 
touch templates/terraform/custom_expand/empty.erb
bundle exec compiler -a -v "ga" -e terraform -f validator -o <output-directory>
rm templates/terraform/custom_expand/empty.erb
git checkout products/container 
```
