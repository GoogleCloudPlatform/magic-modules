### Test that a GCP compute zone exists

    describe google_storage_bucket_object(bucket: 'bucket-buvsjjcndqz',  object: 'bucket-object-pmxbiikq') do
      it { should exist }
    end

### Test that a GCP storage bucket object has non-zero size

    describe google_storage_bucket_object(bucket: 'bucket-buvsjjcndqz',  object: 'bucket-object-pmxbiikq') do
      its('size') { should be > 0 }
    end

### Test that a GCP storage bucket object has the expected content type

    describe google_storage_bucket_object(bucket: 'bucket-buvsjjcndqz',  object: 'bucket-object-pmxbiikq') do
      its('content_type') { should eq "text/plain; charset=utf-8" }
    end


### Test that a GCP storage bucket object was created within a certain time period

    describe google_storage_bucket_object(bucket: 'bucket-buvsjjcndqz',  object: 'bucket-object-pmxbiikq') do
      its('time_created_date') { should be > Time.now - 365*60*60*24*10 }
    end
    
    
### Test that a GCP storage bucket object was last updated within a certain time period

    describe google_storage_bucket_object(bucket: 'bucket-buvsjjcndqz',  object: 'bucket-object-pmxbiikq') do
      its('updated_date') { should be > Time.now - 365*60*60*24*10 }
    end