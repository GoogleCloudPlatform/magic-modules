def resource_to_create(module):
    instance = resource_to_request(module)
    return {
        'instanceId': module.params['name'],
        'instance': instance
    }

def encode_request(request, module):
    del request['name']
    return request

def decode_response(response, module):
    response['name'] = response['name'].split('/')[-1]
    return response
