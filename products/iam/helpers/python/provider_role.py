def resource_to_create(module):
    role = resource_to_request(module)
    return {
        'roleId': module.params['name'],
        'role': role
    }
