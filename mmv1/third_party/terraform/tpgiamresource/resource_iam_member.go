package tpgiamresource

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

func iamMemberCaseDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	isCaseSensitive := tpgresource.IamPrincipalIsCaseSensitive(old) || tpgresource.IamPrincipalIsCaseSensitive(new)
	if isCaseSensitive {
		return old == new
	}
	return tpgresource.CaseDiffSuppress(k, old, new, d)
}

func validateIAMMember(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
	}

	if matched, err := regexp.MatchString("^deleted", v); err != nil {
		return nil, []error{fmt.Errorf("error validating %s: %v", k, err)}
	} else if matched {
		return nil, []error{fmt.Errorf("invalid value for %s (Terraform does not support IAM members for deleted principals)", k)}
	}

	if matched, err := regexp.MatchString("(.+:.+|projectOwners|projectReaders|projectWriters|allUsers|allAuthenticatedUsers)", v); err != nil {
		return nil, []error{fmt.Errorf("error validating %s: %v", k, err)}
	} else if !matched {
		return nil, []error{fmt.Errorf("invalid value \"%s\" for %s (IAM members must have one of the values outlined here: https://cloud.google.com/billing/docs/reference/rest/v1/Policy#Binding)", v, k)}
	}
	return nil, nil
}

var IamMemberBaseSchema = map[string]*schema.Schema{
	"role": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"member": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: iamMemberCaseDiffSuppress,
		ValidateFunc:     validateIAMMember,
	},
	"condition": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"expression": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"title": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"description": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	},
	"etag": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var IamMemberBaseIdentitySchema = map[string]*schema.Schema{
	"member": {
		Type:              schema.TypeString,
		RequiredForImport: true,
	},
	"role": {
		Type:              schema.TypeString,
		RequiredForImport: true,
	},
	"condition_title": {
		Type:              schema.TypeString,
		OptionalForImport: true,
	},
}

func iamMemberImport(newUpdaterFunc NewResourceIamUpdaterFunc, resourceIdParser ResourceIdParserFunc, enableResourceIdentity bool) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if resourceIdParser == nil {
			return nil, errors.New("Import not supported for this IAM resource.")
		}

		if enableResourceIdentity && d.Id() == "" {
			identity, err := d.Identity()
			if err != nil {
				return nil, err
			}
			d.SetId(identity.Get("project").(string) + " " + identity.Get("role").(string) + " " + identity.Get("member").(string) + " " + identity.Get("condition_title").(string))
		}

		config := m.(*transport_tpg.Config)
		s := strings.Fields(d.Id())
		var id, role, member string
		if len(s) < 3 {
			d.SetId("")
			return nil, fmt.Errorf("Wrong number of parts to Member id %s; expected 'resource_name role member [condition_title]'.", s)
		}

		var conditionTitle string
		if len(s) == 3 {
			id, role, member = s[0], s[1], s[2]
		} else {
			// condition titles can have any characters in them, so re-join the split string
			id, role, member, conditionTitle = s[0], s[1], s[2], strings.Join(s[3:], " ")
		}

		// Set the ID only to the first part so all IAM types can share the same ResourceIdParserFunc.
		d.SetId(id)
		if err := d.Set("role", role); err != nil {
			return nil, fmt.Errorf("Error setting role: %s", err)
		}
		if err := d.Set("member", tpgresource.NormalizeIamPrincipalCasing(member)); err != nil {
			return nil, fmt.Errorf("Error setting member: %s", err)
		}

		err := resourceIdParser(d, config, enableResourceIdentity)
		if err != nil {
			return nil, err
		}

		// Set the ID again so that the ID matches the ID it would have if it had been created via TF.
		// Use the current ID in case it changed in the ResourceIdParserFunc.
		d.SetId(d.Id() + "/" + role + "/" + tpgresource.NormalizeIamPrincipalCasing(member))

		// Read the upstream policy so we can set the full condition.
		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return nil, err
		}
		p, err := iamPolicyReadWithRetry(updater)
		if err != nil {
			return nil, err
		}
		var binding *cloudresourcemanager.Binding
		for _, b := range p.Bindings {
			if b.Role == role && conditionKeyFromCondition(b.Condition).Title == conditionTitle {
				containsMember := false
				for _, m := range b.Members {
					if strings.ToLower(m) == strings.ToLower(member) {
						containsMember = true
					}
				}
				if !containsMember {
					continue
				}

				if binding != nil {
					return nil, fmt.Errorf("Cannot import IAM member with condition title %q, it matches multiple conditions", conditionTitle)
				}
				binding = b
			}
		}
		if binding == nil {
			return nil, fmt.Errorf("Cannot find binding for %q with role %q, member %q, and condition title %q", updater.DescribeResource(), role, member, conditionTitle)
		}

		if err := d.Set("condition", FlattenIamCondition(binding.Condition)); err != nil {
			return nil, fmt.Errorf("Error setting condition: %s", err)
		}
		if k := conditionKeyFromCondition(binding.Condition); !k.Empty() {
			d.SetId(d.Id() + "/" + k.String())
		}

		return []*schema.ResourceData{d}, nil
	}
}

