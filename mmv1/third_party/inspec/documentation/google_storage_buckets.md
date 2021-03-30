### Test that there are no more than a specified number of storage buckets for the project

    describe google_storage_buckets(project: 'chef-inspec-gcp') do
      its('count') { should be <= 100}
    end


### Test that an expected named bucket is available

    describe google_storage_buckets do
      its('bucket_names'){ should include "my_expected_bucket" }
    end
    
### Test that all buckets belong to the expected project number

    google_storage_buckets(project: 'chef-inspec-gcp').bucket_names.each do |bucket_name|
      describe google_storage_bucket(name: bucket_name) do
        it { should exist }
        its('project_number'){ should eq 1122334455 }
      end
    end