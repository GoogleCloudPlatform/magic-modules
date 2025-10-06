package datastream_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDatastreamConnectionProfile_update(t *testing.T) {
	// this test uses the random provider
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	random_pass_1 := acctest.RandString(t, 10)
	random_pass_2 := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckDatastreamConnectionProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatastreamConnectionProfile_update(context),
			},
			{
				ResourceName:            "google_datastream_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_without_validation", "connection_profile_id", "location"},
			},
			{
				Config: testAccDatastreamConnectionProfile_update2(context, true),
			},
			{
				ResourceName:            "google_datastream_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_without_validation", "connection_profile_id", "location", "postgresql_profile.0.password"},
			},
			{
				// Disable prevent_destroy
				Config: testAccDatastreamConnectionProfile_update2(context, false),
			},
			{
				Config: testAccDatastreamConnectionProfile_mySQLUpdate(context, true, random_pass_1),
			},
			{
				ResourceName:            "google_datastream_connection_profile.mysql_con_profile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_without_validation", "connection_profile_id", "location", "mysql_profile.0.password"},
			},
			{
				// run once more to update the password. it should update it in-place
				Config: testAccDatastreamConnectionProfile_mySQLUpdate(context, true, random_pass_2),
			},
			{
				ResourceName:            "google_datastream_connection_profile.mysql_con_profile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_without_validation", "connection_profile_id", "location", "mysql_profile.0.password"},
			},
			{
				// Disable prevent_destroy
				Config: testAccDatastreamConnectionProfile_mySQLUpdate(context, false, random_pass_2),
			},
		},
	})
}

