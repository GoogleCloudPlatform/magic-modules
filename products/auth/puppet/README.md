# Google Authentication Puppet Module

## Module Description

This module provides the types to authenticate with Google Cloud Platform.

When executing operations on Google Cloud Platform, e.g. creating a virtual
machine, a SQL database, etc., you need to be authenticated to be able to carry
on with the request.

All Google Cloud Platform modules use an unified authentication mechanism,
provided by this module.

## Setup

To install this module on your Puppet Master (or Puppet Client/Agent), use the
Puppet module installer:

    puppet module install google-gauth

_Note: Google Cloud Platform modules that require authentication will
automatically install this module, as it will be listed in their dependencies._

### Required gem (libraries)

The authentication module depends on gems released by Google. Puppet does not
install gems automatically as part of installing modules. To install, run:

    /opt/puppetlabs/puppet/bin/gem install googleauth google-api-client

As everything related to system configuration, you can install the gem using
Puppet itself ;-)...

```puppet
package { [
    'googleauth',
    'google-api-client',
  ]:
    ensure   => present,
    provider => puppet_gem,
}
```

## Usage

```puppet
gauth_credential { 'mycred':
  path     => $cred_path, # e.g. '/home/nelsonjr/my_account.json'
  provider => serviceaccount,
  scopes   => [
    'https://www.googleapis.com/auth/ndev.clouddns.readwrite',
  ],
}
```

## About Service Accounts

This module uses [service accounts][doc-accounts] to authenticate with Google
Cloud Platform. Google Cloud Platform project administrators manage service
accounts.  They can create, modify and delete accounts and grant account
specific privileges on the projects. Those privileges will be used by Puppet to
carry on the operations on behalf of the user.

### Getting a Service Account key

This module uses the JSON version of the service account key file. When in the
[IAM & Admin][iam-admin] section of the [Developer Console][console] the
administrator can retrieve a key file. Select the JSON as the key format.

The file you download is the one provided in the `path` property.


## Reference

### Parameters

#### `gauth_credential`

##### `provider`

- `serviceaccount` **[preferred]**
	This is the preferred method of specifying credentials, because it does not
	rely on any pre-existing system configuration that Puppet can't track. You'll
	need a credential file, but you can easily manage that with a `file { }`
	resource.

- `defaultuseraccount`
  If you have [`gcloud`][gcloud] setup you can piggyback on the account
  currently set as active for the user running Puppet.

	_Warning! Please be aware that the account is subject to whichever account is
	set as active on `gcloud` tool, so the results will not be always consistent
	if they change under Puppet by the user._

##### `scopes`

The scopes your authentication request will be limited to. When executing
actions against Google Cloud Platform you should choose the minimum amount of
privileges to carry on the operations to avoid accidentally affecting other
resources. For example if I want to manage virtual machines you should request
only "Compute R/W". That way you don't accidentally modify your DNS records.

Google's Puppet modules for Google Cloud Platform list the scopes you can use in
their Puppet Forge documentation page. You can alternatively look at Google
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


##### `path`

If you specify the `serviceaccount` provider this property points to an absolute
path of the service account file (in JSON format).


### Functions

#### `gauth_credential_serviceaccount_for_function`

  Creates the credential token required by other client-side functions to
  retrieve Google Cloud Platform resource properties, e.g. Compute Instance IP
  address.

  This function uses a service account JSON file to provide credentials. A
  service account can be created an managed using Google Cloud IAM service.

  It is common to require dynamic data to be fetched and used by depedent
  resources, for example to lock down a database to a specific machines you
  require the machine's IP address. However the IP address is dynamic and cannot
  be determined upfront, so we need to provide the means to retrieve it
  programmatically. To allow such access GCP requires credentials. This function
  builds such credentials.

##### Arguments

  - `path`:
    the path for the service account credential JSON file

  - `scopes`:
    an array of scopes to allow access to. The service account privileges will
    be constrained to the list provided in this parameter.

##### Examples

```puppet
$fn_auth = gauth_credential_serviceaccount_for_function(
  '/root/my_account.json', ['https://www.googleapis.com/auth/sqlservice.admin']
)
```



[gcloud]: https://cloud.google.com/sdk
[console]: https://cloud.google.com/console
[doc-accounts]: https://cloud.google.com/compute/docs/access/service-accounts
[iam-admin]: https://console.cloud.google.com/iam-admin/serviceaccounts/project
