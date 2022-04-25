### Test that there are no more than a specified number of organizations available

    describe google_organizations do
      its('count') { should be <= 100}
    end

### Test that an expected organization name is available

    describe google_organizations do
      its('names') { should include "organization/1234" }
    end

### Test that an expected organization display name is available

    describe google_organizations do
      its('display_names') { should include "google.com" }
    end
    
### Test that all organizations are ACTIVE

    describe google_organizations do
      its('lifecycle_state'){ should eq 'ACTIVE' }
    end    

### Test that a particular subset of ACTIVE organizations with display name 'goog*' exist

    google_organizations.where(display_name: /^goog/, lifecycle_state: 'ACTIVE').names.each do |name|
      describe google_organization(name: name) do
        it { should exist }
      end
    end