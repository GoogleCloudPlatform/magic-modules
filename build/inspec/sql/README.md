# Google Cloud SQL Chef Cookbook

This cookbook provides the built-in types and services for Chef to manage
Google Cloud Compute resources, as native Chef types.

## Requirements

### Platforms

#### Supported Operating Systems

This cookbook was tested on the following operating systems:

* RedHat 6, 7
* CentOS 6, 7
* Debian 7, 8
* Ubuntu 12.04, 14.04, 16.04, 16.10
* SLES 11-sp4, 12-sp2
* openSUSE 13
* Windows Server 2008 R2, 2012 R2, 2012 R2 Core, 2016 R2, 2016 R2 Core

## Example

```ruby
gauth_credential 'mycred' do
  action :serviceaccount
  path ENV['CRED_PATH'] # e.g. '/path/to/my_account.json'
  scopes [
    'https://www.googleapis.com/auth/sqlservice.admin'
  ]
end

gsql_instance  "sql-test-#{ENV['sql_instance_suffix']}" do
  action :create
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

gsql_database 'webstore' do
  action :create
  charset 'utf8'
  instance "sql-test-#{ENV['sql_instance_suffix']}"
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end
```

## Credentials

All Google Cloud Platform cookbooks use an unified authentication mechanism,
provided by the `google-gauth` cookbook. Don't worry, it is automatically
installed when you install this module.

### Example

```ruby
gauth_credential 'mycred' do
  action :serviceaccount
  path ENV['CRED_PATH'] # e.g. '/path/to/my_account.json'
  scopes [
    'https://www.googleapis.com/auth/sqlservice.admin'
  ]
end

```

For complete details of the authentication cookbook, visit the
[google-gauth][] cookbook documentation.

## Resources

