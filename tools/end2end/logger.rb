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

require 'singleton'

# Simple logger w/ atomic timestamp
class Logger
  include Singleton

  def initialize
    @semaphore = Mutex.new
    @logged = 0
  end

  def log(provider, product, message)
    message.gsub("\n\n", "\n \n")
           .split("\n")
           .each do |l|
             if product.nil?
               printf "%<where>-30s %<when>s: %<log>s\n",
                      where: provider, when: "[#{timestamp}]", log: l
             else
               printf "%<where>-30s %<when>s: %<log>s\n",
                      where: "#{provider}/#{product}", when: "[#{timestamp}]",
                      log: l
             end
           end
  end

  private

  def timestamp
    @semaphore.synchronize do
      @logged += 1
      format '%05d', @logged
    end
  end
end
