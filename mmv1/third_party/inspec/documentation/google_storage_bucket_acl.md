### Test that a GCP storage bucket ACL exists

     describe google_storage_bucket_acl(bucket: 'bucket-buvsjjcndqz',  entity: 'user-object-viewer@spaterson-project.iam.gserviceaccount.com') do
      it { should exist }
    end

### Test that a GCP storage bucket ACL has the expected role (READER, WRITER or OWNER)

    describe google_storage_bucket_acl(bucket: 'bucket-buvsjjcndqz',  entity: 'user-object-viewer@spaterson-project.iam.gserviceaccount.com') do
      its('role') { should eq 'OWNER' }
    end