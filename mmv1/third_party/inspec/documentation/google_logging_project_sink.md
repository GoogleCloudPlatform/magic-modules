### Test that a GCP project logging sink destination is correct

    describe google_logging_project_sink(project: 'chef-inspec-gcp',  sink: 'sink-name-abcd') do
      its('destination') { should eq 'storage.googleapis.com/gcp-inspec-logging-bucket' }
    end

### Test that a GCP project logging sink filter is correct

    describe google_logging_project_sink(project: 'chef-inspec-gcp',  sink: 'sink-name-abcd') do
      its('filter') { should eq "resource.type = gce_instance AND resource.labels.instance_id = \"12345678910123123\"" }
    end

### Test a GCP project logging sink output version format

    describe google_logging_project_sink(project: 'chef-inspec-gcp',  sink: 'sink-name-abcd') do
      its('output_version_format') { should eq "V2" }
    end

### Test a GCP project logging sink writer identity is as expected

    describe google_logging_project_sink(project: 'chef-inspec-gcp',  sink: 'sink-name-abcd') do
      its('writer_identity') { should eq "serviceAccount:my-logging-service-account.iam.gserviceaccount.com" }
    end
