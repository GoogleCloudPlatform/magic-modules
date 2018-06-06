# Puppet
[Puppet][puppet] is a configuration management platform built by Puppet Labs.
Magic Modules auto-generates a series of modules to allow Puppet to create/alter
GCP resources.

Each GCP product (Compute Engine, Cloud Storage, etc) has a module. A Puppet
module is a collection of resources. Each resource corresponds one-to-one with
a GCP resource.

These modules are found on [Puppet Forge][puppet-forge]

## Setting up a Puppet Development Environment
### Installing Puppet
The current version of Puppet (as of 05/2018) is 5.3.

[Install instructions][puppet-install] for Linux are here.

#### Verification
```
  puppet --version
  > 5.3.0
```

### Placing Puppet Modules into proper location.
Puppet modules can be installed from Supermarket using the following command:
```
   puppet module install google-gcompute
```
In many cases, you'll want to be running these modules from source.

All Puppet modules are located in `~/.puppetlabs/etc/code/modules`. You'll have
to place each module inside this folder using the naming convention
`g<product name>`. No underscores are necessary.

Examples: `gcompute`, `gdns`, `gresourcemanager`

Symlinks can be placed here as well if you'd like to symlink these modules
to the output directories of Magic Modules.

#### Verification
```
   ls ~/.puppetlabs/etc/code/modules
   > gauth  gcompute  gcontainer  gdns  giam  glogging  gpubsub ...
```

Make sure that you have the gauth module as part of these.

### Installing Dependent Gems
(If you installed modules from Supermarket, ignore this).
Puppet has its own embedded Ruby with its own embedded gems. You'll need
to install all of the dependent gems for the GCP modules.

```
  sudo /opt/puppetlabs/puppet/bin/gem install googleauth google-api-client
```

#### Verification
```
  /opt/puppetlabs/puppet/bin/gem list | grep "google"
  > google-api-client (0.15.0)
  > googleauth (0.5.3)
```

## Running a Puppet Example
```
  # Project name is hardcoded. You probably aren't using our default project.
  sed -i 's/google.com:graphite-playground/your-project/g' path/to/example
  FACTER_cred_path=<path to service account> FACTER_project=<project name>
    puppet apply <path to example>
```
All environment variables have to be passed in with the `FACTER_` prefix.
You may have to change project names on examples.


[puppet]: https://www.puppet.com
[puppet-forge]: https://forge.puppet.com
[puppet-install]: https://puppet.com/docs/puppet/5.3/install_linux.html
