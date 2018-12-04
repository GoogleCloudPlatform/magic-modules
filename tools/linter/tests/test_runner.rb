require 'tools/linter/tests/tests'

def run_tests(discovery_doc, api, tags={})
  # First context: product name
  RSpec.describe api.prefix do
    discovery_doc.resources.each do |disc_resource|
      api_obj = api.objects.select { |p| p.name == disc_resource.name }.first
      # Second context: resource name
      describe disc_resource.name do
        # Run all resource tests on this resource
        if tags[:resource]
          include_examples 'resource_tests', disc_resource, api_obj
        end

        if tags[:property]
          PropertyFetcher.fetch_property_pairs(disc_resource.properties,
                                               api_obj.all_user_properties) \
                                              do |disc_prop, api_prop, name|
            # Third context: property name
            context name do
              # Run all tests on this property
              include_examples 'property_tests', disc_prop, api_prop
            end
          end
        end
      end
    end
  end
end
