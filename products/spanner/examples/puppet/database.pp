<% if name != 'README.md' -%>
<%= compile 'templates/license.erb' -%>

<%= compile 'templates/autogen_notice.erb' -%>

<%= compile 'templates/puppet/examples~credential.pp.erb' -%>

gspanner_instance_config { 'regional-us-central1':
  project    => 'google.com:graphite-playground',
  credential => 'mycred',
}

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

<% end # name == README.md -%>
gspanner_database { <%= example_resource_name('webstore') -%>:
  ensure           => present,
  extra_statements => [
    'CREATE TABLE customers (
       customer_id INT64 NOT NULL,
       last_name STRING(MAX)
     ) PRIMARY KEY (customer_id)',
  ],
  instance         => <%= example_resource_name('my-spanner') -%>,
  project          => 'google.com:graphite-playground',
  credential       => 'mycred',
}
