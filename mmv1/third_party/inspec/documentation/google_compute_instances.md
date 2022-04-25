### Test that there are no more than a specified number of instances in the project and zone

    describe google_compute_instances(project: 'chef-inspec-gcp',  zone: 'europe-west2-a') do
      its('count') { should be <= 100}
    end

### Test the exact number of instances in the project and zone

    describe google_compute_instances(project: 'chef-inspec-gcp',  zone: 'europe-west2-a') do
      its('instance_ids.count') { should cmp 9 }
    end

### Test that an instance with a particular name exists in the project and zone

    describe google_compute_instances(project: 'chef-inspec-gcp',  zone: 'europe-west2-a') do
      its('instance_names') { should include "my-favourite-instance" }
    end