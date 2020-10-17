### Test that there are no more than a specified number of vpn_tunnels available for the project and region

    describe google_compute_vpn_tunnels(project: 'chef-inspec-gcp', region: 'europe-west2') do
      its('count') { should be <= 100}
    end

### Test that an expected vpn_tunnel name is available for the project and region

    describe google_compute_vpn_tunnels(project: 'chef-inspec-gcp', region: 'europe-west2') do
      its('vpn_tunnel_names') { should include "vpn_tunnel-name" }
    end

### Test that an expected vpn_tunnel target_vpn_gateways name is not present for the project and region

    describe google_compute_vpn_tunnels(project: 'chef-inspec-gcp', region: 'europe-west2') do
      its('vpn_tunnel_target_vpn_gateways') { should not include "gateway-name" }
    end