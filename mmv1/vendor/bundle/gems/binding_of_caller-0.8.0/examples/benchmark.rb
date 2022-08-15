require 'benchmark'

unless Object.const_defined? :BindingOfCaller
  $:.unshift File.expand_path '../../lib', __FILE__
  require 'binding_of_caller'
  require 'binding_of_caller/version'
end


n = 250000

Benchmark.bm(10) do |x|
  x.report("#of_caller") do
    1.upto(n) do
      1.times do
        1.times do
          binding.of_caller(2)
          binding.of_caller(1)
        end
      end
    end
  end

  x.report("#frame_count") do
    1.upto(n) do
      1.times do
        1.times do
          binding.frame_count
        end
      end
    end
  end

  x.report("#callers") do
    1.upto(n) do
      1.times do
        1.times do
          binding.callers
        end
      end
    end
  end

  x.report("#frame_description") do
    1.upto(n) do
      1.times do
        1.times do
          binding.of_caller(1).frame_description
        end
      end
    end
  end

  x.report("#frame_type") do
    1.upto(n) do
      1.times do
        1.times do
          binding.of_caller(1).frame_type
        end
      end
    end
  end
end
