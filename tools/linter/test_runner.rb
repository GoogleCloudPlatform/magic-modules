class TestRunner
  def initialize(discovery_res, api_res)
    @discovery_res = discovery_res
    @api_res = api_res
  end

  def run
    run_on_properties(@discovery_res.properties, @api_res.all_user_properties)
  end

  def run_on_properties(discovery_properties, api_properties, prefix='')
    discovery_properties.each do |disc_prop|
      api_prop = api_properties.select { |p| p.name == disc_prop.name }.first
      test_prop(disc_prop, api_prop, prefix)
      if disc_prop.has_nested_properties?
        run_on_properties(disc_prop.nested_properties, nested_properties_for_api(api_prop), "#{prefix}#{disc_prop.name}.")
      end
    end
  end

  def test_prop(disc, api, prefix)
    unless api
      puts "#{prefix}#{disc.name}: not found"
    end
  end

  private

  def nested_properties_for_api(api)
    if api.is_a?(Api::Type::NestedObject)
      api.properties
    elsif api.is_a?(Api::Type::Array) && api.item_type.is_a?(Api::Type::NestedObject)
      api.item_type.properties
    else
      []
    end
  end
end
