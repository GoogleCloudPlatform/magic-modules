### Test that a GCP Cloud SQL Database instance is in the expected state

    describe google_sql_database_instance(project: 'chef-inspec-gcp',  database: 'my-database') do
      its('state') { should eq 'RUNNABLE' }
    end

### Test that a GCP Cloud SQL Database instance generation type

    describe google_sql_database_instance(project: 'chef-inspec-gcp',  database: 'my-database') do
      its('backend_type') { should eq "SECOND_GEN" }
    end

### Test that a GCP Cloud SQL Database instance connection name is as expected

    describe google_sql_database_instance(project: 'spaterson-project',  database: 'gcp-inspec-db-instance') do
      its('connection_name') { should eq "spaterson-project:europe-west2:gcp-inspec-db-instance" }
    end

### Confirm that a GCP Cloud SQL Database instance has the correct version 

    describe google_sql_database_instance(project: 'spaterson-project',  database: 'gcp-inspec-db-instance') do
      its('database_version') { should eq "MYSQL_5_7" }
    end

### Confirm that a GCP Cloud SQL Database instance is running in the desired region and zone 

    describe google_sql_database_instance(project: 'spaterson-project',  database: 'gcp-inspec-db-instance') do
      its('gce_zone') { should eq "europe-west2-a" }
      its('region') { should eq "europe-west2" }
    end