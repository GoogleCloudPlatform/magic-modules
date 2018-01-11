# Copyright 2017 Google Inc.
# Licensed under the Apache License, Version 2.0 (the 'License');
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an 'AS IS' BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

$LOAD_PATH.unshift(File.expand_path('build/presubmit'))

require 'spec_helper'
require 'perc_bar'

describe Google::ProgressBar do
  context 'same amount' do
    subject { described_class.build(10, 10, 20) }
    it { is_expected.to eq '[--|                 ]' }
  end

  context 'increased a little' do
    subject { described_class.build(10, 10.1, 20) }
    it { is_expected.to include '[-|+                 ]' }
    it { is_expected.to start_with Google::ProgressBar::COLOR_CYAN }
    it { is_expected.to end_with Google::ProgressBar::COLOR_NONE }
  end

  context 'increased a lot' do
    subject { described_class.build(10, 60, 20) }
    it { is_expected.to include '[-|++++++++++        ]' }
    it { is_expected.to start_with Google::ProgressBar::COLOR_CYAN }
    it { is_expected.to end_with Google::ProgressBar::COLOR_NONE }
  end

  context 'decreased a little' do
    subject { described_class.build(10.1, 10, 20) }
    it { is_expected.to include '[-|X                 ]' }
    it { is_expected.to start_with Google::ProgressBar::COLOR_YELLOW }
    it { is_expected.to end_with Google::ProgressBar::COLOR_NONE }
  end

  context 'decreased a lot' do
    subject { described_class.build(60, 10, 20) }
    it { is_expected.to include '[-|XXXXXXXXXX        ]' }
    it { is_expected.to start_with Google::ProgressBar::COLOR_RED }
    it { is_expected.to end_with Google::ProgressBar::COLOR_NONE }
  end

  context 'small start' do
    subject { described_class.build(0, 10, 20) }
    it { is_expected.to include '[|++                 ]' }
    it { is_expected.to start_with Google::ProgressBar::COLOR_CYAN }
    it { is_expected.to end_with Google::ProgressBar::COLOR_NONE }
  end

  context 'small end' do
    subject { described_class.build(10, 0, 20) }
    it { is_expected.to include '[|XX                 ]' }
    it { is_expected.to start_with Google::ProgressBar::COLOR_RED }
    it { is_expected.to end_with Google::ProgressBar::COLOR_NONE }
  end

  context 'good coverage. happy color' do
    subject { described_class.build(60, 95, 20) }
    it { is_expected.to start_with Google::ProgressBar::COLOR_GREEN }
  end

  context 'bad coverage. dropping the ball' do
    subject { described_class.build(90, 85, 20) }
    it { is_expected.to start_with Google::ProgressBar::COLOR_RED }
  end
end
