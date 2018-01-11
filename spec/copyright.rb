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

require 'find'

module Google
  # Enforces copyright notices on source files.
  class CopyrightChecker
    def initialize(files)
      @input_files = files
      @copyrightable_files = {
        /.*\.yaml$/ => /^# Copyright 201[78]/,
        /.*\.rb$/ => /^# Copyright 201[78]/
      }
    end

    def check_missing
      checks = [method(:exist?),
                method(:file?),
                method(:suitable?),
                method(:notice_present?)]
      files = @input_files
      checks.each { |test| files = test.call(files) }
      files
    end

    private

    def exist?(files)
      not_found = files.reject { |f| File.exist?(f) }
      raise "Some files were not found: #{not_found}" unless not_found.empty?
      files
    end

    def file?(files)
      not_files = files.reject { |f| File.file?(f) }
      raise "Some inputs were not files: #{not_files}" unless not_files.empty?
      files
    end

    def suitable?(files)
      files.reject do |f|
        @copyrightable_files.select { |c, _| c =~ f }.empty?
      end
    end

    def notice_present?(files)
      files.select do |f|
        mark = @copyrightable_files.reject { |c, _| (c =~ f).nil? }
        File.readlines(f).select { |l| mark.values[0] =~ l }.empty?
      end
    end
  end
end
