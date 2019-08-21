def resource_to_create(module):
    role = resource_to_request(module)
    del role['name']
    return {
        'roleId': module.params['name'],
        'role': role
    }
