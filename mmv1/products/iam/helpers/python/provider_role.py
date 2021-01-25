def resource_to_create(module):
    role = resource_to_request(module)
    del role['name']
    return {
        'roleId': module.params['name'],
        'role': role
    }

def decode_response(response, module):
    if 'name' in response:
        response['name'] = response['name'].split('/')[-1]
    return response
