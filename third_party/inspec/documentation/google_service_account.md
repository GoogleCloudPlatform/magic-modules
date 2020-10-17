### Test that a GCP project IAM service account has the expected unique identifier

    describe google_service_account(project: 'sample-project', name: 'sample-account@sample-project.iam.gserviceaccount.com') do
      its('unique_id') { should eq 12345678 }
    end

### Test that a GCP project IAM service account has the expected oauth2 client identifier

    describe google_service_account(project: 'sample-project', name: 'sample-account@sample-project.iam.gserviceaccount.com') do
      its('oauth2_client_id') { should eq 12345678 }
    end

### Test that a GCP project IAM service account does not have user managed keys

		describe google_service_account_keys(project: 'chef-gcp-inspec', service_account: "display-name@project-id.iam.gserviceaccount.com") do
		  its('key_types') { should_not include 'USER_MANAGED' }
    end