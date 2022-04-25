### Test that a GCP compute instance does not exist

    describe google_compute_instance(project: 'chef-inspec-gcp',  zone: 'us-east1-b', name: 'inspec-test-vm-not-there') do
      it { should_not exist }
    end

### Test that a GCP compute instance is in the expected state ([explore possible states here](https://cloud.google.com/compute/docs/instances/checking-instance-status))

    describe google_compute_instance(project: 'chef-inspec-gcp',  zone: 'us-east1-b', name: 'inspec-test-vm') do
      its('status') { should eq 'RUNNING' }
    end

### Test that a GCP compute instance is the expected size

    describe google_compute_instance(project: 'chef-inspec-gcp',  zone: 'us-east1-b', name: 'inspec-test-vm') do
      its('machine_type') { should match "f1-micro" }
    end

### Test that a GCP compute instance has the expected CPU platform

    describe google_compute_instance(project: 'chef-inspec-gcp',  zone: 'us-east1-b', name: 'inspec-test-vm') do
      its('cpu_platform') { should match "Intel" }
    end

### Test that a GCP compute instance has the expected number of attached disks

    describe google_compute_instance(project: 'chef-inspec-gcp',  zone: 'us-east1-b', name: 'inspec-test-vm') do
      its('disk_count'){should eq 2}
    end

### Test that a GCP compute instance has the expected number of attached network interfaces

    describe google_compute_instance(project: 'chef-inspec-gcp',  zone: 'us-east1-b', name: 'inspec-test-vm') do
      its('network_interfaces_count'){should eq 1}
    end

### Test that a GCP compute instance has the expected number of tags

    describe google_compute_instance(project: 'chef-inspec-gcp',  zone: 'us-east1-b', name: 'inspec-test-vm') do
      its('tag_count'){should eq 1}
    end

### Test that a GCP compute instance has a single public IP address

    describe google_compute_instance(project: 'chef-inspec-gcp',  zone: 'us-east1-b', name: 'inspec-test-vm') do
      its('first_network_interface_nat_ip_exists'){ should be true }
      its('first_network_interface_name'){ should eq "external-nat" }
      its('first_network_interface_type'){ should eq "one_to_one_nat" }
    end

### Test that a particular compute instance label key is present

    describe google_compute_instance(project: 'chef-inspec-gcp',  zone: 'us-east1-b', name: 'inspec-test-vm') do
      its('labels_keys') { should include 'my_favourite_label' }
    end

### Test that a particular compute instance label value is matching regexp 
    describe google_compute_instance(project: 'chef-inspec-gcp', zone:'us-east1-b', name:'inspec-test-vm').label_value_by_key('business-area') do
      it { should match '^(marketing|research)$' }
    end

### Test that a particular compute instance metadata key is present 
    describe google_compute_instance(project: 'chef-inspec-gcp', zone:'us-east1-b', name:'inspec-test-vm') do
      its('metadata_keys') { should include 'patching-type' }
    end

### Test that a particular compute instance metadata value is matching regexp 
    describe google_compute_instance(project: 'chef-inspec-gcp', zone:'us-east1-b', name:'inspec-test-vm').metadata_value_by_key('patching-window') do
      it { should match '^\d{1}-\d{2}$' }
    end