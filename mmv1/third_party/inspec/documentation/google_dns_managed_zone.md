### Test that a GCP compute zone exists

    describe google_dns_managed_zone(project: 'chef-inspec-gcp',  zone: 'zone-name') do
      it { should exist }
    end

### Test that a GCP DNS managed zone has the expected DNS name

    describe google_dns_managed_zone(project: 'chef-inspec-gcp',  zone: 'zone-name') do
      its('dns_name') { should match 'mydomain.com' }
    end

### Test that a GCP DNS managed zone has expected name server

    describe google_dns_managed_zone(project: 'chef-inspec-gcp',  zone: 'zone-name') do
      its('name_servers') { should include 'ns-cloud-d1.googledomains.com.' }
    end