### Test that there are no more than a specified number of keys for the service account

    describe google_service_account_keys(project: 'sample-project', service_account: 'sample-account@sample-project.iam.gserviceaccount.com') do
      its('count') { should be <= 1000}
    end
    
### Test that a service account with expected name is available

    describe google_service_account_keys(project: 'sample-project', service_account: 'sample-account@sample-project.iam.gserviceaccount.com') do
      its('key_names'){ should include "projects/sample-project/serviceAccounts/test-sa@sample-project.iam.gserviceaccount.com/keys/c6bd986da9fac6d71178db41d1741cbe751a5080" }
    end