func TestAccDatastreamConnectionProfile_sshKey_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	randomPubKey1 := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCjXhptfWIrtflLZ1WeOsjCfHSEKvui0fdNXTqpqIA+2NNlFjwKS4mV3bDJIRlC5FdWG/D5LW4kvSmcTx1eSLUcvqw3i3F73Ii35AR1Rid1bY0LCBYUUgkDKyvZgDzrM7g+MwBtthoud8Axt9/bh28qtzSVNvWfxIYsa2CwtqlkZr5c6Qb6N2B9kxW8WFsCnoAeBaZDMq+LVBRsRJvBBrJm/qhMNPd07Al7wGLEnNPWmwjFT7B12sMjNr7ZNLfI9VckEyUSx3AGBFH7RImeYiWb6vZA9v5DE7kBrCoHtJK5IN9dvqEWXrrDT7RTFXd55xQqT70eZiIDNz1nexDw8ZCn user`
	randomPrivKey1 := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAo14abX1iK7X5S2dVnjrIwnx0hCr7otH3TV06qaiAPtjTZRY8CkuJ
ld2wySEZQuRXVhvw+S1uJL0pnE8dXki1HL6sN4txe9yIt+QEdUYndW2NCwgWFFIJAysr2Y
A86zO4PjMAbbYaLnfAMbff24dvKrc0lTb1n8SGLGtgsLapZGa+XOkG+jdgfZMVvFhbAp6A
HgWmQzKvi1QUbESbwQayZv6oTDT3dOwJe8BixJzT1psIxU+wddrDIza+2TS3yPVXJBMlEs
dwBgRR+0SJnmIlm+r2QPb+QxO5AawqB7SSuSDfXb6hFl66w0+0UxV3eecUKk+9HmYiAzc9
Z3sQ8PGQpwAAA8B2IBoLdiAaCwAAAAdzc2gtcnNhAAABAQCjXhptfWIrtflLZ1WeOsjCfH
SEKvui0fdNXTqpqIA+2NNlFjwKS4mV3bDJIRlC5FdWG/D5LW4kvSmcTx1eSLUcvqw3i3F7
3Ii35AR1Rid1bY0LCBYUUgkDKyvZgDzrM7g+MwBtthoud8Axt9/bh28qtzSVNvWfxIYsa2
CwtqlkZr5c6Qb6N2B9kxW8WFsCnoAeBaZDMq+LVBRsRJvBBrJm/qhMNPd07Al7wGLEnNPW
mwjFT7B12sMjNr7ZNLfI9VckEyUSx3AGBFH7RImeYiWb6vZA9v5DE7kBrCoHtJK5IN9dvq
EWXrrDT7RTFXd55xQqT70eZiIDNz1nexDw8ZCnAAAAAwEAAQAAAQAnvU5kb+mfhGaeBwb2
tIn9dVTKicIoezbTJOiOOKTppMjXgC8euf0/7WuBoYGJmg38rlNR6dEvMqyaj0wvkTQtR9
yQrmTuoljHkrna5TPYBswWcOMeEk6K7Md/4wfulugsiS+DgJah0xN3hKj5t9o848/wtCvP
r3iL+ZrNocFW4Ju+QrArFWTLFuJL4uc69ykgWE7I5Qkm+3Lg6aSoNazMzCu9rCblduetJq
EilQ6AOkv68xTOQ1EDIQc8xr6u6GCUvVVBwYaR3cYV6fWeLWJATqUODkEXdDZfgUerf4Io
3KirdRf0YFyJiHJh4AqWd76jWCkhCwrREx0lfMCZghoxAAAAgHwOfMJtd4wOug2BPKu0SA
HSwQ+yTTibg2xuENstd8akJC3VsU5GC8pngNAyoFpSt3QDlLpvqPqXVJSkkMbUtnPO0SIR
5ffMB97kFvNkMNDUIalwxR9DV1CMPTAnTO7NSfO8UUKRjKivpmpS6ptMjxUM0hPoDBebhx
P37In1a2jDAAAAgQDVCaoMFjHRGds1JaVjm6YviR0C2OsE55GOS7cW+I3SE63DumfHsN8i
r/u5oEQUelaauYVmi9tT3L4lReFX2tYqtyE0mbPUXcY5XfmBxBsjW1sQ6YyHlN/vGLgo33
NZZFpIg2FknTzM4qeddfbyKuqAJX27f7RrSZCf+WrJUKDWqwAAAIEAxFAn6d9na7uHnb31
TQ8PoTvkH7fwugXuG7ACLCTl3PpOSGPQAPI8rCaGOMd+uU1Jyjt3TcdPYlNAtiFQCxWLMH
RNFfeqviC85H6WzQNezNj45QqKTf5gRdHVu2NMRwn2pJjRgdIvsUaL1AY4sC0AivoEMlpx
rQYvdaDG7KsYXfUAAAAEdXNlcgECAwQFBgc=
-----END OPENSSH PRIVATE KEY-----`

	randomPubKey2 := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDmc1i/FqnVtYsTzb6LmoUGom8ISnfRCPTIFf3LLIyRFgO+qD6Dnqn5p2lLE8ksdooAGJ+EyJtV5c+3kYGnjzzH4TlB2pkt562BntrggvJ98sELQbHEDiemiLnJqqIESk5FcSXdcJ/UX/AdkbXLjSR5M8+cGGqKSb0HSnKfOWkjWwZwp/JwbvyWPIJ6IQNKzAS5HVU/J+u8ezhPd1iBdezvAuPlihpjMGQg1KW3APZoELS6/BSMpXcvDy+TwuggEPPZ0Up09BJRtqesHiZur6CnqUIzJcCWCfi5C8IfHzlhawry+iA1V5Lh06Mz7OaySXpf902RITfh+KcLxcSSMmPl user`
	randomPrivKey2 := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEA5nNYvxap1bWLE82+i5qFBqJvCEp30Qj0yBX9yyyMkRYDvqg+g56p
+adpSxPJLHaKABifhMibVeXPt5GBp488x+E5QdqZLeetgZ7a4ILyffLBC0GxxA4npoi5ya
qiBEpORXEl3XCf1F/wHZG1y40keTPPnBhqikm9B0pynzlpI1sGcKfycG78ljyCeiEDSswE
uR1VPyfrvHs4T3dYgXXs7wLj5YoaYzBkINSltwD2aBC0uvwUjKV3Lw8vk8LoIBDz2dFKdP
QSUbanrB4mbq+gp6lCMyXAlgn4uQvCHx85YWsK8vogNVeS4dOjM+zmskl6X/dNkSE34fin
C8XEkjJj5QAAA8CppfYQqaX2EAAAAAdzc2gtcnNhAAABAQDmc1i/FqnVtYsTzb6LmoUGom
8ISnfRCPTIFf3LLIyRFgO+qD6Dnqn5p2lLE8ksdooAGJ+EyJtV5c+3kYGnjzzH4TlB2pkt
562BntrggvJ98sELQbHEDiemiLnJqqIESk5FcSXdcJ/UX/AdkbXLjSR5M8+cGGqKSb0HSn
KfOWkjWwZwp/JwbvyWPIJ6IQNKzAS5HVU/J+u8ezhPd1iBdezvAuPlihpjMGQg1KW3APZo
ELS6/BSMpXcvDy+TwuggEPPZ0Up09BJRtqesHiZur6CnqUIzJcCWCfi5C8IfHzlhawry+i
A1V5Lh06Mz7OaySXpf902RITfh+KcLxcSSMmPlAAAAAwEAAQAAAQEAq2opHRpSgfBj3vsv
PNBXGrRAOr6JmSc8TIhvG22rsU/awTqMJYMjk9v+6iVxgm06ARBPt4kwYhhrBXRqKKTW5S
aWXHGpdwfZe40Z6d39Wcnz5debzuVogOs6ptMRaHeM+QJM1AYuHN6v0I7N1vbJpo3vY4CV
3v8yZ/XshJtDpVNqHFuCh1r07aW4NlqoTy5TEvWD1VPCqAVwTLWuNMfWRGYbwqJrRUxuu3
6vqddE8yMONYMwVRKPADj0DTi3i+LK3v6QfJlxb09EhqJPOOXM+fBVzUWkUXlPjvMP4uUH
/zRrGscSI93n0V/H3/XTOJTskdEZUEFpeFbUXIphloCKEQAAAIEA9CJapVXG9HcKimXX3I
OQdwPoKONM52KnAoWjGO1N5ECydjz2yHQkNJNLFwAUefmKVy0/ce0EdyEJjoHKvCwoTWL6
3CPlWQY+7pk0Fr62iT7UjjGwCtmHB6B5G4qUlsBkVN3WCwfmBwYrziRR+qcS8hSS7m37Uy
rMbGGIHHVGPzIAAACBAP6ouUUlIN7jLdLxyApj1Cx7oW7Gp33j3goXn82WVv6+ubPJymVD
u7zmoWWVegOngoPlR1q/mHBGoB1Ec1Im5IaN5qzVrxVKraJz5Q1XRc/azpkYb1FaDFBW2O
iDaP5PHvNQpYcmE82Dg8bUqa7tYIUgq2vqHJdBZC5IvnYnGrWbAAAAgQDnqf2DVITbK5jK
UJqEmni0YE8PD3PuPGRWLmZeOcxshHR1nQIeUoXWAhCS9G7Rl5Kdr1IXzSln22OvUXMPmE
gZLd7QJVyRQ0bXhYf8nIs/UGhjq83OSoS4iSwHeZ1CrKWmVP74/+Na6fDdfJ65Z8+I4ktM
QC3v6moZVb2wrgGkfwAAAAR1c2VyAQIDBAU=
-----END OPENSSH PRIVATE KEY-----`

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "registry.terraform.io/hashicorp/time",
			},
		},
		CheckDestroy: testAccCheckDatastreamConnectionProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatastreamConnectionProfile_sshKey_update(context, true, randomPrivKey1, randomPubKey1),
			},
			{
				ResourceName:            "google_datastream_connection_profile.ssh_connectivity_profile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "create_without_validation", "forward_ssh_connectivity.0.private_key", "postgresql_profile.0.password"},
			},
			{
				PreConfig: func() {
					fmt.Println("Waiting before proceeding to the next step...")
					time.Sleep(150 * time.Second) // Delay before the next step
				},
				Config: testAccDatastreamConnectionProfile_sshKey_update(context, true, randomPrivKey2, randomPubKey2),
			},
			{
				ResourceName:            "google_datastream_connection_profile.ssh_connectivity_profile",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "location", "create_without_validation", "forward_ssh_connectivity.0.private_key", "postgresql_profile.0.password"},
			},
			{
				PreConfig: func() {
					fmt.Println("Waiting before proceeding to the next step...")
					time.Sleep(150 * time.Second) // Delay before the next step
				},
				Config: testAccDatastreamConnectionProfile_sshKey_update(context, false, randomPrivKey2, randomPubKey2),
			},
		},
	})
}

