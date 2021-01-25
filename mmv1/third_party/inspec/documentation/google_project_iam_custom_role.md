### Test that a GCP project IAM custom role has the expected stage in the launch lifecycle

    describe google_project_iam_custom_role(project: 'chef-inspec-gcp', name: 'chef-inspec-gcp-role-abcd') do
      its('stage') { should eq "GA" }
    end

### Test that a GCP project IAM custom role has the expected included permissions

    describe google_project_iam_custom_role(project: 'chef-inspec-gcp', name: 'chef-inspec-gcp-role-abcd') do
      its('included_permissions') { should eq ["iam.roles.list"] }
    end