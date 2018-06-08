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

<%= lines(autogen_notice :ruby) -%>

require 'google/storage/network/post'

module Google
  module Storage
    module Api
      # A helper class to provide access to (some) Google Compute Storage API.
      class Object
        def initialize(name, bucket, project, cred)
          @name = name
          @bucket = bucket
          @project = project
          @cred = cred
        end

        # TODO(nelsonjr): Implement this as gstorage_object { }
        # TODO(nelsonjr): Make this function wait for the operation to complete
        def upload(source, type)
          upload_request = ::Google::Compute::Network::Post.new(
            gstorage_object_upload, @cred, type, IO.read(source)
          )
          response = JSON.parse(upload_request.send.body)
          raise Puppet::Error, response['error']['errors'][0]['message'] \
             if response['error']
        end

        private

        STORAGE_UPLOAD_URI =
          'https://www.googleapis.com/upload/storage/v1/b'.freeze

        def gstorage_object_upload
          URI.parse([
            [STORAGE_UPLOAD_URI, @bucket, 'o'].join('/'),
            '?',
            {
              'uploadType' => 'media',
              'name' => @name
            }.map { |k, v| "#{k}=#{v}" }.join('&')
          ].join)
        end
      end
    end
  end
end