func testAccDatastreamConnectionProfile_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_datastream_connection_profile" "default" {
	display_name          = "Connection profile"
	location              = "us-central1"
	connection_profile_id = "tf-test-my-profile%{random_suffix}"

	gcs_profile {
		bucket    = "my-bucket"
		root_path = "/path"
	}
	lifecycle {
		prevent_destroy = true
	}
}
`, context)
}

func testAccDatastreamConnectionProfile_update2(context map[string]interface{}, preventDestroy bool) string {
	context["lifecycle_block"] = ""
	if preventDestroy {
		context["lifecycle_block"] = `
		lifecycle {
			prevent_destroy = true
		}`
	}
	return acctest.Nprintf(`
resource "google_sql_database_instance" "instance" {
    name             = "tf-test-my-database-instance%{random_suffix}"
    database_version = "POSTGRES_14"
    region           = "us-central1"
    settings {
      tier = "db-f1-micro"

      ip_configuration {

        // Datastream IPs will vary by region.
        authorized_networks {
            value = "34.71.242.81"
        }

        authorized_networks {
            value = "34.72.28.29"
        }

        authorized_networks {
            value = "34.67.6.157"
        }

        authorized_networks {
            value = "34.67.234.134"
        }

        authorized_networks {
            value = "34.72.239.218"
        }
      }
    }

    deletion_protection  = "false"
}

resource "google_sql_database" "db" {
    instance = google_sql_database_instance.instance.name
    name     = "db"
}

resource "random_password" "pwd" {
    length = 16
    special = false
}

resource "google_sql_user" "user" {
    name = "user"
    instance = google_sql_database_instance.instance.name
    password = random_password.pwd.result
}

resource "google_datastream_connection_profile" "default" {
	display_name          = "Connection profile"
	location              = "us-central1"
	connection_profile_id = "tf-test-my-profile%{random_suffix}"

	postgresql_profile {
		hostname = google_sql_database_instance.instance.public_ip_address
		username = google_sql_user.user.name
		password = google_sql_user.user.password
		database = google_sql_database.db.name
	}
	%{lifecycle_block}
}
`, context)
}

func testAccDatastreamConnectionProfile_mySQLUpdate(context map[string]interface{}, preventDestroy bool, password string) string {
	context["lifecycle_block"] = ""
	if preventDestroy {
		context["lifecycle_block"] = `
		lifecycle {
			prevent_destroy = true
		}`
	}

	context["password"] = password

	return acctest.Nprintf(`
resource "google_sql_database_instance" "mysql_instance" {
    name             = "tf-test-mysql-database-instance%{random_suffix}"
    database_version = "MYSQL_8_0"
    region           = "us-central1"
    settings {
      tier = "db-f1-micro"
        backup_configuration {
            enabled            = true
            binary_log_enabled = true
        }

      ip_configuration {

        // Datastream IPs will vary by region.
        authorized_networks {
            value = "34.71.242.81"
        }

        authorized_networks {
            value = "34.72.28.29"
        }

        authorized_networks {
            value = "34.67.6.157"
        }

        authorized_networks {
            value = "34.67.234.134"
        }

        authorized_networks {
            value = "34.72.239.218"
        }
      }
    }

    deletion_protection  = "false"
}

resource "google_sql_database" "mysql_db" {
    instance = google_sql_database_instance.mysql_instance.name
    name     = "db"
}

resource "google_sql_user" "mysql_user" {
    name = "user"
    instance = google_sql_database_instance.mysql_instance.name
    host     = "%"
    password = "%{password}"
}

