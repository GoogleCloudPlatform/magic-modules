### Test that a GCP KMS crypto key was created recently

    describe google_kms_crypto_key(project: 'chef-inspec-gcp',   location: 'us-east1',  key_ring_name: 'key-ring', name: 'crypto-key') do
      its('create_time_date') { should be > Time.now - 365*60*60*24*10 }
    end

### Test when the next rotation time for a GCP KMS crypto key is scheduled

    describe google_kms_crypto_key(project: 'chef-inspec-gcp',   location: 'us-east1',  key_ring_name: 'key-ring', name: 'crypto-key') do
      its('next_rotation_time_date') { should be > Time.now - 100000 }
    end
    
### Check that the crypto key purpose is as expected
    
    describe google_kms_crypto_key(project: 'chef-inspec-gcp',   location: 'us-east1',  key_ring_name: 'key-ring', name: 'crypto-key') do
      its('purpose') { should eq "ENCRYPT_DECRYPT" }
    end

### Check that the crypto key primary is in "ENABLED" state
    
    describe google_kms_crypto_key(project: 'chef-inspec-gcp',   location: 'us-east1',  key_ring_name: 'key-ring', name: 'crypto-key') do
      its('primary_state') { should eq "ENABLED" }
    end