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
  s.files = [".gemtest".freeze, ".gitignore".freeze, ".travis.yml".freeze, ".yardopts".freeze, "Gemfile".freeze, "HISTORY".freeze, "LICENSE".freeze, "README.md".freeze, "Rakefile".freeze, "binding_of_caller.gemspec".freeze, "examples/benchmark.rb".freeze, "examples/example.rb".freeze, "ext/binding_of_caller/binding_of_caller.c".freeze, "ext/binding_of_caller/extconf.rb".freeze, "ext/binding_of_caller/ruby_headers/192/debug.h".freeze, "ext/binding_of_caller/ruby_headers/192/dln.h".freeze, "ext/binding_of_caller/ruby_headers/192/eval_intern.h".freeze, "ext/binding_of_caller/ruby_headers/192/id.h".freeze, "ext/binding_of_caller/ruby_headers/192/iseq.h".freeze, "ext/binding_of_caller/ruby_headers/192/method.h".freeze, "ext/binding_of_caller/ruby_headers/192/node.h".freeze, "ext/binding_of_caller/ruby_headers/192/regenc.h".freeze, "ext/binding_of_caller/ruby_headers/192/regint.h".freeze, "ext/binding_of_caller/ruby_headers/192/regparse.h".freeze, "ext/binding_of_caller/ruby_headers/192/rubys_gc.h".freeze, "ext/binding_of_caller/ruby_headers/192/thread_pthread.h".freeze, "ext/binding_of_caller/ruby_headers/192/thread_win32.h".freeze, "ext/binding_of_caller/ruby_headers/192/timev.h".freeze, "ext/binding_of_caller/ruby_headers/192/transcode_data.h".freeze, "ext/binding_of_caller/ruby_headers/192/version.h".freeze, "ext/binding_of_caller/ruby_headers/192/vm_core.h".freeze, "ext/binding_of_caller/ruby_headers/192/vm_exec.h".freeze, "ext/binding_of_caller/ruby_headers/192/vm_insnhelper.h".freeze, "ext/binding_of_caller/ruby_headers/192/vm_opts.h".freeze, "ext/binding_of_caller/ruby_headers/193/addr2line.h".freeze, "ext/binding_of_caller/ruby_headers/193/atomic.h".freeze, "ext/binding_of_caller/ruby_headers/193/constant.h".freeze, "ext/binding_of_caller/ruby_headers/193/debug.h".freeze, "ext/binding_of_caller/ruby_headers/193/dln.h".freeze, "ext/binding_of_caller/ruby_headers/193/encdb.h".freeze, "ext/binding_of_caller/ruby_headers/193/eval_intern.h".freeze, "ext/binding_of_caller/ruby_headers/193/id.h".freeze, "ext/binding_of_caller/ruby_headers/193/internal.h".freeze, "ext/binding_of_caller/ruby_headers/193/iseq.h".freeze, "ext/binding_of_caller/ruby_headers/193/method.h".freeze, "ext/binding_of_caller/ruby_headers/193/node.h".freeze, "ext/binding_of_caller/ruby_headers/193/parse.h".freeze, "ext/binding_of_caller/ruby_headers/193/regenc.h".freeze, "ext/binding_of_caller/ruby_headers/193/regint.h".freeze, "ext/binding_of_caller/ruby_headers/193/regparse.h".freeze, "ext/binding_of_caller/ruby_headers/193/revision.h".freeze, "ext/binding_of_caller/ruby_headers/193/rubys_gc.h".freeze, "ext/binding_of_caller/ruby_headers/193/thread_pthread.h".freeze, "ext/binding_of_caller/ruby_headers/193/thread_win32.h".freeze, "ext/binding_of_caller/ruby_headers/193/timev.h".freeze, "ext/binding_of_caller/ruby_headers/193/transcode_data.h".freeze, "ext/binding_of_caller/ruby_headers/193/transdb.h".freeze, "ext/binding_of_caller/ruby_headers/193/version.h".freeze, "ext/binding_of_caller/ruby_headers/193/vm_core.h".freeze, "ext/binding_of_caller/ruby_headers/193/vm_exec.h".freeze, "ext/binding_of_caller/ruby_headers/193/vm_insnhelper.h".freeze, "ext/binding_of_caller/ruby_headers/193/vm_opts.h".freeze, "lib/binding_of_caller.rb".freeze, "lib/binding_of_caller/jruby_interpreted.rb".freeze, "lib/binding_of_caller/mri2.rb".freeze, "lib/binding_of_caller/rubinius.rb".freeze, "lib/binding_of_caller/version.rb".freeze, "test/test_binding_of_caller.rb".freeze]
  s.homepage = "http://github.com/banister/binding_of_caller".freeze
  s.rubygems_version = "2.6.14".freeze
  s.summary = "Retrieve the binding of a method's caller. Can also retrieve bindings even further up the stack.".freeze
  s.test_files = ["test/test_binding_of_caller.rb".freeze]

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
