### Test that there are no more than a specified number of keys in the key ring

    describe google_kms_crypto_keys(project: 'chef-inspec-gcp',   location: 'us-east1',  key_ring_name: 'key-ring') do
      its('count') { should be <= 100}
    end

### Test that an expected key name is present in the key ring 

    describe google_kms_crypto_keys(project: 'chef-inspec-gcp',   location: 'us-east1',  key_ring_name: 'key-ring') do
      its('crypto_key_names') { should include "my-crypto-key-name" }
    end