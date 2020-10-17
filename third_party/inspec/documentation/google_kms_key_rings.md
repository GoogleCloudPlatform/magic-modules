### Test that there are no more than a specified number of kms_key_rings available for the project

    describe google_kms_key_rings(project: 'chef-inspec-gcp', location: 'us-east1') do
      its('count') { should be <= 200}
    end

### Test that an expected kms_key_ring is available for the project

    describe google_kms_key_rings(project: 'chef-inspec-gcp', location: 'us-east1') do
      its('key_ring_names') { should include "a-named-key" }
    end


### Test that all KMS key rings were created in the past year

    describe google_kms_key_rings(project: gcp_project_id, location: 'us-east1').key_ring_names.each do |key_ring_name|
      describe google_kms_key_ring(project: 'chef-inspec-gcp', location: 'us-east1', 'name: key_ring_name) do
        it { should exist }
        its('create_time_date') { should be > Time.now - 365*60*60*24 }
      end
    end