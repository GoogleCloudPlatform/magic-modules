class Test
  def initialize(disc_prop, api_prop, prop_name)
    @disc = disc_prop
    @api = api_prop
    @prop_name = prop_name
  end

  def run
    unless test(@disc, @api)
      puts fail_message(@prop_name)
    end
  end

  def test(disc, api)
    raise "This should be overriden"
  end

  def fail_message(prop_name)
    raise "This should be overriden"
  end
end

class PropExistsTest < Test
  def test(disc, api)
    return api
  end

  def fail_message(prop_name)
    "#{prop_name} does not exist"
  end
end

class TestRunner
  def initialize(discovery_res, api_res)
    @discovery_res = discovery_res
    @api_res = api_res
  end

  def run(&block)
    run_on_properties(@discovery_res.properties, @api_res.all_user_properties, '', &block)
  end

  def run_on_properties(discovery_properties, api_properties, prefix, &block)
    discovery_properties.each do |disc_prop|
      api_prop = api_properties.select { |p| p.name == disc_prop.name }.first
      yield(disc_prop, api_prop, "#{prefix}#{disc_prop.name}")
      if disc_prop.has_nested_properties?
        run_on_properties(disc_prop.nested_properties, nested_properties_for_api(api_prop), "#{prefix}#{disc_prop.name}.", &block)
      end
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

