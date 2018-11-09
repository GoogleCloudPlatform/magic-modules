require 'rspec'

RSpec.shared_examples "property_tests" do |disc_prop, api_prop|
  it 'should exist' do
    expect(api_prop).to be_truthy
  end
end
