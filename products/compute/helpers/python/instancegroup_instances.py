<%
# Instance Logic handling for Instance Groups
# This code handles the adding + removing of instances from an instance group.
# It should be run after all normal create/update/delete logic.

-%>
class InstanceLogic(object):
    def __init__(self, module):
        self.module = module
        self.current_instances = self.list_instances()
        self.module_instances = []

        # Transform module list of instances (dicts of instance responses) into a list of selfLinks.
        instances = self.module.params.get('instances')
        if instances:
            for instance in instances:
                self.module_instances.append(replace_resource_dict(instance, 'selfLink'))

    def run(self):
        # Find all instances to add and add them
        instances_to_add = list(set(self.module_instances) - set(self.current_instances))
        if instances_to_add:
            self.add_instances(instances_to_add)

        # Find all instances to remove and remove them
        instances_to_remove = list(set(self.current_instances) - set(self.module_instances))
        if instances_to_remove:
            self.remove_instances(instances_to_remove)

    def list_instances(self):
        auth = GcpSession(self.module, 'compute')
        response = return_if_object(self.module, auth.post(self._list_instances_url(), {'instanceState': 'ALL'}),
                                    'compute#instanceGroupsListInstances')

        # Transform instance list into a list of selfLinks for diffing with module parameters
        instances = []
        for instance in response.get('items', []):
            instances.append(instance['instance'])
        return instances

    def add_instances(self, instances):
        auth = GcpSession(self.module, 'compute')
        wait_for_operation(self.module, auth.post(self._add_instances_url(), self._build_request(instances)))

    def remove_instances(self, instances):
        auth = GcpSession(self.module, 'compute')
        wait_for_operation(self.module, auth.post(self._remove_instances_url(), self._build_request(instances)))

    def _list_instances_url(self):
        return "https://www.googleapis.com/compute/v1/projects/{project}/zones/{zone}/instanceGroups/{name}/listInstances".format(**self.module.params)

    def _remove_instances_url(self):
        return "https://www.googleapis.com/compute/v1/projects/{project}/zones/{zone}/instanceGroups/{name}/removeInstances".format(**self.module.params)

    def _add_instances_url(self):
        return "https://www.googleapis.com/compute/v1/projects/{project}/zones/{zone}/instanceGroups/{name}/addInstances".format(**self.module.params)

    def _build_request(self, instances):
        request = {
            'instances': []
        }
        for instance in instances:
            request['instances'].append({'instance': instance})
        return request
