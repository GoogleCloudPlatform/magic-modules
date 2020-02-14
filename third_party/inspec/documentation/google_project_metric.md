### Test that a GCP project metric exists

    describe google_project_metric(project: 'chef-inspec-gcp',  metric: 'metric_name') do
      it { should exist }
    end

### Test that a GCP compute zone has an expected CPU platform

    describe google_project_metric(project: 'chef-inspec-gcp',  metric: 'metric_name') do
      its('filter') { should eq "(protoPayload.serviceName=\"cloudresourcemanager.googleapis.com\")" }
    end