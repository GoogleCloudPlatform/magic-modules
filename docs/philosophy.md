# Module Philosophy

There are various processors out there that ingests a manifest and produce a
client library. You can get client libraries for Java, Ruby, .NET, etc. However
for deployment tools a library is not enough. _These tools do not care about an
abstraction layer as they are not a script language_.

## Convergence To Desired State

Our modules care about converging objects from any state into a final desired
state specified by the administrator. The administrator declares his intent in
a product specific DSL, and the infrastructure is supposed to make it happen.
For example consider that the administrator declared:

    gcompute_instance { 'my-middle-tier-vm':
      ensure => running,
      zone   => 'us-central1-a',
    }

If the machine is not running, the infrastructure is supposed to make it
happen: start if stopped, resume if paused, restore if archived. Also if the
machine is not on the us-central1-a zone, it should move it there. If something
described above cannot be accomplished the tool will fail with a clear and
actionable message of why it failed and what the administrator needs to do to
correct the issue.

## Idempotency

Contrary to script languages declarative tools (like Chef or Puppet) specify
intent and the system's final behavior. So that means if the system is already
in the desired state nothing should happen or be changed, or attempted to be
changed, _even if the final result is the same_.

Consider the following bash script:

    #!/bin/bash
    apt-get install apache2
    cp my-template.conf /etc/httpd/conf.d/25-mysite.conf
    chmod 644 /etc/httpd/conf.d/25-mysite.conf
    chattr -i /etc/httpd/conf.d/25-mysite.conf

In the simple script above we're installing Apache, copying our site
definitions and protecting the file. _Note that you cannot run that script
twice without it failing_ as chattr prevents the next cp to work. Also note that
if apache2 is already installed, although apt-get will not attempt to install
it again, the full apt-get install process will happen again. For example
catalog being updated or cleaned, block until another installation is in
progress, add installation to system log, etc.

So if you do the following in Puppet:

    package { 'apache2':
      ensure => installed,
    }

    file { '/etc/httpd/conf.d/25-mysite.conf':
      ensure => file,
      mode   => '0644',
    }

    chattr::attribute_add { '/etc/httpd/conf.d/25-mysite.conf':
      attribute => 'i',
      require   => File['/etc/httpd/conf.d/25-mysite.conf'],
    }

Puppet/Chef will not attempt to start installing apache2 if it is already there
(apt-get will not even be called). Similarly if the properties of the
configuration file are correct (mode, immutability bit, etc) Puppet/Chef will
not attempt to change anything. This leads to a cleaner execution
