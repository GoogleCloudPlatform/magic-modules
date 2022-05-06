package google

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func diffSuppressIamUserName(_, old, new string, d *schema.ResourceData) bool {
	strippedName := strings.Split(new, "@")[0]

	userType := d.Get("type").(string)

	if old == strippedName && strings.Contains(userType, "IAM") {
		return true
	}

	return false
}

func resourceSqlUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlUserCreate,
		Read:   resourceSqlUserRead,
		Update: resourceSqlUserUpdate,
		Delete: resourceSqlUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSqlUserImporter,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		MigrateState:  resourceSqlUserMigrateState,

		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The host the user can connect from. This is only supported for MySQL instances. Don't set this field for PostgreSQL instances. Can be an IP address. Changing this forces a new resource to be created.`,
			},

			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the Cloud SQL instance. Changing this forces a new resource to be created.`,
			},

			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: diffSuppressIamUserName,
				Description:      `The name of the user. Changing this forces a new resource to be created.`,
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Description: `The password for the user. Can be updated. For Postgres instances this is a Required field, unless type is set to
                either CLOUD_IAM_USER or CLOUD_IAM_SERVICE_ACCOUNT.`,
			},

			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: emptyOrDefaultStringSuppress("BUILT_IN"),
				Description: `The user type. It determines the method to authenticate the user during login.
                The default is the database's built-in user type. Flags include "BUILT_IN", "CLOUD_IAM_USER", or "CLOUD_IAM_SERVICE_ACCOUNT".`,
				ValidateFunc: validation.StringInSlice([]string{"BUILT_IN", "CLOUD_IAM_USER", "CLOUD_IAM_SERVICE_ACCOUNT", ""}, false),
			},

			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `If the user has been disabled.`,
			},
			"server_roles": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The server roles for this user in the database.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"deletion_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The deletion policy for the user. Setting ABANDON allows the resource
				to be abandoned rather than deleted. This is useful for Postgres, where users cannot be deleted from the API if they
				have been granted SQL roles. Possible values are: "ABANDON".`,
				ValidateFunc: validation.StringInSlice([]string{"ABANDON", ""}, false),
			},
		},
		UseJSONNumber: true,
	}
}

func resourceSqlUserCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	instance := d.Get("instance").(string)
	password := d.Get("password").(string)
	host := d.Get("host").(string)
	typ := d.Get("type").(string)
	disabled := d.Get("disabled").(bool)
	serverRoles := expandStringArray(d.Get("server_roles"))

	user := &sqladmin.User{
		Name:     name,
		Instance: instance,
		Password: password,
		Host:     host,
		Type:     typ,
		SqlserverUserDetails: &sqladmin.SqlServerUserDetails{
			Disabled:    disabled,
			ServerRoles: serverRoles,
		},
	}

	mutexKV.Lock(instanceMutexKey(project, instance))
	defer mutexKV.Unlock(instanceMutexKey(project, instance))
	var op *sqladmin.Operation
	insertFunc := func() error {
		op, err = config.NewSqlAdminClient(userAgent).Users.Insert(project, instance,
			user).Do()
		return err
	}
	err = retryTimeDuration(insertFunc, d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return fmt.Errorf("Error, failed to insert "+
			"user %s into instance %s: %s", name, instance, err)
	}

	// This will include a double-slash (//) for postgres instances,
	// for which user.Host is an empty string.  That's okay.
	d.SetId(fmt.Sprintf("%s/%s/%s", user.Name, user.Host, user.Instance))

	err = sqlAdminOperationWaitTime(config, op, project, "Insert User", userAgent, d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return fmt.Errorf("Error, failure waiting for insertion of %s "+
			"into %s: %s", name, instance, err)
	}

	return resourceSqlUserRead(d, meta)
}

func resourceSqlUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instance := d.Get("instance").(string)
	name := d.Get("name").(string)
	host := d.Get("host").(string)

	var users *sqladmin.UsersListResponse
	err = nil
	err = retryTime(func() error {
		users, err = config.NewSqlAdminClient(userAgent).Users.List(project, instance).Do()
		return err
	}, 5)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SQL User %q in instance %q", name, instance))
	}

	var user *sqladmin.User
	databaseInstance, err := config.NewSqlAdminClient(userAgent).Instances.Get(project, instance).Do()
	if err != nil {
		return err
	}

	for _, currentUser := range users.Items {
		if !strings.Contains(databaseInstance.DatabaseVersion, "POSTGRES") {
			name = strings.Split(name, "@")[0]
		}

		if currentUser.Name == name {
			// Host can only be empty for postgres instances,
			// so don't compare the host if the API host is empty.
			if host == "" || currentUser.Host == host {
				user = currentUser
				break
			}
		}
	}

	if user == nil {
		log.Printf("[WARN] Removing SQL User %q because it's gone", d.Get("name").(string))
		d.SetId("")

		return nil
	}

	if err := d.Set("host", user.Host); err != nil {
		return fmt.Errorf("Error setting host: %s", err)
	}
	if err := d.Set("instance", user.Instance); err != nil {
		return fmt.Errorf("Error setting instance: %s", err)
	}
	if err := d.Set("name", user.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("type", user.Type); err != nil {
		return fmt.Errorf("Error setting type: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("disabled", user.SqlserverUserDetails.Disabled); err != nil {
		return fmt.Errorf("Error setting disabled: %s", err)
	}
	if err := d.Set("server_roles", user.SqlserverUserDetails.ServerRoles); err != nil {
		return fmt.Errorf("Error setting server_roles: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", user.Name, user.Host, user.Instance))
	return nil
}

func resourceSqlUserUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	if d.HasChange("password") {
		project, err := getProject(d, config)
		if err != nil {
			return err
		}

		name := d.Get("name").(string)
		instance := d.Get("instance").(string)
		password := d.Get("password").(string)
		host := d.Get("host").(string)
		disabled := d.Get("disabled").(bool)
		serverRoles := expandStringArray(d.Get("server_roles").(string))

		user := &sqladmin.User{
			Name:     name,
			Instance: instance,
			Password: password,
			SqlserverUserDetails: &sqladmin.SqlServerUserDetails{
				Disabled:    disabled,
				ServerRoles: serverRoles,
			},
		}

		mutexKV.Lock(instanceMutexKey(project, instance))
		defer mutexKV.Unlock(instanceMutexKey(project, instance))
		var op *sqladmin.Operation
		updateFunc := func() error {
			op, err = config.NewSqlAdminClient(userAgent).Users.Update(project, instance, user).Host(host).Name(name).Do()
			return err
		}
		err = retryTimeDuration(updateFunc, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return fmt.Errorf("Error, failed to update"+
				"user %s into user %s: %s", name, instance, err)
		}

		err = sqlAdminOperationWaitTime(config, op, project, "Insert User", userAgent, d.Timeout(schema.TimeoutUpdate))

		if err != nil {
			return fmt.Errorf("Error, failure waiting for update of %s "+
				"in %s: %s", name, instance, err)
		}

		return resourceSqlUserRead(d, meta)
	}

	return nil
}

func resourceSqlUserDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if deletionPolicy := d.Get("deletion_policy"); deletionPolicy == "ABANDON" {
		// Allows for user to be abandoned without deletion to avoid deletion failing
		// for Postgres users in some circumstances due to existing SQL roles
		return nil
	}

	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	host := d.Get("host").(string)
	instance := d.Get("instance").(string)

	mutexKV.Lock(instanceMutexKey(project, instance))
	defer mutexKV.Unlock(instanceMutexKey(project, instance))

	var op *sqladmin.Operation
	err = retryTimeDuration(func() error {
		op, err = config.NewSqlAdminClient(userAgent).Users.Delete(project, instance).Host(host).Name(name).Do()
		if err != nil {
			return err
		}

		if err := sqlAdminOperationWaitTime(config, op, project, "Delete User", userAgent, d.Timeout(schema.TimeoutDelete)); err != nil {
			return err
		}
		return nil
	}, d.Timeout(schema.TimeoutDelete), isSqlOperationInProgressError, isSqlInternalError)

	if err != nil {
		return fmt.Errorf("Error, failed to delete"+
			"user %s in instance %s: %s", name,
			instance, err)
	}

	return nil
}

func resourceSqlUserImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) == 3 {
		if err := d.Set("project", parts[0]); err != nil {
			return nil, fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("instance", parts[1]); err != nil {
			return nil, fmt.Errorf("Error setting instance: %s", err)
		}
		if err := d.Set("name", parts[2]); err != nil {
			return nil, fmt.Errorf("Error setting name: %s", err)
		}
	} else if len(parts) == 4 {
		if err := d.Set("project", parts[0]); err != nil {
			return nil, fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("instance", parts[1]); err != nil {
			return nil, fmt.Errorf("Error setting instance: %s", err)
		}
		if err := d.Set("host", parts[2]); err != nil {
			return nil, fmt.Errorf("Error setting host: %s", err)
		}
		if err := d.Set("name", parts[3]); err != nil {
			return nil, fmt.Errorf("Error setting name: %s", err)
		}
	} else {
		return nil, fmt.Errorf("Invalid specifier. Expecting {project}/{instance}/{name} for postgres instance and {project}/{instance}/{host}/{name} for MySQL instance")
	}

	return []*schema.ResourceData{d}, nil
}
