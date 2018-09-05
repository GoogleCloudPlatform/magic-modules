<% if false # the license inside this if block pertains to this file -%>
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
<% end -%>
<%# 
This file is complicated, and deserves some documentation here, so that
you, brave adventurer, can understand what's happened.  This file's
goal is to run the examples required by puppet compilation.  Essentially,
you can think of it as a shell script that runs 'puppet apply examples/foo.pp',
then 'puppet apply examples/delete_foo.pp', and checks for errors.

It has to be more complicated than that.  :)
"Fiddling with puppet internal APIs is a lonely and painful path to take."
- David Schmitt, Puppet SWE, September 2018

First, let's talk about the reason that it *isn't* a shell script.
We need unit tests that really, really validate our puppet modules - we
can't have them breaking on us, because the community isn't active enough
that we can trust that we'll find out through regular channels.  We
absolutely must have tests that we can count on to tell us if the module
is working as expected.  And those tests need to run on every PR, because
these modules are complicated enough that they can break in subtle ways.
Consequently, we need to be able to run tests that validate use against the
real GCP APIs, but without spending incredible amounts of time (or money)
or actually creating hordes of resources.  That means a VCR test.

A VCR test is a kind of test that runs once, making real HTTP requests along
the way.  It records the requests and responses.  On later runs, it mocks
those same HTTP requests with those same responses, making your tests
repeatable (desirable), fast (desirable), and in our case, also free!

You can't write a test like that using the puppet command line tool, so
instead we plumb into the guts of puppet.  Our main hook is a function called
apply_compiled_manifest (and its sibling, apply_with_error_check).

These functions are not really very similar to `puppet apply`.  Instead of
taking a filename, they take the entire manifest as a string.  This causes
subtle problems deep within Puppet - Puppet depends on the filename of the
manifest in a handful of ways, not least of which is to find out what
module it is part of (to automagically load dependencies and functions).
Since that isn't possible, we tack on the Puppet override of our own
environment and loader.  Since these functions don't call out to Facter
the way that `puppet apply` does, we need to substitute in our own project
name and credential path (which all our examples use).  Hence the
get_example() function - it does more than just read the file.

You can see all this stuff in spec_helper.rb.erb.

There's also the issue of the begin/rescue pattern you'll see here.  It's
ugly!  It's also necessary.  compile_to_ral() calls out to our mocked up
loader, but fails to successfully compile the manifest due to some
initialization which is mandatory, but isn't possible to complete from
outside puppet.  HOWEVER!  We're in luck!  In calling some cleanup code,
in reporting on the failure, compile_to_ral actually *performs* the
initialization we need.  So, if we try it twice, it works the second
time.  Lucky.

So!  We *can* successfully apply a manifest from
under webmock, as long as we:
- try twice, eating the first failure.
- don't need to use any functions that are loaded from modules
- ensure all our dependent modules have their 'lib' dir in $LOAD_PATH
- do all our own variable substitution
- rewrite any necessary functions in puppet and stick them in
    the manifest.
- import all the magic from puppetlabs_spec_helper/module_spec_helper,
    which I have done without fully understanding it.  The comments in
    the code suggest that the authors also do not fully understand it.

Hopefully that explains the complexity here in a way that will continue
to be useful as the tests develop further.

Let me also explain the order in which these commands are run.  First,
there is a pre-destroy - in the event that an older test failed to clean
up after itself, the next run should take care of it.  Second, we create
the resources in the manifest.  Third, we run the create again, mocking
out any requests during 'flush' and making them crash the test - this is
to check idempotency of the manifests.  Fourth, we delete the resources
in the manifest, and fifth, we check the idempotency of that delete by
confirming that no resources make web requests during 'flush'.

VCR recordings are called 'cassettes', and they are stored in
'spec/cassettes'.  If present, they are used - if absent, created.
#-%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>
require 'spec_helper'
require 'vcr'

VCR.configure do |c|
  c.cassette_library_dir = 'spec/cassettes'
  c.hook_into :webmock
  c.configure_rspec_metadata!
end

describe '<%= Google::StringUtils.underscore(obj.name) -%>.create', vcr: true do
  it 'creates and destroys non-existent <%= Google::StringUtils.underscore(obj.name) -%>' do
    puts 'pre-destroying <%= Google::StringUtils.underscore(obj.name) -%>'
    VCR.use_cassette('pre_destroy_<%= Google::StringUtils.underscore(obj.name) -%>') do
      run_example('delete_<%= Google::StringUtils.underscore(obj.name) -%>')
    end
    puts 'creating <%= Google::StringUtils.underscore(obj.name) -%>'
    VCR.use_cassette('create_<%= Google::StringUtils.underscore(obj.name) -%>') do
      run_example('<%= Google::StringUtils.underscore(obj.name) -%>')
    end
    puts 'checking that <%= Google::StringUtils.underscore(obj.name) -%> is created'
    VCR.use_cassette('check_<%= Google::StringUtils.underscore(obj.name) -%>') do
      validate_no_flush_calls('<%= Google::StringUtils.underscore(obj.name) -%>')
    end
    puts 'destroying <%= Google::StringUtils.underscore(obj.name) -%>'
    VCR.use_cassette('destroy_<%= Google::StringUtils.underscore(obj.name) -%>') do
      run_example('delete_<%= Google::StringUtils.underscore(obj.name) -%>')
    end
    puts 'confirming <%= Google::StringUtils.underscore(obj.name) -%> destroyed'
    VCR.use_cassette('check_destroy_<%= Google::StringUtils.underscore(obj.name) -%>') do
      validate_no_flush_calls('delete_<%= Google::StringUtils.underscore(obj.name) -%>')
    end
  end
end
