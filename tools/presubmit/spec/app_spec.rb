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

describe Presubmit::App do
  before do
    @mod1 = Presubmit::MockedModule::Factory.new([Presubmit::SuccessfulTest])
                                            .create
    @mod2 = Presubmit::MockedModule::Factory.new([Presubmit::FailedTest])
                                            .create
  end

  subject do
    described_class.new(
      [Presubmit::SuccessfulTest, Presubmit::FailedTest],
      [@mod1, @mod2]
    )
  end

  it 'should return two modules worth of results' do
    expect(subject.run.length).to be(2)
  end

  it 'should have a first module with one passing result' do
    expect(subject.run[@mod1].map { |_k, v| v.success? }).to eq([true])
  end

  it 'should have a second module with one failling result' do
    expect(subject.run[@mod2].map { |_k, v| v.success? }).to eq([false])
  end
end
