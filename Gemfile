source 'https://rubygems.org'

gem 'binding_of_caller'

group :test do
  gem 'mocha'
  gem 'rspec'
  # TODO(alexstephen): Monitor rubocop upsteam changes
  # https://github.com/bbatsov/rubocop/pull/4329
  # Change will allow rubocop to use --ignore-parent-exclusion flag
  # Current rubocop upstream will not check Chef files because of
  # AllCops/Exclude
  gem 'rubocop', git: 'https://github.com/nelsonjr/rubocop.git',
                 branch: 'feature/ignore-parent-exclude'
  gem 'simplecov'
end
