### Test that there are no more than a specified number of service accounts for the project

    describe google_service_accounts(project: 'chef-inspec-gcp') do
      its('count') { should be <= 1000}
    end

### Test that an expected service account display name is available

    describe google_service_accounts(project: 'chef-inspec-gcp') do
      its('service_account_display_names'){ should include "gcp_sa_name" }
    end
    
### Test that an expected service account unique identifier is available

    describe google_service_accounts(project: 'chef-inspec-gcp') do
      its('service_account_ids'){ should include 12345678 }
    end    

### Test that a service account with expected name is available

    describe google_service_accounts(project: 'dummy-project') do
      its('service_account_names'){ should include "projects/dummy-project/serviceAccounts/dummy-acct@dummy-project.iam.gserviceaccount.com" }
    end

### Use filtering to retrieve a particular service account

    google_service_accounts(project: 'chef-inspec-gcp').where(service_account_display_names: /^dummyaccount/).service_account_names.each do |sa_name|
      describe google_service_account(name: sa_name) do
        it { should exist }
      end
    end