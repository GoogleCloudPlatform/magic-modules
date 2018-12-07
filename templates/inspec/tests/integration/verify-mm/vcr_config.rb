
require 'vcr'

VCR.configure do |c|
  c.hook_into :webmock
  c.cassette_library_dir = 'inspec-cassettes'
  c.allow_http_connections_when_no_cassette = true

  c.before_record do |i|
    i.response.headers.delete_if { |key| key != 'Content-Type' }
    i.request.headers.delete_if { |key| true }
    if auth_call?(i)
      i.request.body = 'AUTH REQUEST'
      i.response.body = "{\n  \"access_token\": \"ya29.c.samsamsamsamsamsamsamsamsa-thisisnintysixcharactersoftexttolooklikeanauthtokenthisisnintysixcharactersoftexttolooklikeanaut\",\n  \"expires_in\": 3600,\n  \"token_type\": \"Bearer\"\n}"
    end
  end
end

def auth_call?(interaction)
	# Auth calls require extra scrubbing, this method is very broad, this is intentional
  interaction.request.uri.include? 'oauth2'
end