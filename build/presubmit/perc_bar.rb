#!/bin/ruby
# Copyright 2017 Google Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

module Google
  class ProgressBar
    BAD_COVERAGE = -5 # drop in 5% is really bad
    GOOD_COVERAGE = 95 # happy if coverage is over 95%

    COLOR_CYAN = "\x1b[36m"
    COLOR_GREEN = "\x1b[32m"
    COLOR_RED = "\x1b[31m"
    COLOR_YELLOW = "\x1b[33m"
    COLOR_NONE = "\x1b[0m"
    COLOR_NO_COLOR = ''

    def self.build(before, after, slots)
      step = 100.0 / slots
      real_delta = after - before
      # If any delta, we should at minimum 1 char
      delta = [(real_delta.abs / step).round, 1].max
      base = [([before, after].min / step).round - 1, 0].max
      gap = (slots - delta - base - 1)
      color = if real_delta <= BAD_COVERAGE
                COLOR_RED
              elsif real_delta < 0
                COLOR_YELLOW
              elsif real_delta > 0 && after >= GOOD_COVERAGE
                COLOR_GREEN
              elsif real_delta > 0
                COLOR_CYAN
              else
                COLOR_NO_COLOR
              end
      if before > after
        "#{color}[#{'-' * base}|#{'X' * delta}#{' ' * gap}]#{COLOR_NONE}" \
      elsif before == after
        "#{color}[#{'-' * base}#{'-' * delta}|#{' ' * gap}]" \
      else
        "#{color}[#{'-' * base}|#{'+' * delta}#{' ' * gap}]#{COLOR_NONE}" \
      end
    end
  end
end

# Launched from command line
if __FILE__ == $0
  puts Google::ProgressBar.build(ARGV[0].to_f, ARGV[1].to_f, ARGV[2].to_f)
end
