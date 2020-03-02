def deletion_protection_update(module, request, response):
    auth = GcpSession(module, 'compute')
    auth.post(
        ''.join(["https://www.googleapis.com/compute/v1/", "projects/{project}/zones/{zone}/instances/{name}/setDeletionProtection?deletionProtection={deletion_protection}"]).format(**module.params),
        {},
    )

def shielded_instance_config_update(module, request, response):
    auth = GcpSession(module, 'compute')
    auth.post(
        ''.join(["https://www.googleapis.com/compute/v1/", "projects/{project}/zones/{zone}/instances/{name}/updateShieldedInstanceConfig"]).format(**module.params),
        {
            u'enableSecureBoot': navigate_hash(module.params, ['shielded_instance_config', 'enable_secure_boot']),
            u'enableVtpm': navigate_hash(module.params, ['shielded_instance_config', 'enable_vtpm']),
            u'enableIntegrityMonitoring': navigate_hash(module.params, ['shielded_instance_config', 'enable_integrity_monitoring']),
        }
    )
