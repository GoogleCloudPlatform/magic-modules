### Test that there are no more than a specified number of zones available for the project

    describe google_sql_database_instances(project: 'chef-inspec-gcp') do
      its('count') { should be <= 100}
    end


### Test that a database instance exists in the expected zone  

    describe google_sql_database_instances(project: 'chef-inspec-gcp') do
      its('instance_zones') { should include "us-east1-b" }
    end

### Test that a database instance exists in the expected region  

    describe google_sql_database_instances(project: 'chef-inspec-gcp') do
      its('instance_regions') { should include "us-east1" }
    end


### Confirm that at least one database instance is in "RUNNABLE" state

    describe google_sql_database_instances(project: 'chef-inspec-gcp') do
      its('instance_states') { should include "RUNNABLE" }
    end

### Test that a subset of all database instances matching "*mysqldb*" are all version "MYSQL_5_7"

    google_sql_database_instances(project: 'chef-inspec-gcp').where(instance_name: /mysqldb/).instance_names.each do |instance_name|
      describe google_sql_database_instance(project: 'chef-inspec-gcp',  database: instance_name) do
        it { should exist }
        its('database_version') { should eq "MYSQL_5_7" }
      end
    end