* [`gsql_instance`](#gsql_instance) -
    Represents a Cloud SQL instance. Cloud SQL instances are SQL databases
    hosted in Google's cloud. The Instances resource provides methods for
    common configuration and management tasks.
* [`gsql_database`](#gsql_database) -
    Represents a SQL database inside the Cloud SQL instance, hosted in
    Google's cloud.
* [`gsql_user`](#gsql_user) -
    The Users resource represents a database user in a Cloud SQL instance.


### gsql_instance
Represents a Cloud SQL instance. Cloud SQL instances are SQL databases
hosted in Google's cloud. The Instances resource provides methods for
common configuration and management tasks.


#### Example

```ruby
gsql_instance "sql-test-#{ENV['sql_instance_suffix']}" do
  action :create
  database_version 'MYSQL_5_7'
  settings({
    tier: 'db-n1-standard-1',
    ip_configuration:  {
      authorized_networks: [
        # The ACL below is for example only. (do NOT use in production as-is)
        {
          name: 'google dns server',
          value: '8.8.8.8/32'
        }
      ]
    }
  })
  region 'us-central1'
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

```

#### Reference

```ruby
gsql_instance 'id-for-resource' do
  backend_type          'FIRST_GEN', 'SECOND_GEN' or 'EXTERNAL'
  connection_name       string
  database_version      'MYSQL_5_5', 'MYSQL_5_6', 'MYSQL_5_7' or 'POSTGRES_9_6'
  failover_replica      {
    available boolean,
    name      string,
  }
  instance_type         'CLOUD_SQL_INSTANCE', 'ON_PREMISES_INSTANCE' or 'READ_REPLICA_INSTANCE'
  ip_addresses          [
    {
      ip_address     string,
      time_to_retire time,
      type           'PRIMARY' or 'OUTGOING',
    },
    ...
  ]
  ipv6_address          string
  master_instance_name  string
  max_disk_size         integer
  name                  string
  region                string
  replica_configuration {
    failover_target               boolean,
    mysql_replica_configuration   {
      ca_certificate            string,
      client_certificate        string,
      client_key                string,
      connect_retry_interval    integer,
      dump_file_path            string,
      master_heartbeat_period   integer,
      password                  string,
      ssl_cipher                string,
      username                  string,
      verify_server_certificate boolean,
    },
    replica_names                 [
      string,
      ...
    ],
    service_account_email_address string,
  }
  settings              {
    ip_configuration {
      authorized_networks [
        {
          expiration_time time,
          name            string,
          value           string,
        },
        ...
      ],
      ipv4_enabled        boolean,
      require_ssl         boolean,
    },
    settings_version integer,
    tier             string,
  }
  project               string
  credential            reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gsql_instance` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gsql_instance` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `backend_type` -
  * FIRST_GEN: First Generation instance. MySQL only.
  * SECOND_GEN: Second Generation instance or PostgreSQL instance.
  * EXTERNAL: A database server that is not managed by Google.

* `connection_name` -
  Connection name of the Cloud SQL instance used in connection strings.

* `database_version` -
  The database engine type and version. For First Generation instances,
  can be MYSQL_5_5, or MYSQL_5_6. For Second Generation instances, can
  be MYSQL_5_6 or MYSQL_5_7. Defaults to MYSQL_5_6.
  PostgreSQL instances: POSTGRES_9_6
  The databaseVersion property can not be changed after instance
  creation.

* `failover_replica` -
  The name and status of the failover replica. This property is
  applicable only to Second Generation instances.

* `failover_replica/available`
  Output only. The availability status of the failover replica. A false status
  indicates that the failover replica is out of sync. The master
  can only failover to the falover replica when the status is true.

* `failover_replica/name`
  The name of the failover replica. If specified at instance
  creation, a failover replica is created for the instance. The name
  doesn't include the project ID. This property is applicable only
  to Second Generation instances.

* `instance_type` -
  The instance type. This can be one of the following.
  * CLOUD_SQL_INSTANCE: A Cloud SQL instance that is not replicating
  from a master.
  * ON_PREMISES_INSTANCE: An instance running on the customer's
  premises.
  * READ_REPLICA_INSTANCE: A Cloud SQL instance configured as a
  read-replica.

* `ip_addresses` -
  Output only. The assigned IP addresses for the instance.

* `ip_addresses[]/ip_address`
  The IP address assigned.

* `ip_addresses[]/time_to_retire`
  The due time for this IP to be retired in RFC 3339 format, for
  example 2012-11-15T16:19:00.094Z. This field is only available
  when the IP is scheduled to be retired.

* `ip_addresses[]/type`
  The type of this IP address. A PRIMARY address is an address
  that can accept incoming connections. An OUTGOING address is the
  source address of connections originating from the instance, if
  supported.

* `ipv6_address` -
  The IPv6 address assigned to the instance. This property is applicable
  only to First Generation instances.

* `master_instance_name` -
  The name of the instance which will act as master in the replication
  setup.

* `max_disk_size` -
  The maximum disk size of the instance in bytes.

* `name` -
  Required. Name of the Cloud SQL instance. This does not include the project
  ID.

* `region` -
  The geographical region. Defaults to us-central or us-central1
  depending on the instance type (First Generation or Second
  Generation/PostgreSQL).

* `replica_configuration` -
  Configuration specific to failover replicas and read replicas.

* `replica_configuration/failover_target`
  Specifies if the replica is the failover target. If the field is
  set to true the replica will be designated as a failover replica.
  In case the master instance fails, the replica instance will be
  promoted as the new master instance.
  Only one replica can be specified as failover target, and the
  replica has to be in different zone with the master instance.

* `replica_configuration/mysql_replica_configuration`
  MySQL specific configuration when replicating from a MySQL
  on-premises master. Replication configuration information such as
  the username, password, certificates, and keys are not stored in
  the instance metadata.  The configuration information is used
  only to set up the replication connection and is stored by MySQL
  in a file named master.info in the data directory.

* `replica_configuration/mysql_replica_configuration/ca_certificate`
  PEM representation of the trusted CA's x509 certificate.

* `replica_configuration/mysql_replica_configuration/client_certificate`
  PEM representation of the slave's x509 certificate

* `replica_configuration/mysql_replica_configuration/client_key`
  PEM representation of the slave's private key. The
  corresponsing public key is encoded in the client's asf asd
  certificate.

* `replica_configuration/mysql_replica_configuration/connect_retry_interval`
  Seconds to wait between connect retries. MySQL's default is 60
  seconds.

* `replica_configuration/mysql_replica_configuration/dump_file_path`
  Path to a SQL dump file in Google Cloud Storage from which the
  slave instance is to be created. The URI is in the form
  gs://bucketName/fileName. Compressed gzip files (.gz) are
  also supported. Dumps should have the binlog co-ordinates from
  which replication should begin. This can be accomplished by
  setting --master-data to 1 when using mysqldump.

* `replica_configuration/mysql_replica_configuration/master_heartbeat_period`
  Interval in milliseconds between replication heartbeats.

* `replica_configuration/mysql_replica_configuration/password`
  The password for the replication connection.

* `replica_configuration/mysql_replica_configuration/ssl_cipher`
  A list of permissible ciphers to use for SSL encryption.

* `replica_configuration/mysql_replica_configuration/username`
  The username for the replication connection.

* `replica_configuration/mysql_replica_configuration/verify_server_certificate`
  Whether or not to check the master's Common Name value in the
  certificate that it sends during the SSL handshake.

* `replica_configuration/replica_names`
  The replicas of the instance.

* `replica_configuration/service_account_email_address`
  The service account email address assigned to the instance. This
  property is applicable only to Second Generation instances.

* `settings` -
  The user settings.

* `settings/ip_configuration`
  The settings for IP Management. This allows to enable or disable
  the instance IP and manage which external networks can connect to
  the instance. The IPv4 address cannot be disabled for Second
  Generation instances.

* `settings/ip_configuration/ipv4_enabled`
  Whether the instance should be assigned an IP address or not.

* `settings/ip_configuration/authorized_networks`
  The list of external networks that are allowed to connect to
  the instance using the IP. In CIDR notation, also known as
  'slash' notation (e.g. 192.168.100.0/24).

* `settings/ip_configuration/authorized_networks[]/expiration_time`
  The time when this access control entry expires in RFC
  3339 format, for example 2012-11-15T16:19:00.094Z.

* `settings/ip_configuration/authorized_networks[]/name`
  An optional label to identify this entry.

* `settings/ip_configuration/authorized_networks[]/value`
  The whitelisted value for the access control list. For
  example, to grant access to a client from an external IP
  (IPv4 or IPv6) address or subnet, use that address or
  subnet here.

* `settings/ip_configuration/require_ssl`
  Whether the mysqld should default to 'REQUIRE X509' for
  users connecting over IP.

* `settings/tier`
  The tier or machine type for this instance, for
  example db-n1-standard-1. For MySQL instances, this field
  determines whether the instance is Second Generation (recommended)
  or First Generation.

* `settings/settings_version`
  Output only. The version of instance settings. This is a required field for
  update method to make sure concurrent updates are handled properly.
  During update, use the most recent settingsVersion value for this
  instance and do not try to update this value.

#### Label
Set the `i_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gsql_database
Represents a SQL database inside the Cloud SQL instance, hosted in
Google's cloud.


#### Example

```ruby
# Tip: Remember to define gsql_instance to match the 'instance' property.
gsql_database 'webstore' do
  action :create
  charset 'utf8'
  instance "sql-test-#{ENV['sql_instance_suffix']}"
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

```

#### Reference

```ruby
gsql_database 'id-for-resource' do
  charset    string
  collation  string
  instance   reference to gsql_instance
  name       string
  project    string
  credential reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gsql_database` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gsql_database` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `charset` -
  The MySQL charset value.

* `collation` -
  The MySQL collation value.

* `name` -
  The name of the database in the Cloud SQL instance.
  This does not include the project ID or instance name.

* `instance` -
  Required. The name of the Cloud SQL instance. This does not include the project
  ID.

#### Label
Set the `d_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

### gsql_user
The Users resource represents a database user in a Cloud SQL instance.


#### Example

```ruby
# Tip: Remember to define gsql_instance to match the 'instance' property.
gsql_user 'john.doe' do
  action :create
  password 'secret-password'
  host '10.1.2.3'
  instance "sql-test-#{ENV['sql_instance_suffix']}"
  project ENV['PROJECT'] # ex: 'my-test-project'
  credential 'mycred'
end

```

#### Reference

```ruby
gsql_user 'id-for-resource' do
  host       string
  instance   reference to gsql_instance
  name       string
  password   string
  project    string
  credential reference to gauth_credential
end
```

#### Actions

* `create` -
  Converges the `gsql_user` resource into the final
  state described within the block. If the resource does not exist, Chef will
  attempt to create it.
* `delete` -
  Ensures the `gsql_user` resource is not present.
  If the resource already exists Chef will attempt to delete it.

#### Properties

* `host` -
  Required. The host name from which the user can connect. For insert operations,
  host defaults to an empty string. For update operations, host is
  specified as part of the request URL. The host name cannot be updated
  after insertion.

* `name` -
  Required. The name of the user in the Cloud SQL instance.

* `instance` -
  Required. The name of the Cloud SQL instance. This does not include the project
  ID.

* `password` -
  The password for the user.

#### Label
Set the `u_label` property when attempting to set primary key
of this object. The primary key will always be referred to by the initials of
the resource followed by "_label"

[google-gauth]: https://supermarket.chef.io/cookbooks/google-gauth
