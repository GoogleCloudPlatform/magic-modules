package google

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/iterator"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/genproto/protobuf/field_mask"
)

func resourceKmsCryptoKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceKmsCryptoKeyCreate,
		Read:   resourceKmsCryptoKeyRead,
		Update: resourceKmsCryptoKeyUpdate,
		Delete: resourceKmsCryptoKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"key_ring": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: kmsCryptoKeyRingsEquivalent,
			},

			"rotation_period": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsCryptoKeyRotationPeriodsEquivalent,
			},

			"purpose": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				Default:          "encrypt_decrypt",
				ValidateFunc:     validation.StringInSlice(kmsPurposeNames(), true),
				StateFunc:        caseLowerStateFunc,
				DiffSuppressFunc: caseDiffSuppress,
			},

			"version_template": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:             schema.TypeString,
							ForceNew:         true,
							Optional:         true,
							Default:          "symmetric_encryption",
							ValidateFunc:     validation.StringInSlice(kmsAlgorithmNames(), true),
							StateFunc:        caseLowerStateFunc,
							DiffSuppressFunc: caseDiffSuppress,
						},

						"protection_level": {
							Type:             schema.TypeString,
							ForceNew:         true,
							Optional:         true,
							Default:          "software",
							ValidateFunc:     validation.StringInSlice(kmsProtectionLevelNames(), true),
							StateFunc:        caseLowerStateFunc,
							DiffSuppressFunc: caseDiffSuppress,
						},
					},
				},
			},

			"next_rotation_rfc3339": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"next_rotation_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"rotation_period_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKmsCryptoKeyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	tfKeyRing, err := parseKmsKeyRingId(d.Get("key_ring").(string), config)
	if err != nil {
		return err
	}

	cryptoKeyId := &kmsCryptoKeyId{
		KeyRingId: *tfKeyRing,
		Name:      d.Get("name").(string),
	}

	key := &kmspb.CryptoKey{
		Purpose:         kmsPurposes[strings.ToLower(d.Get("purpose").(string))],
		VersionTemplate: expandVersionTemplate(d.Get("version_template").([]interface{})),
	}

	if d.Get("rotation_period") != "" {
		rotationPeriodStr := d.Get("rotation_period").(string)
		rotationSchedule, nextRotationTime, err := kmsCryptoKeyNextRotation(rotationPeriodStr, time.Now())
		if err != nil {
			return fmt.Errorf("Failed to parse rotation_period: %s", err)
		}

		key.RotationSchedule = rotationSchedule
		key.NextRotationTime = nextRotationTime
	}

	ctx := context.Background()
	cryptoKey, err := config.clientKms.CreateCryptoKey(ctx, &kmspb.CreateCryptoKeyRequest{
		Parent:      tfKeyRing.keyRingId(),
		CryptoKeyId: cryptoKeyId.Name,
		CryptoKey:   key,
	})
	if err != nil {
		return fmt.Errorf("Error creating CryptoKey: %s", err)
	}

	log.Printf("[DEBUG] Created CryptoKey %s", cryptoKey.Name)

	d.SetId(cryptoKeyId.cryptoKeyId())

	return resourceKmsCryptoKeyRead(d, meta)
}

func resourceKmsCryptoKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}

	var key kmspb.CryptoKey
	var updatedFields []string

	key.Name = cryptoKeyId.cryptoKeyId()

	if d.HasChange("rotation_period") {
		if d.Get("rotation_period") != "" {
			rotationPeriodStr := d.Get("rotation_period").(string)
			rotationSchedule, nextRotationTime, err := kmsCryptoKeyNextRotation(rotationPeriodStr, time.Now())
			if err != nil {
				return fmt.Errorf("Failed to parse rotation_period: %s", err)
			}

			key.RotationSchedule = rotationSchedule
			key.NextRotationTime = nextRotationTime
		}

		// If the rotation period changed but was empty, still send the update
		// fields which triggers a removal of the rotation period.
		updatedFields = append(updatedFields, "rotation_period")
		updatedFields = append(updatedFields, "next_rotation_time")
	}

	ctx := context.Background()
	cryptoKey, err := config.clientKms.UpdateCryptoKey(ctx, &kmspb.UpdateCryptoKeyRequest{
		CryptoKey: &key,
		UpdateMask: &field_mask.FieldMask{
			Paths: updatedFields,
		},
	})
	if err != nil {
		return fmt.Errorf("Error updating CryptoKey: %s", err.Error())
	}

	log.Printf("[DEBUG] Updated CryptoKey %s", cryptoKey.Name)

	d.SetId(cryptoKeyId.cryptoKeyId())

	return resourceKmsCryptoKeyRead(d, meta)
}

func resourceKmsCryptoKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Executing read for KMS CryptoKey %s", cryptoKeyId.cryptoKeyId())

	ctx := context.Background()
	cryptoKey, err := config.clientKms.GetCryptoKey(ctx, &kmspb.GetCryptoKeyRequest{
		Name: cryptoKeyId.cryptoKeyId(),
	})
	if err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}
	d.Set("name", cryptoKeyId.Name)
	d.Set("key_ring", cryptoKeyId.KeyRingId.terraformId())

	if rotationPeriod := cryptoKey.GetRotationPeriod(); rotationPeriod != nil {
		sec := rotationPeriod.GetSeconds()
		d.Set("rotation_period", fmt.Sprintf("%ds", sec))
		d.Set("rotation_period_seconds", sec)
	} else {
		d.Set("rotation_period", "")
		d.Set("rotation_period_seconds", 0)
	}

	if nextRotationTime := cryptoKey.GetNextRotationTime(); nextRotationTime != nil {
		sec := nextRotationTime.GetSeconds()
		d.Set("next_rotation_rfc3339", time.Unix(sec, 0).Format(time.RFC3339))
		d.Set("next_rotation_seconds", sec)
	} else {
		d.Set("next_rotation_rfc3339", "")
		d.Set("next_rotation_seconds", 0)
	}

	d.Set("purpose", kmsPurposeToString(cryptoKey.Purpose))
	d.Set("self_link", cryptoKey.Name)

	if err = d.Set("version_template", flattenVersionTemplate(cryptoKey.VersionTemplate)); err != nil {
		return fmt.Errorf("Error setting version_template in state: %s", err.Error())
	}

	d.SetId(cryptoKeyId.cryptoKeyId())

	return nil
}

/*
	Because KMS CryptoKey resources cannot be deleted on GCP, we are only going to remove it from state
	and destroy all its versions, rendering the key useless for encryption and decryption of data.
	Re-creation of this resource through Terraform will produce an error.
*/

func resourceKmsCryptoKeyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cryptoKeyId, err := parseKmsCryptoKeyId(d.Id(), config)
	if err != nil {
		return err
	}

	log.Printf("[WARN] Cloud KMS CryptoKey resources cannot be deleted. "+
		"The CryptoKey %s will be removed from Terraform state, and all its "+
		"CryptoKeyVersions will be destroyed, but it will still be present on "+
		"the server.", cryptoKeyId.cryptoKeyId())

	var ckvs []string
	ctx := context.Background()
	it := config.clientKms.ListCryptoKeyVersions(ctx, &kmspb.ListCryptoKeyVersionsRequest{
		Parent: cryptoKeyId.cryptoKeyId(),
	})
	for {
		resp, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return fmt.Errorf("Failed to list crypto key versions: %s", err)
		}

		if resp.State != kmspb.CryptoKeyVersion_DESTROYED &&
			resp.State != kmspb.CryptoKeyVersion_DESTROY_SCHEDULED {
			ckvs = append(ckvs, resp.Name)
		}
	}

	var mu sync.Mutex
	var errs *multierror.Error
	pool := workerpool.New(runtime.NumCPU() - 1)
	for _, ckv := range ckvs {
		ckv := ckv

		pool.Submit(func() {
			if _, err := config.clientKms.DestroyCryptoKeyVersion(ctx, &kmspb.DestroyCryptoKeyVersionRequest{
				Name: ckv,
			}); err != nil {
				mu.Lock()
				errs = multierror.Append(errs, fmt.Errorf("Failed to destroy crypto key version %s: %s", ckv, err))
				mu.Unlock()
			}
		})
	}
	pool.StopWait()

	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func expandVersionTemplate(configured []interface{}) *kmspb.CryptoKeyVersionTemplate {
	if configured == nil || len(configured) == 0 {
		return nil
	}

	data := configured[0].(map[string]interface{})
	return &kmspb.CryptoKeyVersionTemplate{
		Algorithm:       kmsAlgorithms[strings.ToLower(data["algorithm"].(string))],
		ProtectionLevel: kmsProtectionLevels[strings.ToLower(data["protection_level"].(string))],
	}
}

func flattenVersionTemplate(versionTemplate *kmspb.CryptoKeyVersionTemplate) []map[string]interface{} {
	if versionTemplate == nil {
		return nil
	}

	versionTemplateSchema := make([]map[string]interface{}, 0, 1)
	data := map[string]interface{}{
		"algorithm":        kmsAlgorithmToString(versionTemplate.Algorithm),
		"protection_level": kmsProtectionLevelToString(versionTemplate.ProtectionLevel),
	}

	versionTemplateSchema = append(versionTemplateSchema, data)
	return versionTemplateSchema
}

