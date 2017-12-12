# Copyright: (c) 2018, Google Inc.
# GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)


class ModuleDocFragment(object):
        # GCP doc fragment.
        DOCUMENTATION = '''
options:
    state:
        description:
            - Whether the given zone should or should not be present.
        required: true
        choices: ["present", "absent"]
        default: "present"
    project:
        description:
            - The Google Cloud Platform project to use.
        default: null
    auth_kind:
        description:
            - The type of credential used.
        required: true
        choices: ["machineaccount", "serviceaccount", "application"]
    service_account_file:
        description:
            - The path of a Service Account JSON file if serviceaccount is selected as type.
    service_account_email:
        description:
            - An optional service account email address if machineaccount is selected
              and the user does not wish to use the default email.
    scopes:
      description:
          - Array of scopes to be used.
      required: true
notes:
  - For authentication, you can set service_account_file using the
    GCP_SERVICE_ACCOUNT_FILE env variable.
  - For authentication, you can set service_account_email using the
    GCP_SERVICE_ACCOUNT_EMAIL env variable.
  - For authentication, you can set auth_kind using the GCP_AUTH_KIND env
    variable.
  - For authentication, you can set scopes using the GCP_SCOPES env variable.
  - Environment variables values will only be used if the playbook values are
    not set.
'''
