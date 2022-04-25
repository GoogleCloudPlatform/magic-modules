### Test that there are no more than a specified number of forwarding_rules available for the project and region

    describe google_compute_forwarding_rules(project: 'chef-inspec-gcp', region: 'europe-west2') do
      its('count') { should be <= 100}
    end

### Test that an expected forwarding_rule identifier is present in the project and region

    describe google_compute_forwarding_rules(project: 'chef-inspec-gcp', region: 'europe-west2') do
      its('forwarding_rule_ids') { should include 12345678975432 }
    end


### Test that an expected forwarding_rule name is available for the project and region

    describe google_compute_forwarding_rules(project: 'chef-inspec-gcp', region: 'europe-west2') do
      its('forwarding_rule_names') { should include "forwarding_rule-name" }
    end

### Test that an expected forwarding_rule network name is not present for the project and region

    describe google_compute_forwarding_rules(project: 'chef-inspec-gcp', region: 'europe-west2') do
      its('forwarding_rule_networks') { should not include "network-name" }
    end