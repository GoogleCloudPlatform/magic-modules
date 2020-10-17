def resource_to_create(module):
    instance = resource_to_request(module)
    if 'name' in instance:
        del instance['name']

    clusters = []
    if 'clusters' in instance:
        clusters = instance['clusters']
        del instance['clusters']

    return {
        'instanceId': module.params['name'].split('/')[-1],
        'instance': instance,
        'clusters': clusters
    }

def encode_request(request, module):
    if 'name' in request:
        del request['name']

    if 'clusters' in request:
        request['clusters'] = convert_clusters_to_map(request['clusters'])
    return request

def decode_response(response, module):
    if 'name' in response:
        response['name'] = response['name'].split('/')[-1]

    if 'clusters' in response:
        response['clusters'] = convert_map_to_clusters(response['clusters'])
    return response

def convert_clusters_to_map(clusters):
    cmap = {}
    for cluster in clusters:
        cmap[cluster['name']] = cluster
        del cmap[cluster['name']]['name']
    return cmap

def convert_map_to_clusters(clusters):
    carray = []
    for key, cluster in clusters.items():
        cluster['name'] = key.split('/')[-1]
        carray.append(cluster)
    return carray

def bigtable_async_url(module, extra_data=None):
    if extra_data is None:
        extra_data = {}
    location_name = module.params['clusters'][0]['location'].split('/')[-1]

    url = ('https://bigtableadmin.googleapis.com/v2/operations/projects/%s'
       '/instances/%s/locations/%s/operations/{op_id}' %
       (module.params['project'], module.params['name'], location_name))

    return url.format(**extra_data)
