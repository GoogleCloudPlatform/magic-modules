email, ok := d.GetOk("email")
if !ok {
return fmt.Errorf("error registering ACME account, email address is required")
}
privateKeyPem, _ := d.GetOk("private_key_pem")
basePath, err := tpgresource.ReplaceVars(d, config, "{{PublicCABasePath}}")
if err != nil {
return err
}
isStagingEnv := strings.Contains(basePath, "preprod-")
eabKeyId, ok := d.GetOk("eab_key_id")
if !ok {
return fmt.Errorf("error registering ACME account, registration server is required")
}
eabHmacKeyUrlEncoded, ok := d.GetOk("eab_hmac_key")
if !ok {
return fmt.Errorf("error registering ACME account, registration server is required")
}

log.Printf("[DEBUG] Registering ACME account")
accountUri, err := createNewAccountUsingEab(email.(string), isStagingEnv, privateKeyPem.(string), eabKeyId.(string), eabHmacKeyUrlEncoded.(string))
if err != nil {
return fmt.Errorf("couldn't register the account: %s", err)
}

// Store the ID now
id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/global/acmeRegistrations/{{name}}")
if err != nil {
return fmt.Errorf("error constructing id: %s", err)
}
d.SetId(id)
log.Printf("[DEBUG] Finished registering ACME account: %s", accountUri)
d.Set("account_uri", accountUri)
return nil