func ConvertToIdentitySchema(parentSchema map[string]*schema.Schema) map[string]*schema.Schema {
	identitySchema := make(map[string]*schema.Schema)
	for k, v := range parentSchema {
		identitySchema[k] = &schema.Schema{
			Type: v.Type,
		}
		// If the field has RequiredForImport or OptionalForImport set, preserve them
		if v.RequiredForImport {
			identitySchema[k].RequiredForImport = true
		}
		if v.OptionalForImport {
			identitySchema[k].OptionalForImport = true
		}
		// If not explicitly set, infer from Required+ForceNew pattern
		if !v.RequiredForImport && !v.OptionalForImport {
			if v.Required && v.ForceNew {
				identitySchema[k].RequiredForImport = true
			} else if v.Optional && v.ForceNew {
				identitySchema[k].OptionalForImport = true
			}
		}
	}
	return identitySchema
}

func ResourceIamMember(parentSpecificSchema map[string]*schema.Schema, newUpdaterFunc NewResourceIamUpdaterFunc, resourceIdParser ResourceIdParserFunc, options ...func(*IamSettings)) *schema.Resource {
	settings := NewIamSettings(options...)

	createTimeOut := time.Duration(settings.CreateTimeOut) * time.Minute

	resourceSchema := &schema.Resource{
		Create: resourceIamMemberCreate(newUpdaterFunc, settings.EnableBatching, settings.EnableResourceIdentity),
		Read:   resourceIamMemberRead(newUpdaterFunc, settings.EnableResourceIdentity),
		Delete: resourceIamMemberDelete(newUpdaterFunc, settings.EnableBatching, settings.EnableResourceIdentity),

		// if non-empty, this will be used to send a deprecation message when the
		// resource is used.
		DeprecationMessage: settings.DeprecationMessage,

		Schema:         tpgresource.MergeSchemas(IamMemberBaseSchema, parentSpecificSchema),
		SchemaVersion:  settings.SchemaVersion,
		StateUpgraders: settings.StateUpgraders,
		Importer: &schema.ResourceImporter{
			State: iamMemberImport(newUpdaterFunc, resourceIdParser, settings.EnableResourceIdentity),
		},
		UseJSONNumber: true,
	}

	if settings.EnableResourceIdentity {
		resourceSchema.Identity = &schema.ResourceIdentity{
			Version: 1,
			SchemaFunc: func() map[string]*schema.Schema {
				return tpgresource.MergeSchemas(IamMemberBaseIdentitySchema, ConvertToIdentitySchema(parentSpecificSchema))
			},
		}
	}

	if createTimeOut > 0 {
		resourceSchema.Timeouts = &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(createTimeOut),
		}
	}
	return resourceSchema
}

func getResourceIamMember(d *schema.ResourceData) *cloudresourcemanager.Binding {
	b := &cloudresourcemanager.Binding{
		Members: []string{d.Get("member").(string)},
		Role:    d.Get("role").(string),
	}
	if c := ExpandIamCondition(d.Get("condition")); c != nil {
		b.Condition = c
	}
	return b
}

func resourceIamMemberCreate(newUpdaterFunc NewResourceIamUpdaterFunc, enableBatching bool, enableResourceIdentity bool) schema.CreateFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)

		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		memberBind := getResourceIamMember(d)
		modifyF := func(ep *cloudresourcemanager.Policy) error {
			// Merge the bindings together
			ep.Bindings = MergeBindings(append(ep.Bindings, memberBind))
			ep.Version = IamPolicyVersion
			return nil
		}
		if enableBatching {
			err = BatchRequestModifyIamPolicy(updater, modifyF, config,
				fmt.Sprintf("Create IAM Members %s %+v for %s", memberBind.Role, memberBind.Members[0], updater.DescribeResource()))
		} else {
			err = iamPolicyReadModifyWrite(updater, modifyF)
		}
		if err != nil {
			return err
		}
		d.SetId(updater.GetResourceId() + "/" + memberBind.Role + "/" + tpgresource.NormalizeIamPrincipalCasing(memberBind.Members[0]))
		if k := conditionKeyFromCondition(memberBind.Condition); !k.Empty() {
			d.SetId(d.Id() + "/" + k.String())
		}

		if enableResourceIdentity {
			identity, err := d.Identity()
			if err != nil {
				return err
			}
			identity.Set("project", d.Get("project").(string))
			identity.Set("role", memberBind.Role)
			identity.Set("member", tpgresource.NormalizeIamPrincipalCasing(memberBind.Members[0]))
			if memberBind.Condition != nil {
				identity.Set("condition_title", memberBind.Condition.Title)
			}
		}

		return resourceIamMemberRead(newUpdaterFunc, enableResourceIdentity)(d, meta)
	}
}

