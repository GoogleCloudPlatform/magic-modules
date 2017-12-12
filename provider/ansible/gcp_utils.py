# Copyright (c), Google Inc, 2017
# All rights reserved.
#
# Redistribution and use in source and binary forms, with or without modification,
# are permitted provided that the following conditions are met:
#
#    * Redistributions of source code must retain the above copyright
#      notice, this list of conditions and the following disclaimer.
#    * Redistributions in binary form must reproduce the above copyright notice,
#      this list of conditions and the following disclaimer in the documentation
#      and/or other materials provided with the distribution.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
# ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
# WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
# IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
# INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
# PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
# INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
# LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE
# USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
#

try:
    import requests
    HAS_REQUESTS = True
except ImportError:
    HAS_REQUESTS = False

try:
    import google.auth
    import google.auth.compute_engine
    from google.oauth2 import service_account
    from google.auth.transport.requests import AuthorizedSession
    HAS_GOOGLE_LIBRARIES = True
except ImportError:
    HAS_GOOGLE_LIBRARIES = False

from ansible.module_utils.basic import AnsibleModule
import os


def navigate_hash(source, path, default=None):
    key = path[0]
    path = path[1:]
    if key not in source:
        return default
    result = source[key]
    if path:
        return navigate_hash(result, path, default)
    else:
        return result


class GcpAuthentication(object):
    def __init__(self, module):
        self.module = module
        self._set_values
        self._validate()

    def session(self):
        return AuthorizedSession(
            self._credentials().with_scopes(self.module.params['scopes']))

    def _set_values(self):
        if not self.module.params['auth_kind']:
            self.module.params['auth_kind'] = os.getenv('GCP_AUTH_KIND')

        if not self.module.params['service_account_email']:
            self.module.params['service_account_email'] = os.getenv('GCP_SERVICE_ACCOUNT_EMAIL')

        if not self.module.params['service_account_file']:
            self.module.params['service_account_file'] = os.getenv('GCP_SERVICE_ACCOUNT_FILE')

        if not self.module.params['scopes']:
            self.module.params['scopes'] = os.getenv('GCP_SCOPES')

    def _validate(self):
        if not HAS_REQUESTS:
            self.module.fail_json(msg="Please install the requests library")

        if not HAS_GOOGLE_LIBRARIES:
            self.module.fail_json(msg="Please install the google-auth library")

        if self.module.params['service_account_email'] is not None and self.module.params['auth_kind'] != 'machineaccount':
            self.module.fail_json(
                msg="Service Acccount Email only works with Machine Account-based authentication"
            )

        if self.module.params['service_account_file'] is not None and self.module.params['auth_kind'] != 'serviceaccount':
            self.module.fail_json(
                msg="Service Acccount File only works with Service Account-based authentication"
            )

    def _credentials(self):
        cred_type = self.module.params['auth_kind']
        if cred_type == 'application':
            credentials, project_id = google.auth.default()
            return credentials
        elif cred_type == 'serviceaccount':
            return service_account.Credentials.from_service_account_file(
                self.module.params['service_account_file'])
        elif cred_type == 'machineaccount':
            return google.auth.compute_engine.Credentials(
                self.module.params['service_account_email'])
        else:
   
