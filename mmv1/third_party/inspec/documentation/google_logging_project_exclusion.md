### Test that a GCP project logging exclusion name is as expected

    describe google_logging_project_exclusion(project: 'chef-inspec-gcp',  exclusion: 'exclusion-name-abcd') do
      its('name') { should eq 'exclusion-name-abcd' }
    end

### Test that a GCP project logging exclusion filter is set correctly

    describe google_logging_project_exclusion(project: 'chef-inspec-gcp',  exclusion: 'exclusion-name-abcd') do
      its('filter') { should eq 'resource.type = gce_instance AND severity <= DEBUG' }
    end

### Test that a GCP project logging exclusion description is as expected

    describe google_logging_project_exclusion(project: 'chef-inspec-gcp',  exclusion: 'exclusion-name-abcd') do
      its('description') { should eq 'Exclude GCE instance debug logs' }
    end