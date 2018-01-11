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

require 'spec_helper'
require 'find'
require 'copyright'

describe 'ensure files have copyright notice' do
  let(:my_path) { File.expand_path(__dir__) }
  let(:my_root) { File.expand_path('..', __dir__) }
  it do
    files = Find.find(File.expand_path(File.join(my_path, '..')))
                .select { |f| File.file?(f) }
                .select do |f|
                  my_tests = f.start_with?("#{my_path}/data/copyright_")
                  artifacts = f.start_with?("#{my_root}/build/")
                  presubmit = f.start_with?("#{my_root}/build/presubmit/")

                  !my_tests && !artifacts || presubmit
                end
    checker = Google::CopyrightChecker.new(files)
    missing = checker.check_missing.collect { |f| "  - #{f}" }
    raise "Files missing (or outdated) copyright:\n#{missing.join("\n")}" \
      unless missing.empty?
  end
end

describe 'check the checker' do
  it 'should pass if all files are good' do
    expect(Google::CopyrightChecker.new(['spec/data/copyright_good1.rb',
                                         'spec/data/copyright_good2.rb'])
                                   .check_missing)
      .to eq []
  end

  it 'should fail if files do not exist' do
    expect do
      Google::CopyrightChecker.new(['spec/data/copyright_good1.rb',
                                    'spec/data/copyright_bad1.rb',
                                    'spec/data/copyright_missing1.rb'])
                              .check_missing
    end.to raise_error(StandardError, /not found.*missing1.rb/)
  end

  it 'should trigger missing copyright' do
    expect(Google::CopyrightChecker.new(['spec/data/copyright_good1.rb',
                                         'spec/data/copyright_bad1.rb',
                                         'spec/data/copyright_good2.rb'])
                                   .check_missing)
      .to contain_exactly 'spec/data/copyright_bad1.rb'
  end

  it 'should fail if year is incorrect' do
    expect(Google::CopyrightChecker.new(['spec/data/copyright_bad2.rb'])
                                   .check_missing)
      .to contain_exactly 'spec/data/copyright_bad2.rb'
  end
end
