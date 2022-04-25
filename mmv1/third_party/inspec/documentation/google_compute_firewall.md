### Test that a GCP compute firewall allows SSH access on port 22

    describe google_compute_firewall(project: 'chef-inspec-gcp', name: 'firewall-rule') do
      its('allowed_ssh?')  { should be true }
    end

### Test that a GCP compute firewall does not allow HTTP access on port 80

    describe google_compute_firewall(project: 'chef-inspec-gcp', name: 'firewall-rule') do
      its('allowed_http?')  { should be false }
    end

### Test that a GCP compute firewall allows HTTPS access on port 443

    describe google_compute_firewall(project: 'chef-inspec-gcp', name: 'firewall-rule') do
      its('allowed_https?')  { should be true }
    end

### Test the direction of a GCP compute firewall e.g. "INGRESS" or "EGRESS"

    describe google_compute_firewall(project: 'chef-inspec-gcp', name: 'firewall-rule') do
      its('direction') { should eq "INGRESS" }
    end

### Test the source IP range list for the GCP compute firewall is not open to the world

    describe google_compute_firewall(project: 'chef-inspec-gcp', name: 'firewall-rule') do
      its('source_ranges') { should_not eq ["0.0.0.0/0"] }
      # or using helpers
      its('direction') { should eq "INGRESS" }
      it { should_not allow_ip_ranges ["0.0.0.0/0"] }
    end

### Test whether or not a port/protocol is defined for a given firewall rule

    describe google_compute_firewall(project: 'chef-inspec-gcp', name: 'firewall-rule') do
      # similar to the http  / ssh helpers above
      it { should_not allow_port_protocol("22", "tcp") }
      it { should allow_port_protocol("80", "tcp") }
    end

### Test whether firewall rule allows ingress/egress for specified tags 

    describe google_compute_firewall(project: 'chef-inspec-gcp', name: 'firewall-rule') do
      its('direction') { should eq "INGRESS" }
      it { should allow_source_tags ["allow-gcp-tag"] }
      it { should allow_target_tags ["allow-gcp-other-tag"] }
      # stricter
      it { should allow_source_tags_only ["allow-gcp-tag"] }
      it { should allow_target_tags_only ["allow-gcp-other-tag"] }
    end
