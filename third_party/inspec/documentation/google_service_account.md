### Test that a GCP project IAM service account has the expected unique identifier

    describe google_service_account(name: 'projects/sample-project/serviceAccounts/sample-account@sample-project.iam.gserviceaccount.com') do
      its('unique_id') { should eq 12345678 }
    end

### Test that a GCP project IAM service account has the expected oauth2 client identifier

    describe google_service_account(name: 'projects/sample-project/serviceAccounts/sample-account@sample-project.iam.gserviceaccount.com') do
      its('oauth2_client_id') { should eq 12345678 }
    end

### Test that a GCP project IAM service account does not have user managed keys

    describe google_service_account(name: 'projects/sample-project/serviceAccounts/sample-account@sample-project.iam.gserviceaccount.com') do
      it { should have_user_managed_keys }
    end