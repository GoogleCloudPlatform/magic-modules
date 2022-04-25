### Test that a GCP compute vpn_tunnel exists

    describe google_compute_vpn_tunnel(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-vpn-tunnel') do
      it { should exist }
    end

### Test when a GCP compute vpn_tunnel was created

    describe google_compute_vpn_tunnel(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-vpn-tunnel') do
      its('creation_timestamp_date') { should be > Time.now - 365*60*60*24*10 }
    end

### Test for an expected vpn_tunnel identifier 

    describe google_compute_vpn_tunnel(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-vpn-tunnel') do
      its('id') { should eq 12345567789 }
    end    

### Test that a vpn_tunnel peer address is as expected

    describe google_compute_vpn_tunnel(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-vpn-tunnel') do
      its('peer_ip') { should eq "123.123.123.123" }
    end  

### Test that a vpn_tunnel status is as expected

    describe google_compute_vpn_tunnel(project: 'chef-inspec-gcp', region: 'europe-west2', name: 'gcp-inspec-vpn_tunnel') do
      its('status') { should eq "ESTABLISHED" }
    end 