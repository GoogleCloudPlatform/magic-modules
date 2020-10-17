### Test that there are no more than a specified number of users available for the project

    describe google_sql_users(project: 'chef-inspec-gcp', database: 'database-instance') do
      its('count') { should be <= 100}
    end

### Test that an expected user is available for the project

    describe google_sql_users(project: 'chef-inspec-gcp') do
      its('user_names') { should include "us-east1-b" }
    end

### Test whether any users are in status "DOWN"

    describe google_sql_users(project: 'chef-inspec-gcp') do
      its('user_statuses') { should_not include "DOWN" }
    end

### Test users exist for all database instances in a project

    google_sql_database_instances(project: 'chef-inspec-gcp').instance_names.each do |instance_name|
      describe google_sql_users(project: 'chef-inspec-gcp', database: instance_name) do
        it { should exist }
      end
    end