// kmsCryptoKeyNextRotation accepts a period as parseable by time.ParseDuration
// and returns the appropriate protos for the rotation period and next rotation
// time. If an error occurs, it is returned and the protos will be nil.
func kmsCryptoKeyNextRotation(periodStr string, now time.Time) (*kmspb.CryptoKey_RotationPeriod, *timestamp.Timestamp, error) {
	dur, err := time.ParseDuration(periodStr)
	if err != nil {
		return nil, nil, err
	}

	rotationPeriod := &kmspb.CryptoKey_RotationPeriod{
		RotationPeriod: &duration.Duration{
			Seconds: int64(dur.Seconds()),
		},
	}

	nextRotationTime := &timestamp.Timestamp{
		Seconds: now.UTC().Add(dur).Unix(),
	}

	return rotationPeriod, nextRotationTime, nil
}

// kmsCryptoKeyRingsEquivalent determines if the two keyrings are equivalent,
// taking into account short vs long names.
//
// The following KMS key rings are equivalent:
//
//     projects/my-project/locations/us-east4/keyRings/my-keyring
//     my-project/us-east4/my-keyring
//
// This function compares the values to see if they are equivalent.
func kmsCryptoKeyRingsEquivalent(k, old, new string, d *schema.ResourceData) bool {
	keyRingIdWithSpecifiersRegex := regexp.MustCompile("^projects/(" + ProjectRegex + ")/locations/([a-z0-9-])+/keyRings/([a-zA-Z0-9_-]{1,63})$")
	normalizedKeyRingIdRegex := regexp.MustCompile("^(" + ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})$")
	if matches := keyRingIdWithSpecifiersRegex.FindStringSubmatch(new); matches != nil {
		normMatches := normalizedKeyRingIdRegex.FindStringSubmatch(old)
		return normMatches != nil && normMatches[1] == matches[1] && normMatches[2] == matches[2] && normMatches[3] == matches[3]
	}
	return false
}

// kmsCryptoKeyRotationPeriodsEquivalent determines if two rotation periods are
// equivalent. Since the API converts the duration into seconds, the user might
// specify 24h, but that will come back as 86400s. This checks if both values
// are equivalent and suppresses the diff if so.
func kmsCryptoKeyRotationPeriodsEquivalent(k, old, new string, d *schema.ResourceData) bool {
	oldD, err := time.ParseDuration(old)
	if err != nil {
		return false
	}

	newD, err := time.ParseDuration(new)
	if err != nil {
		return false
	}

	return oldD == newD
}

// kmsPurposes is the list of purposes to key types
var kmsPurposes = map[string]kmspb.CryptoKey_CryptoKeyPurpose{
	"asymmetric_decrypt": kmspb.CryptoKey_ASYMMETRIC_DECRYPT,
	"asymmetric_sign":    kmspb.CryptoKey_ASYMMETRIC_SIGN,
	"encrypt_decrypt":    kmspb.CryptoKey_ENCRYPT_DECRYPT,
	"unspecified":        kmspb.CryptoKey_CRYPTO_KEY_PURPOSE_UNSPECIFIED,
}

