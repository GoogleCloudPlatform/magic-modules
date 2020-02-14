### Test that there are no more than a specified number of instance groups available for the project

    describe google_compute_instance_groups(project: 'chef-inspec-gcp') do
      its('count') { should be <= 100}
    end

### Test that an expected instance_group is available for the project

    describe google_compute_instance_groups(project: 'chef-inspec-gcp', zone: 'europe-west2-a') do
      its('instance_group_names') { should include "my-instance-group-name" }
    end

### Test that a subset of all instance_groups matching "mig*" have size greater than zero

    google_compute_instance_groups(project: 'chef-inspec-gcp', zone: 'europe-west2-a').where(instance_group_name: /^mig/).instance_group_names.each do |instance_group_name|
      describe google_compute_instance_group(project: 'chef-inspec-gcp', zone: 'europe-west2-a', name: instance_group_name) do
        it { should exist }
        its('size') { should be > 0 }
      end
    end