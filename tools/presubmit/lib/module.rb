module Presubmit
  # An interface class for a Submodule
  class Module
    def add_tests(tests)
      @tests = tests
    end

    # Creates modules with properly initialized tests.
    # This is where per-module and per-product test differences are made.
    class Factory
      def initialize(testers)
        @testers = testers
      end

      def create
        raise 'children must implement'
      end
    end

    # Run all tests and return the results.
    def run
      Hash[
        @tests.map do |test|
          [test, test.run]
        end
      ]
    end
  end
end
