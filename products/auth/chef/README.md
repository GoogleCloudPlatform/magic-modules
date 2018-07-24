# Google Authentication Chef Cookbook

## Module Description

This cookbook provides the resources to authenticate with Google Cloud Platform.

When executing operations on Google Cloud Platform, e.g. creating a virtual
machine, a SQL database, etc., you need to be authenticated to be able to carry
on with the request.

All Google Cloud Platform cookbooks use an unified authentication mechanism,
provided by this cookbook.

## Example

```chef
gauth_credential 'mycred' do
  action :serviceaccount
  path '/home/nelsonjr/my_account.json'
  scopes [
    'https://www.googleapis.com/auth/compute'
  ]
end
```

## Setup

To install this cookbook using `knife` tool:

    knife cookbook site install google-gauth
    knife cookbook site download google-gauth

_Note: Google Cloud Platform cookbooks that require authentication will
automatically install this cookbook, as it will be listed in their
dependencies._

## Platforms

### Supported Operating Systems

This cookbook was tested on the following operating systems:

* RedHat 6, 7
* CentOS 6, 7
* Debian 7, 8
* Ubuntu 12.04, 14.04, 16.04, 16.10
* SLES 11-sp4, 12-sp2
* openSUSE 13
* Windows Server 2008 R2, 2012 R2, 2012 R2 Core, 2016 R2, 2016 R2 Core

## About Service Accounts

This cookbook uses [service accounts][doc-accounts] to authenticate with Google
Cloud Platform. Google Cloud Platform project administrators manage service
accounts.  They can create, modify and delete accounts and grant account
specific privileges on the projects. Those privileges will be used by Chef to
carry on the operations on behalf of the user.

### Getting a Service Account key

The provider uses the JSON version of the service account key file. When in the
[IAM & Admin][iam-admin] section of the [Developer Console][console] the
administrator can retrieve a key file. Select the JSON as the key format.

The file you download is the one provided in the `path` property.

## Providers

### `gauth_credential`

#### `provider`

- `serviceaccount` **[preferred]**
	This is the preferred method of specifying credentials, because it does not
	rely on any pre-existing system configuration that Chef can't track. You'll
	need a credential file, but you can easily manage that with a `file do .. end`
	resource.

- `defaultuseraccount`
  If you have [`gcloud`][gcloud] setup you can piggyback on the account
  currently set as active for the user running Chef.

	_Warning! Please be aware that the account is subject to whichever account is
	set as active on `gcloud` tool, so the results will not be always consistent
	if they change under Chef by the user._

#### `scopes`

The scopes your authentication request will be limited to. When executing
actions against Google Cloud Platform you should choose the minimum amount of
privileges to carry on the operations to avoid accidentally affecting other
resources. For example if I want to manage virtual machines you should request
only "Compute R/W". That way you don't accidentally modify your DNS records.

Google's Chef cookbooks for Google Cloud Platform list the scopes you can use in
their Chef Supermarket documentation page. You can alternatively look at Google
Cloud Platform documentation for the product you're interacting with.

A few examples:

<table>
  <tr>
    <th>Product</th>
    <th colspan='2'>Scope</th>
  </tr>
  <tr>
    <td>Compute Engine (VMs, Disks, ...)</td>
    <td>Read Write</td>
    <td><code>https://www.googleapis.com/auth/compute</code></td>
  </tr>
  <tr>
    <td>Cloud SQL</td>
    <td>Read Write</td>
    <td><code>https://www.googleapis.com/auth/sqlservice.admin</code></td>
  </tr>
  <tr>
    <td rowspan='2'>Cloud DNS</td>
    <td>Read Only</td>
    <td><code>https://www.googleapis.com/auth/ndev.clouddns.readonly</code></td>
  </tr>
  <tr>
    <td>Read Write</td>
    <td><code>https://www.googleapis.com/auth/ndev.clouddns.readwrite</code></td>
  </tr>
</table>


#### `path`

If you specify the `serviceaccount` provider this property points to an absolute
path of the service account file (in JSON format).


[gcloud]: https://cloud.google.com/sdk
[console]: https://cloud.google.com/console
[doc-accounts]: https://cloud.google.com/compute/docs/access/service-accounts
[iam-admin]: https://console.cloud.google.com/iam-admin/serviceaccounts/project
