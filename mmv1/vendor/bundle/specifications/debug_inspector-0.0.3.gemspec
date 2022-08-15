# -*- encoding: utf-8 -*-
# stub: debug_inspector 0.0.3 ruby lib
# stub: ext/debug_inspector/extconf.rb

Gem::Specification.new do |s|
  s.name = "debug_inspector".freeze
  s.version = "0.0.3"

  s.required_rubygems_version = Gem::Requirement.new(">= 0".freeze) if s.respond_to? :required_rubygems_version=
  s.require_paths = ["lib".freeze]
  s.authors = ["John Mair (banisterfiend)".freeze]
  s.date = "2017-05-08"
  s.description = "A Ruby wrapper for the MRI 2.0 debug_inspector API".freeze
  s.email = ["jrmair@gmail.com".freeze]
  s.extensions = ["ext/debug_inspector/extconf.rb".freeze]
  s.files = ["ext/debug_inspector/extconf.rb".freeze]
  s.homepage = "https://github.com/banister/debug_inspector".freeze
  s.licenses = ["MIT".freeze]
  s.rubygems_version = "3.0.9".freeze
  s.summary = "A Ruby wrapper for the MRI 2.0 debug_inspector API".freeze

  s.installed_by_version = "3.0.9" if s.respond_to? :installed_by_version

  if s.respond_to? :specification_version then
    s.specification_version = 4

    if Gem::Version.new(Gem::VERSION) >= Gem::Version.new('1.2.0') then
      s.add_development_dependency(%q<minitest>.freeze, [">= 5"])
    else
      s.add_dependency(%q<minitest>.freeze, [">= 5"])
    end
  else
    s.add_dependency(%q<minitest>.freeze, [">= 5"])
  end
end
