## Examples

The following examples show how to use this InSpec audit resource.

### Test that a GCP alert policy is enabled 

    describe google_project_alert_policy(policy: 'spaterson', name: '9271751234503117449') do
      it { should be_enabled }
    end

### Test that a GCP compute alert policy display name is correct

    describe google_project_alert_policy(policy: 'spaterson-project', name: '9271751234503117449') do
      its('display_name') { should eq 'policy name' }
    end