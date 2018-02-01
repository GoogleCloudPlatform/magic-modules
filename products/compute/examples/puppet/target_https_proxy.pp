<% if false # the license inside this if block assertains to this file -%>
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
<% unless name == "README.md" -%>
<%= compile 'templates/license.erb' -%>

<%= compile 'templates/autogen_notice.erb' -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

gcompute_zone { 'us-central1-a':
  project    => 'google.com:graphite-playground',
  credential => 'mycred',
}

gcompute_instance_group { <%= example_resource_name('my-puppet-masters') -%>:
  ensure     => present,
  zone       => 'us-central1-a',
  project    => 'google.com:graphite-playground',
  credential => 'mycred',
}

gcompute_backend_service { <%= example_resource_name('my-app-backend') -%>:
  ensure        => present,
  backends      => [
    { group => <%= example_resource_name('my-puppet-masters') -%> },
  ],
  enable_cdn    => true,
  health_checks => [
    gcompute_health_check_ref('another-hc', 'google.com:graphite-playground'),
  ],
  project       => 'google.com:graphite-playground',
  credential    => 'mycred',
}

gcompute_url_map { <%= example_resource_name('my-url-map') -%>:
  ensure          => present,
  default_service => <%= example_resource_name('my-app-backend') -%>,
  project         => 'google.com:graphite-playground',
  credential      => 'mycred',
}

# *******
# WARNING: This manifest is for example purposes only. It is *not* advisable to
# have the key embedded like this because if you check this file into source
# control you are publishing the private key to whomever can access the source
# code. Instead you should protect the key, and for example, use the file()
# function to read it from disk without writing it verbatim to the manifest:
#
# gcompute_ssl_certificate { ...
#   ...
#   private_key => file('/path/to/my/private/key.pem'),
#   ...
# }
# *******

gcompute_ssl_certificate { <%= example_resource_name('sample-certificate') -%>:
  ensure      => present,
  description => 'A certificate for test purposes only.',
  project     => 'google.com:graphite-playground',
  credential  => 'mycred',
  certificate => '-----BEGIN CERTIFICATE-----
MIICqjCCAk+gAwIBAgIJAIuJ+0352Kq4MAoGCCqGSM49BAMCMIGwMQswCQYDVQQG
EwJVUzETMBEGA1UECAwKV2FzaGluZ3RvbjERMA8GA1UEBwwIS2lya2xhbmQxFTAT
BgNVBAoMDEdvb2dsZSwgSW5jLjEeMBwGA1UECwwVR29vZ2xlIENsb3VkIFBsYXRm
b3JtMR8wHQYDVQQDDBZ3d3cubXktc2VjdXJlLXNpdGUuY29tMSEwHwYJKoZIhvcN
AQkBFhJuZWxzb25hQGdvb2dsZS5jb20wHhcNMTcwNjI4MDQ1NjI2WhcNMjcwNjI2
MDQ1NjI2WjCBsDELMAkGA1UEBhMCVVMxEzARBgNVBAgMCldhc2hpbmd0b24xETAP
BgNVBAcMCEtpcmtsYW5kMRUwEwYDVQQKDAxHb29nbGUsIEluYy4xHjAcBgNVBAsM
FUdvb2dsZSBDbG91ZCBQbGF0Zm9ybTEfMB0GA1UEAwwWd3d3Lm15LXNlY3VyZS1z
aXRlLmNvbTEhMB8GCSqGSIb3DQEJARYSbmVsc29uYUBnb29nbGUuY29tMFkwEwYH
KoZIzj0CAQYIKoZIzj0DAQcDQgAEHGzpcRJ4XzfBJCCPMQeXQpTXwlblimODQCuQ
4mzkzTv0dXyB750fOGN02HtkpBOZzzvUARTR10JQoSe2/5PIwaNQME4wHQYDVR0O
BBYEFKIQC3A2SDpxcdfn0YLKineDNq/BMB8GA1UdIwQYMBaAFKIQC3A2SDpxcdfn
0YLKineDNq/BMAwGA1UdEwQFMAMBAf8wCgYIKoZIzj0EAwIDSQAwRgIhALs4vy+O
M3jcqgA4fSW/oKw6UJxp+M6a+nGMX+UJR3YgAiEAvvl39QRVAiv84hdoCuyON0lJ
zqGNhIPGq2ULqXKK8BY=
-----END CERTIFICATE-----',
  private_key => '-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIObtRo8tkUqoMjeHhsOh2ouPpXCgBcP+EDxZCB/tws15oAoGCCqGSM49
AwEHoUQDQgAEHGzpcRJ4XzfBJCCPMQeXQpTXwlblimODQCuQ4mzkzTv0dXyB750f
OGN02HtkpBOZzzvUARTR10JQoSe2/5PIwQ==
-----END EC PRIVATE KEY-----',
}

<% end # name == README.md -%>
gcompute_target_https_proxy { <%= example_resource_name('my-https-proxy') -%>:
  ensure           => present,
  ssl_certificates => [
    <%= example_resource_name('sample-certificate') -%>,
  ],
  url_map          => <%= example_resource_name('my-url-map') -%>,
  project          => 'google.com:graphite-playground',
  credential       => 'mycred',
}
