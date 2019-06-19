class Kubectl(object):
    def initialize(self, module):
        self.module = module

    """
    Writes a kubectl config file
    kubectl_path must be set or this will fail.
    """
    def write_file(self):
        with open(module.params['kubectl_path'], 'w') as f:
            f.write(yaml.dump(self._contents()))

    """
    Returns the contents of a kubectl file
    """
    def _contents(self):
        token = self._auth_token()
        endpoint = "https://{}".format(fetch["endpoint"])
        context = module.params.get('kubectl_context')
        if not context:
            context = module.params['name']

        return {
          'apiVersion': 'v1',
          'clusters': [
            {
              'name': context,
              'cluster': {
                'certificate-authority-data':
                  self.fetch['masterAuth']['clusterCaCertificate'],
                'server': endpoint,
              }
            }
          ],
          'contexts': [
            {
              'name': context,
              'context': {
                'cluster': context,
                'user': context
              }
            }
          ],
          'current-context': context,
          'kind': 'Config',
          'preferences': {},
          'users': [
            {
              'name': context,
              'user': {
                'auth-provider': {
                  'config': {
                    'access-token': token,
                    'cmd-args': 'config config-helper --format=json',
                    'cmd-path': '/usr/lib64/google-cloud-sdk/bin/gcloud',
                    'expiry-key': '{.credential.token_expiry}',
                    'token-key': '{.credential.access_token}'
                  },
                  'name': 'gcp'
                },
                'username': self.fetch['masterAuth']['username'],
                'password': self.fetch['masterAuth']['password']
              }
            }
          ]
        }

    """
    Returns the auth token used in kubectl
    This also sets the 'fetch' variable used in creating the kubectl
    """
    def _auth_token(self):
        auth = GcpSession(module, 'auth')
        response = auth.get(self_link(module))
        self.fetch = response.json()
        return response.request.headers.authorization.split(' ')[1]
