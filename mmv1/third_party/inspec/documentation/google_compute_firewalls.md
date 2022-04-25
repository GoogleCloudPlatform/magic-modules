### Test that there are no more than a specified number of firewalls available for the project

    describe google_compute_firewalls(project: 'chef-inspec-gcp') do
      its('count') { should be <= 100}
    end

### Test that an expected firewall is available for the project

    describe google_compute_firewalls(project: 'chef-inspec-gcp') do
      its('firewall_names') { should include "my-app-firewall-rule" }
    end

### Test that a particular named rule does not exist

    describe google_compute_firewalls(project: 'chef-inspec-gcp') do
      its('firewall_names') { should_not include "default-allow-ssh" }
    end

### Test there are no firewalls for the "INGRESS" direction

    describe google_compute_firewalls(project: 'chef-inspec-gcp').where(firewall_direction: 'INGRESS') do
      it { should_not exist }
    end