resource "google_datastream_connection_profile" "mysql_con_profile" {
    display_name          = "Source connection profile"
	location              = "us-central1"
	connection_profile_id = "tf-test-mysql-profile%{random_suffix}"

    mysql_profile {
		hostname = google_sql_database_instance.mysql_instance.public_ip_address
		username = google_sql_user.mysql_user.name
		password = google_sql_user.mysql_user.password
	}
	%{lifecycle_block}
}
`, context)
}

func testAccDatastreamConnectionProfile_sshKey_update(context map[string]interface{}, preventDestroy bool, private_key string, public_key string) string {
	context["lifecycle_block"] = ""
	if preventDestroy {
		context["lifecycle_block"] = `
        lifecycle {
            prevent_destroy = true
        }`
	}
	context["private_key"] = private_key
	context["public_key"] = public_key

	return acctest.Nprintf(`
resource "google_compute_network" "default" {
		name = "tf-test-datastream-ssh%{random_suffix}"
}

resource "google_sql_database_instance" "instance" {
    depends_on         = [google_compute_instance.default]
    name            	= "tf-test-my-database-instance%{random_suffix}"
    database_version	= "POSTGRES_14"
    region           	= "us-central1"
    settings {
        tier = "db-f1-micro"
        ip_configuration {
			ipv4_enabled = true

			authorized_networks {
				value = google_compute_instance.default.network_interface.0.access_config.0.nat_ip
			}
        }
    }
    
    deletion_protection  = "false"
}

resource "google_sql_database" "db" {
	depends_on = [google_sql_database_instance.instance]
	instance = google_sql_database_instance.instance.name
	name     = "db"
}

resource "google_sql_user" "user" {
	depends_on	= [google_sql_database_instance.instance]
	name		= "user"
	instance	= google_sql_database_instance.instance.name
	password	= "password%{random_suffix}"
}

resource "google_compute_instance" "default" {
	name         = "tf-test-instance-%{random_suffix}"
	machine_type = "e2-small"
	zone         = "us-central1-a"
	boot_disk {
		initialize_params {
			image = "debian-11-bullseye-v20241009"
		}
	}

	network_interface {
		network    = google_compute_network.default.name
		access_config {}
		}

	metadata = {
		"ssh-keys" = "user:%{public_key}"
	}

	metadata_startup_script = <<-EOT
	#!/bin/bash
	echo "Updating SSHD config for SSH forwarding..."

	# Backup sshd_config
	echo "AllowTcpForwarding yes" >> /etc/ssh/sshd_config
	echo "PasswordAuthentication no" >> /etc/ssh/sshd_config
	echo "PubkeyAuthentication yes" >> /etc/ssh/sshd_config
	echo "AuthorizedKeysFile .ssh/authorized_keys" >> /etc/ssh/sshd_config
	
	# Restart SSH service
	systemctl restart sshd
	EOT

	tags = ["ssh-host"]

	depends_on = [google_compute_firewall.ssh, google_compute_firewall.datastream_sql_access]

}

resource "time_sleep" "ssh_host_wait" {
	depends_on = [google_compute_instance.default]
	create_duration = "12m"
}

resource "google_compute_firewall" "ssh" {
	name 	= "tf-test-%{random_suffix}"
	network =  google_compute_network.default.name

	allow {
		protocol = "tcp"
		ports    = ["22"]
	}

	direction     = "INGRESS"
	priority      = 1000
	source_ranges = ["34.71.242.81", "34.72.28.29", "34.67.6.157", "34.67.234.134", "34.72.239.218"]

	target_tags = ["ssh-host"]
}

resource "google_compute_firewall" "datastream_sql_access" {
    name    	= "datastream-to-cloudsql-%{random_suffix}"
    network 	=  google_compute_network.default.name

    allow {
        protocol = "tcp"
        ports    = ["5432"]
    }

    direction     = "INGRESS"
    priority      = 1000
    source_ranges = ["34.71.242.81", "34.72.28.29", "34.67.6.157", "34.67.234.134", "34.72.239.218"]  #Datastream IPs

}

