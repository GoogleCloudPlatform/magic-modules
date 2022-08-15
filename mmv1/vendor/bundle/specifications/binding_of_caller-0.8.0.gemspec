# -*- encoding: utf-8 -*-
# stub: binding_of_caller 0.8.0 ruby lib
# stub: ext/binding_of_caller/extconf.rb

Gem::Specification.new do |s|
  s.name = "binding_of_caller".freeze
  s.version = "0.8.0"

  s.required_rubygems_version = Gem::Requirement.new(">= 0".freeze) if s.respond_to? :required_rubygems_version=
  s.require_paths = ["lib".freeze]
  s.authors = ["John Mair (banisterfiend)".freeze]
  s.date = "2018-01-10"
  s.description = "Retrieve the binding of a method's caller. Can also retrieve bindings even further up the stack.".freeze
  s.email = "jrmair@gmail.com".freeze
  s.extensions = ["ext/binding_of_caller/extconf.rb".freeze]
  s.files = ["ext/binding_of_caller/extconf.rb".freeze]
  s.homepage = "http://github.com/banister/binding_of_caller".freeze
  s.rubygems_version = "3.0.9".freeze
  s.summary = "Retrieve the binding of a method's caller. Can also retrieve bindings even further up the stack.".freeze

  s.installed_by_version = "3.0.9" if s.respond_to? :installed_by_version

  if s.respond_to? :specification_version then
    s.specification_version = 4

    if Gem::Version.new(Gem::VERSION) >= Gem::Version.new('1.2.0') then
      s.add_runtime_dependency(%q<debug_inspector>.freeze, [">= 0.0.1"])
      s.add_development_dependency(%q<bacon>.freeze, [">= 0"])
      s.add_development_dependency(%q<rake>.freeze, [">= 0"])
    else
      s.add_dependency(%q<debug_inspector>.freeze, [">= 0.0.1"])
      s.add_dependency(%q<bacon>.freeze, [">= 0"])
      s.add_dependency(%q<rake>.freeze, [">= 0"])
    end
  else
    s.add_dependency(%q<debug_inspector>.freeze, [">= 0.0.1"])
    s.add_dependency(%q<bacon>.freeze, [">= 0"])
    s.add_dependency(%q<rake>.freeze, [">= 0"])
  end
end
