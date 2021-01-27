### Test that a GCP compute zone exists

    describe google_compute_zone(project: 'chef-inspec-gcp',  zone: 'us-east1-b') do
      it { should exist }
    end

### Test that a GCP compute zone is in the expected state

    describe google_compute_zone(project: 'chef-inspec-gcp',  zone: 'us-east1-b') do
      its('status') { should eq 'UP' }
      # or equivalently
      it { should be_up }
    end

### Test that a GCP compute zone has an expected CPU platform

    describe google_compute_zone(project: 'chef-inspec-gcp',  zone: 'us-east1-b') do
      its('available_cpu_platforms') { should include "Intel Skylake" }
    end