### Test that there are no more than a specified number of zones available for the project

    describe google_compute_zones(project: 'chef-inspec-gcp') do
      its('count') { should be <= 100}
    end

### Test the exact number of zones in the project

    describe google_compute_zones(project: 'chef-inspec-gcp') do
      its('zone_ids.count') { should cmp 9 }
    end

### Test that an expected zone is available for the project

    describe google_compute_zones(project: 'chef-inspec-gcp') do
      its('zone_names') { should include "us-east1-b" }
    end

### Test whether any zones are in status "DOWN"

    describe google_compute_zones(project: 'chef-inspec-gcp') do
      its('zone_statuses') { should_not include "DOWN" }
    end

### Test that a subset of all zones matching "us*" are "UP"

    google_compute_zones(project: 'chef-inspec-gcp').where(zone_name: /^us/).zone_names.each do |zone_name|
      describe google_compute_zone(project: 'chef-inspec-gcp',  zone: zone_name) do
        it { should exist }
        its('status') { should eq 'UP' }
      end
    end