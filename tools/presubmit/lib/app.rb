module Presubmit
  # An class to run a app full of tests.
  class App
    def initialize(testers, modules)
      @testers = testers
      @modules = modules
    end

    # Run all tests and return the results.
    def run
      Hash[@modules.map { |mod| [mod, mod.run] }]
    end
  end
end
