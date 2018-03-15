<% if name != 'README.md' -%>
<%= compile 'templates/license.erb' -%>

<%= lines(autogen_notice :puppet) -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

<% end # name == README.md -%>
gspanner_instance_config { 'regional-us-central1':
  project    => 'google.com:graphite-playground',
  credential => 'mycred',
}
