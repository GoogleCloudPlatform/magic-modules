def loop_resources_in_api(api, discovery, &block)
  names = api.all_resource_names
  names.each { |name| yield api.resource(name), discovery.resource(name) }
end
