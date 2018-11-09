class TestRunner
  def initialize(discovery_res, api_res)
    @discovery_res = discovery_res
    @api_res = api_res
  end

  def run
    @discovery_res.properties.each do |disc_prop|
      api_prop = @api_res.all_user_properties.select { |p| disc_prop.name == p.name }.first
      unless test_prop(disc_prop, api_prop)
        puts "#{disc_prop.name}: not found"
      end
    end
  end

  def test_prop(disc, api)
    return api
  end
end
