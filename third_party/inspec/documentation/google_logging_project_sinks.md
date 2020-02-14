### Test that there are no more than a specified number of sinks available for the project

    describe google_logging_project_sinks(project: 'chef-inspec-gcp') do
      its('count') { should be <= 100}
    end

### Test that an expected sink name is available for the project

    describe google_logging_project_sinks(project: 'chef-inspec-gcp') do
      its('sink_names') { should include "my-sink" }
    end

### Test that an expected sink destination is available for the project

    describe google_logging_project_sinks(project: 'chef-inspec-gcp') do
      its('sink_destinations') { should include "storage.googleapis.com/a-logging-bucket" }
    end

### Test that a subset of all sinks matching "project*" have a particular writer identity 

    google_logging_project_sinks(project: 'chef-inspec-gcp').where(sink_name: /project/).sink_names.each do |sink_name|
      describe google_logging_project_sink(project: 'chef-inspec-gcp',  sink: sink_name) do
        its('writer_identity') { should eq "serviceAccount:my-logging-service-account.iam.gserviceaccount.com" }
      end
    end