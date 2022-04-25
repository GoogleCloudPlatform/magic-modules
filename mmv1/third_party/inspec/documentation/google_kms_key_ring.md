### Test that a GCP kms key ring exists

    describe google_kms_key_ring(project: 'chef-inspec-gcp',  location: 'us-east1', name: 'key-ring-name') do
      it { should exist }
    end

### Test that a GCP kms key ring is in the expected state 

For any existing key ring, below should definitely be true!

    describe google_kms_key_ring(project: 'chef-inspec-gcp',  location: 'us-east1', name: 'key-ring-name') do
      its('create_time_date') { should be > Time.now - 365*60*60*24*50 }
    end