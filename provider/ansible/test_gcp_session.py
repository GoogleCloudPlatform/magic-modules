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
from contextlib import contextmanager
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
        raise kwargs["msg"]


class GcpSessionTestCase(unittest.TestCase):
    success_json = {"status": "SUCCESS"}
    user_agent = "Google-Ansible-MM-mock"
    url = "http://www.googleapis.com/compute/test_instance"

    @contextmanager
    def setup_auth(self):
        """
        This is a context manager that mocks out
        the google-auth library and uses the built-in
        AnonymousCredentials for sending requests.
        """
        with patch(
            "google.oauth2.service_account.Credentials.from_service_account_file"
        ) as mock:
            with patch.object(
                AnonymousCredentials, "with_scopes", create=True
            ) as mock2:
                creds = AnonymousCredentials()
                mock2.return_value = creds
                mock.return_value = creds
                yield

    @responses.activate
    def test_get(self):
        responses.add(responses.GET, self.url, status=200, json=self.success_json)

        with self.setup_auth():
            module = FakeModule(
                {
                    "scopes": "foo",
                    "service_account_file": "file_name",
                    "project": "test_project",
                    "auth_kind": "serviceaccount",
                }
            )

            session = GcpSession(module, "mock")
            resp = session.get(self.url)

            assert responses.calls[0].request.headers["User-Agent"] == self.user_agent
            assert resp.json() == self.success_json
            assert resp.status_code == 200

    @responses.activate
    def test_post(self):
        responses.add(responses.POST, self.url, status=200, json=self.success_json)

        with self.setup_auth():
            body = {"content": "some_content"}
            module = FakeModule(
                {
                    "scopes": "foo",
                    "service_account_file": "file_name",
                    "project": "test_project",
                    "auth_kind": "serviceaccount",
                }
            )

            session = GcpSession(module, "mock")
            resp = session.post(
                self.url, body=body, headers={"x-added-header": "my-header"}
            )

            # Ensure Google header added.
            assert responses.calls[0].request.headers["User-Agent"] == self.user_agent

            # Ensure all content was passed along.
            assert responses.calls[0].request.headers["x-added-header"] == "my-header"

            # Ensure proper request was made.
            assert resp.json() == self.success_json
            assert resp.status_code == 200

    @responses.activate
    def test_delete(self):
        responses.add(responses.DELETE, self.url, status=200, json=self.success_json)

        with self.setup_auth():
            body = {"content": "some_content"}
            module = FakeModule(
                {
                    "scopes": "foo",
                    "service_account_file": "file_name",
                    "project": "test_project",
                    "auth_kind": "serviceaccount",
                }
            )

            session = GcpSession(module, "mock")
            resp = session.delete(self.url)

            # Ensure Google header added.
            assert responses.calls[0].request.headers["User-Agent"] == self.user_agent

            # Ensure proper request was made.
            assert resp.json() == self.success_json
            assert resp.status_code == 200

    @responses.activate
    def test_put(self):
        responses.add(responses.PUT, self.url, status=200, json=self.success_json)

        with self.setup_auth():
            body = {"content": "some_content"}
            module = FakeModule(
                {
                    "scopes": "foo",
                    "service_account_file": "file_name",
                    "project": "test_project",
                    "auth_kind": "serviceaccount",
                }
            )

            session = GcpSession(module, "mock")
            resp = session.put(self.url, body={"foo": "bar"})

            # Ensure Google header added.
            assert responses.calls[0].request.headers["User-Agent"] == self.user_agent

            # Ensure proper request was made.
            assert resp.json() == self.success_json
            assert resp.status_code == 200

    @responses.activate
    def test_patch(self):
        responses.add(responses.PATCH, self.url, status=200, json=self.success_json)

        with self.setup_auth():
            body = {"content": "some_content"}
            module = FakeModule(
                {
                    "scopes": "foo",
                    "service_account_file": "file_name",
                    "project": "test_project",
                    "auth_kind": "serviceaccount",
                }
            )

            session = GcpSession(module, "mock")
            resp = session.patch(self.url, body={"foo": "bar"})

            # Ensure Google header added.
            assert responses.calls[0].request.headers["User-Agent"] == self.user_agent

            # Ensure proper request was made.
            assert resp.json() == self.success_json
            assert resp.status_code == 200
