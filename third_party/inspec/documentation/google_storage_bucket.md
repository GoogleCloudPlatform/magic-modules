### Test that a GCP storage bucket is in the expected location

    describe google_storage_bucket(name: 'chef-inspec-gcp-storage-bucket-abcd') do
      its('location') { should eq "EUROPE-WEST2" }
    end

### Test that a GCP storage bucket has the expected project number

    describe google_storage_bucket(name: 'chef-inspec-gcp-storage-bucket-abcd') do
      its('project_number') {should eq 12345678 }
    end

### Test that a GCP storage bucket has the expected storage class

    describe google_storage_bucket(name: 'chef-inspec-gcp-storage-bucket-abcd') do
      its('storage_class') { should eq 'STANDARD' }
    end