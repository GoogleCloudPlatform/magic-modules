require 'rspec'

RSpec.shared_examples 'property_tests' do |_disc_prop, api_prop|
  it 'should exist' do
    expect(api_prop).to be_truthy
  end
end

RSpec.shared_examples 'resource_tests' do |disc_res, api_res|
  it 'should have kind', skip: !disc_res.schema.dig('properties', 'kind', 'default') do
    expect(disc_res.schema.dig('properties', 'kind', 'default')).to eq(api_res.kind)
  end
end