resource "google_datastream_connection_profile" "ssh_connectivity_profile" {
    display_name          = "Source connection profile"
    location              = "us-central1"
    connection_profile_id = "tf-test-pg-profile%{random_suffix}"

    postgresql_profile {
        hostname 			= google_sql_database_instance.instance.public_ip_address
        username 			= google_sql_user.user.name
        password 			= google_sql_user.user.password
        database 			= google_sql_database.db.name
        port 				= 5432
    }

    forward_ssh_connectivity {
        hostname 	= google_compute_instance.default.network_interface.0.access_config.0.nat_ip
        username 	= google_sql_user.user.name
        port    	= 22
        private_key 	= <<EOT
%{private_key}
EOT
	}

	depends_on = [time_sleep.ssh_host_wait]
	timeouts {
         create = "10m"
	}
    %{lifecycle_block}
}
`, context)
}

func TestAccDatastreamConnectionProfile_datastreamConnectionProfileMongodbSrv(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"deletion_protection": false,
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDatastreamConnectionProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatastreamConnectionProfile_datastreamConnectionProfileMongodbSrv(context),
			},
			{
				ResourceName:            "google_datastream_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "create_without_validation", "labels", "location", "terraform_labels", "mongodb_profile.0.ssl_config.0.client_key", "mongodb_profile.0.ssl_config.0.client_certificate", "mongodb_profile.0.ssl_config.0.ca_certificate", "mongodb_profile.0.password", "mongodb_profile.0.srv_connection_format"},
			},
		},
	})
}

func testAccDatastreamConnectionProfile_datastreamConnectionProfileMongodbSrv(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_datastream_connection_profile" "default" {
    display_name              = "Connection profile MongoDB SRV"
    location                  = "us-central1"
    connection_profile_id     = "tf-test-source-profile%{random_suffix}"
    create_without_validation = true

    mongodb_profile {
        srv_connection_format {}

        host_addresses {
          hostname = "example.mongodb.com"
        }
        
        username       = "username%{random_suffix}"
        password       = "password%{random_suffix}"

		    ssl_config {
		        ca_certificate     = "-----BEGIN CERTIFICATE-----\nMIIDfzCCAmegAwIBAgIJAMz9tB11EHe1MA0GCSqGSIb3DQEBCwUAMG4xCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApUZXN0IFN0YXRlMRYwFAYDVQQHDA1UZXN0IExvY2F0\naW9uMRowGAYDVQQKDBFUZXN0IE9yZ2FuaXphdGlvbjEWMBQGA1UEAwwNRHVtbXkg\nVGVzdCBDQTAeFw0yNTEwMDYxMzEzMjRaFw0yNjEwMDYxMzEzMjRaMG4xCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApUZXN0IFN0YXRlMRYwFAYDVQQHDA1UZXN0IExvY2F0\naW9uMRowGAYDVQQKDBFUZXN0IE9yZ2FuaXphdGlvbjEWMBQGA1UEAwwNRHVtbXkg\nVGVzdCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANllxriVNxHY\nGPs3V/rk/oTePcxY4ARYNRVIpw/G7t6828yUGkNRowTNcF0Lf/28VroqTyuC4oPC\nognoEslMdHUC1X5+8EDChz6yO/N3gefL6BFkBYKCaU/1MTSpWuxzLnOVpcMYNBeS\nuDtt5TroeKdZqSh5pTtNriQ5QYorujA6wFeDODqN0eluJwBAL/I8JAkEMRKdv+md\nhjNMCRyQJMUAVTV5cKYa3hwdyaRFORid00RUzPvxJdj40WbJPdT+OhGXCL0yrlWa\njEN3odNvTWhhMzfJZfh6Q3f1sDs/ubNFKGRrcQukb2sK7GQO2ltxTLR6CEh9tT5h\nbhfx022SMukCAwEAAaMgMB4wDwYDVR0TAQH/BAUwAwEB/zALBgNVHQ8EBAMCAQYw\nDQYJKoZIhvcNAQELBQADggEBAHt7xETI2AZEUZfS2WxfFLNmh8WMFxD597kdDdsj\nX4ZXLiUfkOFIkcWwdKYrYibG09Ps4rR4BAtV+2JNwHut059lOBOPR6gBJ44sjpBf\nHHHjrqJ1a6D+wwRUKuK5qSGlnB+l8qy1OZjcRxq9dljw+zFooRUCbZps8mk+lc+a\nwlVXHXZgVM9y4RDxb2CeWwt8al0gakf/vH3XBagxXj0oYS86eVGzS7rpAxRDPROy\nBNzzNFCkymVAQEO7XvLcOf6nD/jFYvfYwHCGfVMpUdAxG2oSzi+Oa1U40NqJKrUg\ndnnyl8jlciOkAslweooS0KUfvAVGSOxC1dmtKtHsguPsbC8=\n-----END CERTIFICATE-----"
	          client_certificate = "-----BEGIN CERTIFICATE-----\nMIIDXDCCAkQCCQCM9siYSDyyCDANBgkqhkiG9w0BAQsFADBuMQswCQYDVQQGEwJV\nUzETMBEGA1UECAwKVGVzdCBTdGF0ZTEWMBQGA1UEBwwNVGVzdCBMb2NhdGlvbjEa\nMBgGA1UECgwRVGVzdCBPcmdhbml6YXRpb24xFjAUBgNVBAMMDUR1bW15IFRlc3Qg\nQ0EwHhcNMjUxMDA2MTMxMzM0WhcNMjYxMDA2MTMxMzM0WjByMQswCQYDVQQGEwJV\nUzETMBEGA1UECAwKVGVzdCBTdGF0ZTEWMBQGA1UEBwwNVGVzdCBMb2NhdGlvbjEa\nMBgGA1UECgwRVGVzdCBPcmdhbml6YXRpb24xGjAYBgNVBAMMEUR1bW15IFRlc3Qg\nQ2xpZW50MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsuJVeZAuAaPs\nTZwtqCfWzxK3QSehXMQz8p3zjQxPAMgZW5N9VozWiQ5i6E1mVZhW4ugPlD6z+znu\n9CCfuXiAJ5jm8+piOZUON4bRWtuUikQSe6OD4kdgt668liT8ozuN/jrfpxKKnpR0\nFMITfGYmOCU4dRR2K8FTNRBQZXXqBqck9XNXWibr/1b78U+IVXZHq66Ofg59o77e\nYTyzSuSNWZ73pYvhFfcW5tnA6j3ERULWWvYn5SglHX+NLJXHxcort8IUToxZMVn/\n14QApWY5HKSznWI/KgClQPrsWHIf1S2QbUwip41A+SAfP1ZZBqRVMhVp6cGkLbNh\nZ9bOT4OvIQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAh2AUtc5q30ACsM8KODDwY\nzEXPDINBGuEcBdTyEb/RWBGyt/QZzbnDhA72IfJ8WX9Etlo6WUJgk4JF7hBJ+wra\nY4EirsOt9l4zk8fjiLADgYMB+sySQl3sy4cVY9dI1UhQPpVoyjZ9JV7qZVwCAXqV\nZug7ulwQUULJDXQnmC7ZCQDJSolJ2wg0cC1FyLakZdZfRS9dhtnrzcbn5th1e/0n\nngQHunR8C94LuCFr+qJZa1+D813T5+OwFUThzwk7B+Ke35wZHALtgvnM0Jc6OR5+\nXVAk0fYAh/fKvIIOoKgF+qRNGcUP83rXQeg69J5mAezb4mw638LuOmKia+EOhoAx\n-----END CERTIFICATE-----"
	          client_key         = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpgIBAAKCAQEAsuJVeZAuAaPsTZwtqCfWzxK3QSehXMQz8p3zjQxPAMgZW5N9\nVozWiQ5i6E1mVZhW4ugPlD6z+znu9CCfuXiAJ5jm8+piOZUON4bRWtuUikQSe6OD\n4kdgt668liT8ozuN/jrfpxKKnpR0FMITfGYmOCU4dRR2K8FTNRBQZXXqBqck9XNX\nWibr/1b78U+IVXZHq66Ofg59o77eYTyzSuSNWZ73pYvhFfcW5tnA6j3ERULWWvYn\n5SglHX+NLJXHxcort8IUToxZMVn/14QApWY5HKSznWI/KgClQPrsWHIf1S2QbUwi\np41A+SAfP1ZZBqRVMhVp6cGkLbNhZ9bOT4OvIQIDAQABAoIBAQCH6oatVcpO/sD1\n+xuJr7N8JKlOfRESzhT2W+MIoXiJjIAP35GVKG99NYwbG2wMzzH9N/tWVQolcVBI\n91zE7HTbIUchv02gmMtzjyEU2tAS+kPc41G6pScsiTzLDBFU6VQq/Yqfg+wFL6C/\ngPKTS33wnP83njNnbX2OTPX5EU2efS2XiJYJ1MeKszP/gLzNT/nzPVPDQbwy2BBv\nlcUQm5h3kBYUxpFd6JIEohmxgFBfNY4r8UoC00/5xIJmUVJ30cwxJnGd9jiLT0h4\nCW+j/uOS51c7Z60IfoXV7ggkHvwdlx0Plh5SPfqXjmUJxHXKNYwipMxyso6Qt2AI\nO3tLD/8RAoGBAN7NcIoO5+IF9t4h7jZIxJm8nhkIIXnGAOHjEI3w9j/zFyCCYDxx\n6GBPGXY/icNx/JTRt3749s94nX+aIvsniNsKaVwpOWin5bqq12As8DDUiK0CAKVc\n1fxjAa0w6V6kedU33pNXAXn1q2azgtFc5OXk9FWoVKYqbDw6agXHH6ZFAoGBAM2J\nrNMhzhEmCt5e5Qx9k+6KaSxbTYLD/TcJP3KDBSYfMgBqcr0FxOXlETSWgkomWXit\nQCFklwf/egMcYT5dIwoN7cFXUmG6/+h+oPVPwaT3puvvtKhlZXDndFbUDC4R97AX\nEFeelgZAS6fx38c79VUCQpFQnMrqCmjpxVXww3EtAoGBAKrwAXzajMubedjZPXMG\nh1fQH5fi5hQQdtLXq/bKvZM4xTCa9ozJc9iYN1fCzcZWqMvgzqCrEGkDCAtDTb1V\niqlLJqSfuDz0O8vokQ9nyuwb07SwyaAVRtO5firLUPDczeBpWem/IhHZCyTjauWI\nGNHMxC0H1dIa0CmxQ3ClYkHlAoGBALjQT890+SbgTyuOpmRp0ofOey2AV5z6gAhp\nz1w3RXz21e4byVoAAwE4zRS9NSBZhWAGYMDmAwwVA3Aip6n881HKHnwX+aKZFBzJ\nKBAMjDG64aQK4SX+Lo2sASdF+kG+tDnpMy+mEH5EeALmcXJjjoDGzHZ/xsyKT5vw\ngBl7qTFtAoGBAI9POAsCVU9rwK3BqZB0iWyvynBarR6QmSkZf6TAOPcKc7M+QJO5\nKWJ4R1a58g1gpVdcKgU15Nym6WLSCr86+XYNPUWr+TbGxdAsRb8PJwDJVrFvi1W+\ni7dqZQELjScVyttHDE82OmQvt8OocEyhLB/zFXnwc+nNycGYr5H21dOp\n-----END RSA PRIVATE KEY-----"
		    }
    }
}
`, context)
}

