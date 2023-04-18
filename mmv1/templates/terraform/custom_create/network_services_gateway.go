userAgent, err := generateUserAgentString(d, config.UserAgent)
if err != nil {
	return err
}

obj := make(map[string]interface{})
labelsProp, err := expandNetworkServicesGatewayLabels(d.Get("labels"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
	obj["labels"] = labelsProp
}
descriptionProp, err := expandNetworkServicesGatewayDescription(d.Get("description"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
	obj["description"] = descriptionProp
}
typeProp, err := expandNetworkServicesGatewayType(d.Get("type"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("type"); !isEmptyValue(reflect.ValueOf(typeProp)) && (ok || !reflect.DeepEqual(v, typeProp)) {
	obj["type"] = typeProp
}
portsProp, err := expandNetworkServicesGatewayPorts(d.Get("ports"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("ports"); !isEmptyValue(reflect.ValueOf(portsProp)) && (ok || !reflect.DeepEqual(v, portsProp)) {
	obj["ports"] = portsProp
}
scopeProp, err := expandNetworkServicesGatewayScope(d.Get("scope"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("scope"); !isEmptyValue(reflect.ValueOf(scopeProp)) && (ok || !reflect.DeepEqual(v, scopeProp)) {
	obj["scope"] = scopeProp
}
serverTlsPolicyProp, err := expandNetworkServicesGatewayServerTlsPolicy(d.Get("server_tls_policy"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("server_tls_policy"); !isEmptyValue(reflect.ValueOf(serverTlsPolicyProp)) && (ok || !reflect.DeepEqual(v, serverTlsPolicyProp)) {
	obj["serverTlsPolicy"] = serverTlsPolicyProp
}
addressesProp, err := expandNetworkServicesGatewayAddresses(d.Get("addresses"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("addresses"); !isEmptyValue(reflect.ValueOf(addressesProp)) && (ok || !reflect.DeepEqual(v, addressesProp)) {
	obj["addresses"] = addressesProp
}
subnetworkProp, err := expandNetworkServicesGatewaySubnetwork(d.Get("subnetwork"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("subnetwork"); !isEmptyValue(reflect.ValueOf(subnetworkProp)) && (ok || !reflect.DeepEqual(v, subnetworkProp)) {
	obj["subnetwork"] = subnetworkProp
}
networkProp, err := expandNetworkServicesGatewayNetwork(d.Get("network"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("network"); !isEmptyValue(reflect.ValueOf(networkProp)) && (ok || !reflect.DeepEqual(v, networkProp)) {
	obj["network"] = networkProp
}
gatewaySecurityPolicyProp, err := expandNetworkServicesGatewayGatewaySecurityPolicy(d.Get("gateway_security_policy"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("gateway_security_policy"); !isEmptyValue(reflect.ValueOf(gatewaySecurityPolicyProp)) && (ok || !reflect.DeepEqual(v, gatewaySecurityPolicyProp)) {
	obj["gatewaySecurityPolicy"] = gatewaySecurityPolicyProp
}
certificateUrlsProp, err := expandNetworkServicesGatewayCertificateUrls(d.Get("certificate_urls"), d, config)
if err != nil {
	return err
} else if v, ok := d.GetOkExists("certificate_urls"); !isEmptyValue(reflect.ValueOf(certificateUrlsProp)) && (ok || !reflect.DeepEqual(v, certificateUrlsProp)) {
	obj["certificateUrls"] = certificateUrlsProp
}

url, err := ReplaceVars(d, config, "{{NetworkServicesBasePath}}projects/{{project}}/locations/{{location}}/gateways?gatewayId={{name}}")
if err != nil {
	return err
}

log.Printf("[DEBUG] Creating new Gateway: %#v", obj)
billingProject := ""

project, err := getProject(d, config)
if err != nil {
	return fmt.Errorf("Error fetching project for Gateway: %s", err)
}
billingProject = project

// err == nil indicates that the billing_project value was found
if bp, err := getBillingProject(d, config); err == nil {
	billingProject = bp
}

// resourceNetworkServicesGatewayCreate method is pretty much the same as the auto generated via MMv1 except for the following piece.
// Both the request and its operation must be retried at once instead of retrying only one of them.
// This was required because when the gateway of type SECURE_WEB_GATEWAY is being created it also creates a "swg-autogen-router":
// 1- if there is more than one gateway being created at the same time,
// they all try to create them "swg-autogen-router" at the same time but only one is necessary under the same netowork;
// 2- it sends the request of gateway creation to the api but the operation may fail since there might be a conflict of the swg-autogen-router creation.
// 3- all the operation must be retried so the subsequent request will reuse the swg-autogen-router instead of trying to create a new one.

// BEGIN
maxRetries := 3
for i := 1; i <= maxRetries; i++ {
	res, err := SendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Gateway: %s", err)
	}

	// Store the ID now
	id, err := ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/gateways/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = NetworkServicesOperationWaitTime(
		config, res, project, "Creating Gateway", userAgent,
		d.Timeout(schema.TimeoutCreate))

	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		if i < maxRetries {
			time.Sleep(5 * time.Second)
			continue
		}
		return fmt.Errorf("Error waiting to create Gateway: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Gateway %q: %#v", d.Id(), res)
	// all good
	break
}
// END

return resourceNetworkServicesGatewayRead(d, meta)