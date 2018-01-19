module Presubmit
  # An interface class for a Presubmit test.
  class Test
    def initialize(mod)
      @mod = mod
    end

    # This will run a test and return its results.
    def run
      raise 'This must be implemented by the test'
    end
  end

  # A class that stores the results of a test.
  class Results
    attr_reader :tester
    attr_reader :status
    attr_reader :output

    def initialize(tester, status, output)
      @tester = tester
      @status = status
      @output = output
    end

    def success?
      @status.zero?
    end

    def failed?
      !success?
    end

    def warning?
      false
    end
  end
end