func TestAccDatastreamConnectionProfile_datastreamConnectionProfileMongodbStrdFull(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"deletion_protection": false,
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDatastreamConnectionProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDatastreamConnectionProfile_datastreamConnectionProfileMongodbStrdFull(context),
			},
			{
				ResourceName:            "google_datastream_connection_profile.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_profile_id", "create_without_validation", "labels", "location", "terraform_labels", "mongodb_profile.0.ssl_config.0.client_key", "mongodb_profile.0.ssl_config.0.client_certificate", "mongodb_profile.0.ssl_config.0.ca_certificate", "mongodb_profile.0.password", "mongodb_profile.0.srv_connection_format"},
			},
		},
	})
}

func testAccDatastreamConnectionProfile_datastreamConnectionProfileMongodbStrdFull(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "password" {
  secret_id = "password-%{random_suffix}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "password" {
  secret = google_secret_manager_secret.password.id

  secret_data = "password%{random_suffix}"
}

resource "google_secret_manager_secret" "client_key" {
  secret_id = "client-key-%{random_suffix}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "client_key" {
  secret = google_secret_manager_secret.client_key.id

  secret_data = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpgIBAAKCAQEAsuJVeZAuAaPsTZwtqCfWzxK3QSehXMQz8p3zjQxPAMgZW5N9\nVozWiQ5i6E1mVZhW4ugPlD6z+znu9CCfuXiAJ5jm8+piOZUON4bRWtuUikQSe6OD\n4kdgt668liT8ozuN/jrfpxKKnpR0FMITfGYmOCU4dRR2K8FTNRBQZXXqBqck9XNX\nWibr/1b78U+IVXZHq66Ofg59o77eYTyzSuSNWZ73pYvhFfcW5tnA6j3ERULWWvYn\n5SglHX+NLJXHxcort8IUToxZMVn/14QApWY5HKSznWI/KgClQPrsWHIf1S2QbUwi\np41A+SAfP1ZZBqRVMhVp6cGkLbNhZ9bOT4OvIQIDAQABAoIBAQCH6oatVcpO/sD1\n+xuJr7N8JKlOfRESzhT2W+MIoXiJjIAP35GVKG99NYwbG2wMzzH9N/tWVQolcVBI\n91zE7HTbIUchv02gmMtzjyEU2tAS+kPc41G6pScsiTzLDBFU6VQq/Yqfg+wFL6C/\ngPKTS33wnP83njNnbX2OTPX5EU2efS2XiJYJ1MeKszP/gLzNT/nzPVPDQbwy2BBv\nlcUQm5h3kBYUxpFd6JIEohmxgFBfNY4r8UoC00/5xIJmUVJ30cwxJnGd9jiLT0h4\nCW+j/uOS51c7Z60IfoXV7ggkHvwdlx0Plh5SPfqXjmUJxHXKNYwipMxyso6Qt2AI\nO3tLD/8RAoGBAN7NcIoO5+IF9t4h7jZIxJm8nhkIIXnGAOHjEI3w9j/zFyCCYDxx\n6GBPGXY/icNx/JTRt3749s94nX+aIvsniNsKaVwpOWin5bqq12As8DDUiK0CAKVc\n1fxjAa0w6V6kedU33pNXAXn1q2azgtFc5OXk9FWoVKYqbDw6agXHH6ZFAoGBAM2J\nrNMhzhEmCt5e5Qx9k+6KaSxbTYLD/TcJP3KDBSYfMgBqcr0FxOXlETSWgkomWXit\nQCFklwf/egMcYT5dIwoN7cFXUmG6/+h+oPVPwaT3puvvtKhlZXDndFbUDC4R97AX\nEFeelgZAS6fx38c79VUCQpFQnMrqCmjpxVXww3EtAoGBAKrwAXzajMubedjZPXMG\nh1fQH5fi5hQQdtLXq/bKvZM4xTCa9ozJc9iYN1fCzcZWqMvgzqCrEGkDCAtDTb1V\niqlLJqSfuDz0O8vokQ9nyuwb07SwyaAVRtO5firLUPDczeBpWem/IhHZCyTjauWI\nGNHMxC0H1dIa0CmxQ3ClYkHlAoGBALjQT890+SbgTyuOpmRp0ofOey2AV5z6gAhp\nz1w3RXz21e4byVoAAwE4zRS9NSBZhWAGYMDmAwwVA3Aip6n881HKHnwX+aKZFBzJ\nKBAMjDG64aQK4SX+Lo2sASdF+kG+tDnpMy+mEH5EeALmcXJjjoDGzHZ/xsyKT5vw\ngBl7qTFtAoGBAI9POAsCVU9rwK3BqZB0iWyvynBarR6QmSkZf6TAOPcKc7M+QJO5\nKWJ4R1a58g1gpVdcKgU15Nym6WLSCr86+XYNPUWr+TbGxdAsRb8PJwDJVrFvi1W+\ni7dqZQELjScVyttHDE82OmQvt8OocEyhLB/zFXnwc+nNycGYr5H21dOp\n-----END RSA PRIVATE KEY-----"
}

resource "google_datastream_connection_profile" "default" {
    display_name              = "Connection profile MongoDB SRV"
    location                  = "us-central1"
    connection_profile_id     = "tf-test-source-profile%{random_suffix}"
    create_without_validation = true

    mongodb_profile {
			  standard_connection_format {
				   direct_connection = true
				}
		    replica_set = "replica_set_name"
        host_addresses {
					  hostname = "example.mongodb1.com"
					  port = 27017
        }
        host_addresses {
            hostname = "example.mongodb2.com"
					  port = 27017
        }
        host_addresses {
            hostname = "example.mongodb3.com"
					  port = 27017
        }

        username                       = "username%{random_suffix}"
		    secret_manager_stored_password = google_secret_manager_secret_version.password.name
		    
		    ssl_config {
		        ca_certificate     = "-----BEGIN CERTIFICATE-----\nMIIDfzCCAmegAwIBAgIJAMz9tB11EHe1MA0GCSqGSIb3DQEBCwUAMG4xCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApUZXN0IFN0YXRlMRYwFAYDVQQHDA1UZXN0IExvY2F0\naW9uMRowGAYDVQQKDBFUZXN0IE9yZ2FuaXphdGlvbjEWMBQGA1UEAwwNRHVtbXkg\nVGVzdCBDQTAeFw0yNTEwMDYxMzEzMjRaFw0yNjEwMDYxMzEzMjRaMG4xCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApUZXN0IFN0YXRlMRYwFAYDVQQHDA1UZXN0IExvY2F0\naW9uMRowGAYDVQQKDBFUZXN0IE9yZ2FuaXphdGlvbjEWMBQGA1UEAwwNRHVtbXkg\nVGVzdCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANllxriVNxHY\nGPs3V/rk/oTePcxY4ARYNRVIpw/G7t6828yUGkNRowTNcF0Lf/28VroqTyuC4oPC\nognoEslMdHUC1X5+8EDChz6yO/N3gefL6BFkBYKCaU/1MTSpWuxzLnOVpcMYNBeS\nuDtt5TroeKdZqSh5pTtNriQ5QYorujA6wFeDODqN0eluJwBAL/I8JAkEMRKdv+md\nhjNMCRyQJMUAVTV5cKYa3hwdyaRFORid00RUzPvxJdj40WbJPdT+OhGXCL0yrlWa\njEN3odNvTWhhMzfJZfh6Q3f1sDs/ubNFKGRrcQukb2sK7GQO2ltxTLR6CEh9tT5h\nbhfx022SMukCAwEAAaMgMB4wDwYDVR0TAQH/BAUwAwEB/zALBgNVHQ8EBAMCAQYw\nDQYJKoZIhvcNAQELBQADggEBAHt7xETI2AZEUZfS2WxfFLNmh8WMFxD597kdDdsj\nX4ZXLiUfkOFIkcWwdKYrYibG09Ps4rR4BAtV+2JNwHut059lOBOPR6gBJ44sjpBf\nHHHjrqJ1a6D+wwRUKuK5qSGlnB+l8qy1OZjcRxq9dljw+zFooRUCbZps8mk+lc+a\nwlVXHXZgVM9y4RDxb2CeWwt8al0gakf/vH3XBagxXj0oYS86eVGzS7rpAxRDPROy\nBNzzNFCkymVAQEO7XvLcOf6nD/jFYvfYwHCGfVMpUdAxG2oSzi+Oa1U40NqJKrUg\ndnnyl8jlciOkAslweooS0KUfvAVGSOxC1dmtKtHsguPsbC8=\n-----END CERTIFICATE-----"
	          client_certificate = "-----BEGIN CERTIFICATE-----\nMIIDXDCCAkQCCQCM9siYSDyyCDANBgkqhkiG9w0BAQsFADBuMQswCQYDVQQGEwJV\nUzETMBEGA1UECAwKVGVzdCBTdGF0ZTEWMBQGA1UEBwwNVGVzdCBMb2NhdGlvbjEa\nMBgGA1UECgwRVGVzdCBPcmdhbml6YXRpb24xFjAUBgNVBAMMDUR1bW15IFRlc3Qg\nQ0EwHhcNMjUxMDA2MTMxMzM0WhcNMjYxMDA2MTMxMzM0WjByMQswCQYDVQQGEwJV\nUzETMBEGA1UECAwKVGVzdCBTdGF0ZTEWMBQGA1UEBwwNVGVzdCBMb2NhdGlvbjEa\nMBgGA1UECgwRVGVzdCBPcmdhbml6YXRpb24xGjAYBgNVBAMMEUR1bW15IFRlc3Qg\nQ2xpZW50MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsuJVeZAuAaPs\nTZwtqCfWzxK3QSehXMQz8p3zjQxPAMgZW5N9VozWiQ5i6E1mVZhW4ugPlD6z+znu\n9CCfuXiAJ5jm8+piOZUON4bRWtuUikQSe6OD4kdgt668liT8ozuN/jrfpxKKnpR0\nFMITfGYmOCU4dRR2K8FTNRBQZXXqBqck9XNXWibr/1b78U+IVXZHq66Ofg59o77e\nYTyzSuSNWZ73pYvhFfcW5tnA6j3ERULWWvYn5SglHX+NLJXHxcort8IUToxZMVn/\n14QApWY5HKSznWI/KgClQPrsWHIf1S2QbUwip41A+SAfP1ZZBqRVMhVp6cGkLbNh\nZ9bOT4OvIQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAh2AUtc5q30ACsM8KODDwY\nzEXPDINBGuEcBdTyEb/RWBGyt/QZzbnDhA72IfJ8WX9Etlo6WUJgk4JF7hBJ+wra\nY4EirsOt9l4zk8fjiLADgYMB+sySQl3sy4cVY9dI1UhQPpVoyjZ9JV7qZVwCAXqV\nZug7ulwQUULJDXQnmC7ZCQDJSolJ2wg0cC1FyLakZdZfRS9dhtnrzcbn5th1e/0n\nngQHunR8C94LuCFr+qJZa1+D813T5+OwFUThzwk7B+Ke35wZHALtgvnM0Jc6OR5+\nXVAk0fYAh/fKvIIOoKgF+qRNGcUP83rXQeg69J5mAezb4mw638LuOmKia+EOhoAx\n-----END CERTIFICATE-----"
		        secret_manager_stored_client_key = google_secret_manager_secret_version.client_key.name
		    }
    }
}
`, context)
}
