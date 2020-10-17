### Test that GCP compute project information exists

    describe google_compute_project_info(project: 'chef-inspec-gcp') do
      it { should exist }
    end

### Test that GCP compute project default service account is as expected

    describe google_compute_project_info(project: 'chef-inspec-gcp') do
      its('default_service_account') { should eq '12345-compute@developer.gserviceaccount.com' }
    end