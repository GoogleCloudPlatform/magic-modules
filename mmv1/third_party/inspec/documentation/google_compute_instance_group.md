### Test that a GCP compute instance group has the expected size

    describe google_compute_instance_group(project: 'chef-inspec-gcp', zone: 'europe-west2-a', name: 'gcp-inspec-test') do
      its('size') { should eq 2 }
    end

### Test that a GCP compute instance group has a port with supplied name and value

    describe google_compute_instance_group(project: 'chef-inspec-gcp', zone: 'europe-west2-a', name: 'gcp-inspec-test') do
      its('port_name') { should eq "http" }
      its('port_value') { should eq 80 }
    end
