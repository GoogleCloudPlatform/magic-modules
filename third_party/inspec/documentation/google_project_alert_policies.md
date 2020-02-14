### Test that there are no more than a specified number of project alert policies available for the project

    describe google_project_alert_policies(project: 'chef-inspec-gcp') do
      its('count') { should be <= 100}
    end

### Test that an expected policy name is available for the project

    describe google_project_alert_policies(project: 'chef-inspec-gcp') do
      its('policy_names') { should include 'projects/spaterson-project/alertPolicies/9271751234503117449' }
    end

### Test whether any expected policy display name is available for the project

    describe google_project_alert_policies(project: 'chef-inspec-gcp') do
      its('policy_display_names') { should_not include 'banned policy' }
    end

### Ensure no existing policies are inactive

    describe google_project_alert_policies(project: 'chef-inspec-gcp') do
      its('policy_enabled_states') { should_not include false }
    end