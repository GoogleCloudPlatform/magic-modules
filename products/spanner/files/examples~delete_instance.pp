<% if name != 'README.md' -%>
<%= compile 'templates/license.erb' -%>

<%= compile 'templates/autogen_notice.erb' -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

<% end # name == README.md -%>
gspanner_instance { <%= example_resource_name('my-spanner') -%>:
  ensure     => absent,
  project    => 'google.com:graphite-playground',
  credential => 'mycred',
}