// kmsPurposeNames returns the list of key purposes.
func kmsPurposeNames() []string {
	list := make([]string, 0, len(kmsPurposes))
	for k := range kmsPurposes {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

// kmsPurposeToString accepts a kmspb and maps that to the user readable purpose.
// Instead of maintaining two maps, this iterates over the purposes map because
// N will always be ridiculously small.
func kmsPurposeToString(p kmspb.CryptoKey_CryptoKeyPurpose) string {
	for k, v := range kmsPurposes {
		if p == v {
			return k
		}
	}
	return "unspecified"
}

// kmsAlgorithms is the list of key algorithms.
var kmsAlgorithms = map[string]kmspb.CryptoKeyVersion_CryptoKeyVersionAlgorithm{
	"symmetric_encryption":         kmspb.CryptoKeyVersion_GOOGLE_SYMMETRIC_ENCRYPTION,
	"rsa_sign_pss_2048_sha256":     kmspb.CryptoKeyVersion_RSA_SIGN_PSS_2048_SHA256,
	"rsa_sign_pss_3072_sha256":     kmspb.CryptoKeyVersion_RSA_SIGN_PSS_3072_SHA256,
	"rsa_sign_pss_4096_sha256":     kmspb.CryptoKeyVersion_RSA_SIGN_PSS_4096_SHA256,
	"rsa_sign_pkcs1_2048_sha256":   kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_2048_SHA256,
	"rsa_sign_pkcs1_3072_sha256":   kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_3072_SHA256,
	"rsa_sign_pkcs1_4096_sha256":   kmspb.CryptoKeyVersion_RSA_SIGN_PKCS1_4096_SHA256,
	"rsa_decrypt_oaep_2048_sha256": kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_2048_SHA256,
	"rsa_decrypt_oaep_3072_sha256": kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_3072_SHA256,
	"rsa_decrypt_oaep_4096_sha256": kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_4096_SHA256,
	"ec_sign_p256_sha256":          kmspb.CryptoKeyVersion_EC_SIGN_P256_SHA256,
	"ec_sign_p384_sha384":          kmspb.CryptoKeyVersion_EC_SIGN_P384_SHA384,
}

// kmsAlgorithmNames returns the list of key algorithms.
func kmsAlgorithmNames() []string {
	list := make([]string, 0, len(kmsAlgorithms))
	for k := range kmsAlgorithms {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

// kmsAlgorithmToString accepts a kmspb and maps that to the user readable algorithm.
// Instead of maintaining two maps, this iterates over the algorithms map because
// N will always be ridiculously small.
func kmsAlgorithmToString(p kmspb.CryptoKeyVersion_CryptoKeyVersionAlgorithm) string {
	for k, v := range kmsAlgorithms {
		if p == v {
			return k
		}
	}
	return "unspecified"
}

// kmsProtectionLevels is the list of key protection levels.
var kmsProtectionLevels = map[string]kmspb.ProtectionLevel{
	"hsm":      kmspb.ProtectionLevel_HSM,
	"software": kmspb.ProtectionLevel_SOFTWARE,
}

// kmsProtectionLevelNames returns the list of key protection levels.
func kmsProtectionLevelNames() []string {
	list := make([]string, 0, len(kmsProtectionLevels))
	for k := range kmsProtectionLevels {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

// kmsProtectionLevelToString accepts a kmspb and maps that to the user readable algorithm.
// Instead of maintaining two maps, this iterates over the algorithms map because
// N will always be ridiculously small.
func kmsProtectionLevelToString(p kmspb.ProtectionLevel) string {
	for k, v := range kmsProtectionLevels {
		if p == v {
			return k
		}
	}
	return "unknown"
}

func parseKmsCryptoKeyId(id string, config *Config) (*kmsCryptoKeyId, error) {
	parts := strings.Split(id, "/")

	cryptoKeyIdRegex := regexp.MustCompile("^(" + ProjectRegex + ")/([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})/([a-zA-Z0-9_-]{1,63})$")
	cryptoKeyIdWithoutProjectRegex := regexp.MustCompile("^([a-z0-9-])+/([a-zA-Z0-9_-]{1,63})/([a-zA-Z0-9_-]{1,63})$")
	cryptoKeyRelativeLinkRegex := regexp.MustCompile("^projects/(" + ProjectRegex + ")/locations/([a-z0-9-]+)/keyRings/([a-zA-Z0-9_-]{1,63})/cryptoKeys/([a-zA-Z0-9_-]{1,63})$")

	if cryptoKeyIdRegex.MatchString(id) {
		return &kmsCryptoKeyId{
			KeyRingId: kmsKeyRingId{
				Project:  parts[0],
				Location: parts[1],
				Name:     parts[2],
			},
			Name: parts[3],
		}, nil
	}

	if cryptoKeyIdWithoutProjectRegex.MatchString(id) {
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{location}/{keyRingName}/{cryptoKeyName}` id format.")
		}

		return &kmsCryptoKeyId{
			KeyRingId: kmsKeyRingId{
				Project:  config.Project,
				Location: parts[0],
				Name:     parts[1],
			},
			Name: parts[2],
		}, nil
	}

	if parts := cryptoKeyRelativeLinkRegex.FindStringSubmatch(id); parts != nil {
		return &kmsCryptoKeyId{
			KeyRingId: kmsKeyRingId{
				Project:  parts[1],
				Location: parts[2],
				Name:     parts[3],
			},
			Name: parts[4],
		}, nil
	}
	return nil, fmt.Errorf("Invalid CryptoKey id format, expecting `{projectId}/{locationId}/{KeyringName}/{cryptoKeyName}` or `{locationId}/{keyRingName}/{cryptoKeyName}.`")
}

type kmsCryptoKeyId struct {
	KeyRingId kmsKeyRingId
	Name      string
}

func (s *kmsCryptoKeyId) cryptoKeyId() string {
	return fmt.Sprintf("%s/cryptoKeys/%s", s.KeyRingId.keyRingId(), s.Name)
}

func (s *kmsCryptoKeyId) terraformId() string {
	return fmt.Sprintf("%s/%s", s.KeyRingId.terraformId(), s.Name)
}
