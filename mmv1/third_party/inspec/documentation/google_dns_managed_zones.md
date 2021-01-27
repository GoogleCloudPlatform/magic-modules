### Test that there are no more than a specified number of zones available for the project

    describe google_dns_managed_zones(project: 'chef-inspec-gcp') do
      its('count') { should be <= 100}
    end

### Test that an expected, named managed zone is available for the project

    describe google_dns_managed_zones(project: 'chef-inspec-gcp') do
      its('zone_names') { should include "zone-name" }
    end

### Test that a subset of all zones matching "myzone*" exist

    google_dns_managed_zones(project: 'chef-inspec-gcp').where(zone_name: /^myzone/).zone_names.each do |zone_name|
      describe google_dns_managed_zone(project: 'chef-inspec-gcp',  zone: zone_name) do
        it { should exist }
      end
    end