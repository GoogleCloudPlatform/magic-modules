### Test that a GCP organization has the expected name

    describe google_organization(name: 'organizations/1234') do
      its('name') { should eq 'organizations/1234' }
    end

### Test that a GCP organization has the expected lifecycle state e.g. "ACTIVE"

    describe google_organization(display_name: 'google.com') do
      its('lifecycle_state') { should eq "ACTIVE" }
    end