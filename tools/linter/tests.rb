require 'rspec'

RSpec.shared_examples "property_tests" do |disc_prop, api_prop|
  it 'should exist' do
    expect(api_prop).to be_truthy
  end

  it 'should have the same type' do
    disc_prop_type = disc_prop.type
    disc_prop_type = 'Api::Type::ResourceRef' if api_prop.class.name == 'Api::Type::ResourceRef'
    expect(api_prop.class.to_s).to eq(disc_prop.type)
  end
end
