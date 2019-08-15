# -*- coding: utf-8 -*-
# (c) 2019, Google Inc.
#
# This file is part of Ansible
#
# Ansible is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Ansible is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with Ansible.  If not, see <http://www.gnu.org/licenses/>.

from pytest import importorskip
from units.compat import unittest
from units.compat.mock import patch
from ansible.module_utils.gcp_utils import GcpSession
import responses
import tempfile

importorskip("requests")
importorskip("google.auth")
importorskip("responses")

from google.auth.credentials import AnonymousCredentials

class FakeModule(object):
    def __init__(self, params):
        self.params = params

    def fail_json(self, **kwargs):
        raise kwargs['msg']


class GcpSessionTestCase(unittest.TestCase):

    @responses.activate
    def test_full_get(self):
        url = 'http://www.googleapis.com/compute/test_instance'
        responses.add(responses.GET, url,
                      status=200, json={'status': 'SUCCESS'})

        with patch('google.oauth2.service_account.Credentials.from_service_account_file') as mock:
            with patch.object (AnonymousCredentials, 'with_scopes', create=True) as mock2:
                creds = AnonymousCredentials()
                mock2.return_value = creds
                mock.return_value = creds

                module = FakeModule({ 'scopes': 'foo', 'service_account_file': 'file_name', 'project': 'test_project', 'auth_kind': 'serviceaccount'})

                session = GcpSession(module, 'mock')
                resp = session.get(url)

                assert resp.request.headers['User-Agent'] == 'Google-Ansible-MM-mock'
                assert resp.json() == {'status': 'SUCCESS'}
                assert resp.status_code == 200
