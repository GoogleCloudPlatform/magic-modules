dlext = RbConfig::CONFIG['DLEXT']

mri_2 = defined?(RUBY_ENGINE) && RUBY_ENGINE == "ruby" &&
    RUBY_VERSION =~ /^2/

if mri_2
  require 'binding_of_caller/mri2'
elsif defined?(RUBY_ENGINE) && RUBY_ENGINE == "ruby"
  require "binding_of_caller.#{dlext}"
elsif defined?(Rubinius)
  require 'binding_of_caller/rubinius'
elsif defined?(JRuby)
  require 'binding_of_caller/jruby_interpreted'
end
