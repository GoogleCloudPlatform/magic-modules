PROVIDER_FOLDERS = {
  ansible: 'build/ansible',
  puppet: 'build/puppet/%s',
  chef: 'build/chef/%s'
}.freeze

# Give a list of all providers served by MM
def provider_list
  PROVIDER_FOLDERS.keys
end

# Give a list of all products served by a provider
def modules_for_provider(provider)
  products = File.join(File.dirname(__FILE__), '..', 'products')
  files = Dir.glob("#{products}/**/#{provider}.yaml")
  files.map do |file|
    match = file.match(%r{^.*products\/([a-z]*)\/.*yaml.*})
    match.captures[0] if match
  end.compact
end
