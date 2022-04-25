### Test that there are no more than a specified number of metrics available for the project

    describe google_project_metrics(project: 'chef-inspec-gcp') do
      its('count') { should be <= 100}
    end

### Test that an expected metric name is available for the project

    describe google_project_metrics(project: 'chef-inspec-gcp') do
      its('metric_names') { should include "metric-name" }
    end

### Test that a subset of all metrics with name matching "*project*" have a particular writer identity 

    google_project_metrics(project: 'chef-inspec-gcp').where(metric_name: /project/).metric_names.each do |metric_name|
      describe google_project_metric(project: 'chef-inspec-gcp',  metric: metric_name) do
        its('filter') { should eq "(protoPayload.serviceName=\"cloudresourcemanager.googleapis.com\")" }
      end
    end