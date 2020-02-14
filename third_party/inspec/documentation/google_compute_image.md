### Test that a GCP compute image is in a particular status e.g. "READY" means available for use

    describe google_compute_image(project: 'chef-inspec-gcp', location: 'europe-west2', name: 'compute-address') do
      its('status') { should eq "READY" }
    end

### Test that a GCP compute image has the expected family

    describe google_compute_image(project: 'chef-inspec-gcp', name: 'ubuntu') do
      its('family') { should match "ubuntu" }
    end