func resourceIamMemberRead(newUpdaterFunc NewResourceIamUpdaterFunc, enableResourceIdentity bool) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)

		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		eMember := getResourceIamMember(d)
		eCondition := conditionKeyFromCondition(eMember.Condition)
		p, err := iamPolicyReadWithRetry(updater)
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Resource %q with IAM Member: Role %q Member %q", updater.DescribeResource(), eMember.Role, eMember.Members[0]))
		}
		log.Print(spew.Sprintf("[DEBUG]: Retrieved policy for %s: %#v\n", updater.DescribeResource(), p))
		log.Printf("[DEBUG]: Looking for binding with role %q and condition %#v", eMember.Role, eCondition)

		var binding *cloudresourcemanager.Binding
		for _, b := range p.Bindings {
			if b.Role == eMember.Role && conditionKeyFromCondition(b.Condition) == eCondition {
				binding = b
				break
			}
		}

		if binding == nil {
			log.Printf("[DEBUG]: Binding for role %q with condition %#v does not exist in policy of %s, removing member %q from state.", eMember.Role, eCondition, updater.DescribeResource(), eMember.Members[0])
			d.SetId("")
			return nil
		}

		log.Printf("[DEBUG]: Looking for member %q in found binding", eMember.Members[0])
		var member string
		for _, m := range binding.Members {
			if strings.ToLower(m) == strings.ToLower(eMember.Members[0]) {
				member = m
			}
		}

		if member == "" {
			log.Printf("[DEBUG]: Member %q for binding for role %q with condition %#v does not exist in policy of %s, removing from state.", eMember.Members[0], eMember.Role, eCondition, updater.DescribeResource())
			d.SetId("")
			return nil
		}

		if err := d.Set("etag", p.Etag); err != nil {
			return fmt.Errorf("Error setting etag: %s", err)
		}
		if err := d.Set("member", member); err != nil {
			return fmt.Errorf("Error setting member: %s", err)
		}
		if err := d.Set("role", binding.Role); err != nil {
			return fmt.Errorf("Error setting role: %s", err)
		}
		if err := d.Set("condition", FlattenIamCondition(binding.Condition)); err != nil {
			return fmt.Errorf("Error setting condition: %s", err)
		}

		if enableResourceIdentity {
			identity, err := d.Identity()
			if err != nil {
				return err
			}
			identity.Set("project", d.Get("project").(string))
			identity.Set("role", binding.Role)
			identity.Set("member", tpgresource.NormalizeIamPrincipalCasing(eMember.Members[0]))
			if binding.Condition != nil {
				identity.Set("condition_title", binding.Condition.Title)
			}
		}
		return nil
	}
}

func resourceIamMemberDelete(newUpdaterFunc NewResourceIamUpdaterFunc, enableBatching bool, enableResourceIdentity bool) schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)

		updater, err := newUpdaterFunc(d, config)
		if err != nil {
			return err
		}

		memberBind := getResourceIamMember(d)
		modifyF := func(ep *cloudresourcemanager.Policy) error {
			// Merge the bindings together
			ep.Bindings = subtractFromBindings(ep.Bindings, memberBind)
			return nil
		}
		if enableBatching {
			err = BatchRequestModifyIamPolicy(updater, modifyF, config,
				fmt.Sprintf("Delete IAM Members %s %s for %q", memberBind.Role, memberBind.Members[0], updater.DescribeResource()))
		} else {
			err = iamPolicyReadModifyWrite(updater, modifyF)
		}
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Resource %s for IAM Member (role %q, %q)", updater.GetResourceId(), memberBind.Members[0], memberBind.Role))
		}
		return resourceIamMemberRead(newUpdaterFunc, enableResourceIdentity)(d, meta)
	}
}
