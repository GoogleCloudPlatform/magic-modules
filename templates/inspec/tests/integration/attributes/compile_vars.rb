require 'erb'
require 'yaml'

loaded = YAML.load_file('attributes.yaml')
template = ERB.new(File.read('terraform.tfvars.erb'))
puts template.result_with_hash(loaded)