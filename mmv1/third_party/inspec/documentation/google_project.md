### Test that a GCP project has the expected project number

    describe google_project(project: 'chef-inspec-gcp') do
      its('project_number') { should eq 12345678 }
    end

### Test that a GCP project has the expected lifecycle state e.g. "ACTIVE"

    describe google_project(project: 'chef-inspec-gcp') do
      its('lifecycle_state') { should eq "ACTIVE" }
    end

### Validate that a GCP project has some arbitrary label with expected content (for example defined by regexp )

    describe google_project(project: 'chef-inspec-gcp').label_value_by_key('season') do
      it {should match '^(winter|spring|summer|autumn)$' }
    end