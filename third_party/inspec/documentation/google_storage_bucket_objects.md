### Test that there are no more than a specified number of storage buckets for the project

    describe google_storage_bucket_objects(bucket: 'bucket-name') do
      its('count') { should be <= 100 }
    end


### Test that an expected named bucket is available

    describe google_storage_bucket_objects(bucket: 'bucket-name') do
      its('object_buckets'){ should include 'my_expected_bucket' }
    end

### Test that an expected named bucket is available

    describe google_storage_bucket_objects(bucket: 'bucket-name') do
      its('object_names'){ should include 'my_expected_object' }
    end
        
### Test a filtered group of bucket objects created within the last 24hrs

    describe google_storage_bucket_objects(bucket: 'bucket-name').where(object_created_time > Time.now - 60*60*24) do
      it { should exist }
    end