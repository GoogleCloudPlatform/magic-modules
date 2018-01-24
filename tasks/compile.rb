require 'tasks/common'

def compile_module(provider, mod)
  output = PROVIDER_FOLDERS[provider.to_sym] % mod
  `./compiler -p products/#{mod} -e #{provider} -o #{output}`
end

def all_compile
  provider_list.map do |x|
    all_tasks_for_provider(x, 'compile:')
  end
end

def all_tasks_for_provider(prov, prefix = '')
  modules_for_provider(prov).map { |x| "#{prefix}#{prov}:#{x}".to_sym }
end

# Compiling Tasks
desc 'Compile all modules'
task compile: all_compile

namespace :compile do
  provider_list.each do |provider|
    # Each provider should default to compiling everything.
    desc "Compile all modules for #{provider.capitalize}"
    multitask provider.to_sym => all_tasks_for_provider(provider)

    namespace provider.to_sym do
      modules_for_provider(provider).each do |mod|
        # Each module should have its own task for compiling.
        desc "Compile the #{mod} module for #{provider.capitalize}"
        task mod.to_sym do
          compile_module(provider, mod)
        end
      end
    end
  end
end
