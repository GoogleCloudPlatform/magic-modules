<% if name != 'README.md' -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

gspanner_instance_config { 'regional-us-central1':
  project    => 'google.com:graphite-playground',
  credential => 'mycred',
}

<% end # name == README.md -%>
gspanner_instance { <%= example_resource_name('my-spanner') -%>:
  display_name => 'My Spanner Instance',
  node_count   => 2,
  labels       => [
    {
      'cost-center' => 'ti-1700004',
    },
  ],
  config       => 'regional-us-central1',
  project      => 'google.com:graphite-playground',
  credential   => 'mycred',
}
