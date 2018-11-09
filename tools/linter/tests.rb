require 'rspec'

RSpec.shared_examples "property_tests" do |disc_prop, api_prop|
  it 'should exist' do
    expect(api_prop).to be_truthy
  end

  it 'should have the same type', skip: api_prop.nil? do
    disc_prop_type = disc_prop.type
    disc_prop_type = 'Api::Type::ResourceRef' if api_prop.class.name == 'Api::Type::ResourceRef'
    expect(api_prop.class.to_s).to eq(disc_prop_type)
  end
end

RSpec.shared_examples "resource_tests" do |disc_res, api_res|
  it 'should have kind', skip: !disc_res.schema.dig('properties', 'kind', 'default') do
    expect(disc_res.schema.dig('properties', 'kind', 'default')).to eq(api_res.kind)
  end
end
