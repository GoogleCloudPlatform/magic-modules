package google

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"crypto/md5"
	"encoding/base64"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

func resourceStorageBucketObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageBucketObjectCreate,
		Read:   resourceStorageBucketObjectRead,
		Update: resourceStorageBucketObjectUpdate,
		Delete: resourceStorageBucketObjectDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the containing bucket.`,
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the object. If you're interpolating the name of this object, see output_name instead.`,
			},

			"cache_control": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: `Cache-Control directive to specify caching behavior of object data. If omitted and object is accessible to all anonymous users, the default will be public, max-age=3600`,
			},

			"content_disposition": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: `Content-Disposition of the object data.`,
			},

			"content_encoding": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: `Content-Encoding of the object data.`,
			},

			"content_language": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: `Content-Language of the object data.`,
			},

			"content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `Content-Type of the object data. Defaults to "application/octet-stream" or "text/plain; charset=utf-8".`,
			},

			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source"},
				Sensitive:     true,
				Description:   `Data as string to be uploaded. Must be defined if source is not. Note: The content field is marked as sensitive. To view the raw contents of the object, please define an output.`,
			},

			"crc32c": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Base 64 CRC32 hash of the uploaded data.`,
			},

			"md5hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Base 64 MD5 hash of the uploaded data.`,
			},

			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"content"},
				Description:   `A path to the data you want to upload. Must be defined if content is not.`,
			},

			// Detect changes to local file or changes made outside of Terraform to the file stored on the server.
			"detect_md5hash": {
				Type: schema.TypeString,
				// This field is not Computed because it needs to trigger a diff.
				Optional: true,
				ForceNew: true,
				// Makes the diff message nicer:
				// detect_md5hash:       "1XcnP/iFw/hNrbhXi7QTmQ==" => "different hash" (forces new resource)
				// Instead of the more confusing:
				// detect_md5hash:       "1XcnP/iFw/hNrbhXi7QTmQ==" => "" (forces new resource)
				Default: "different hash",
				// 1. Compute the md5 hash of the local file
				// 2. Compare the computed md5 hash with the hash stored in Cloud Storage
				// 3. Don't suppress the diff iff they don't match
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					localMd5Hash := ""
					if source, ok := d.GetOkExists("source"); ok {
						localMd5Hash = getFileMd5Hash(source.(string))
					}

					if content, ok := d.GetOkExists("content"); ok {
						localMd5Hash = getContentMd5Hash([]byte(content.(string)))
					}

					// If `source` or `content` is dynamically set, both field will be empty.
					// We should not suppress the diff to avoid the following error:
					// 'Mismatch reason: extra attributes: detect_md5hash'
					if localMd5Hash == "" {
						return false
					}

					// `old` is the md5 hash we retrieved from the server in the ReadFunc
					if old != localMd5Hash {
						return false
					}

					return true
				},
			},

			"storage_class": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The StorageClass of the new bucket object. Supported values include: MULTI_REGIONAL, REGIONAL, NEARLINE, COLDLINE, ARCHIVE. If not provided, this defaults to the bucket's default storage class or to a standard class.`,
			},
			"kms_key_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: compareCryptoKeyVersions,
				Description:      `Resource name of the Cloud KMS key that will be used to encrypt the object. Overrides the object metadata's kmsKeyName value, if any.`,
			},
			"event_based_hold": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether an object is under event-based hold. Event-based hold is a way to retain objects until an event occurs, which is signified by the hold's release (i.e. this value is set to false). After being released (set to false), such objects will be subject to bucket-level retention (if any).`,
			},
			"temporary_hold": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether an object is under temporary hold. While this flag is set to true, the object is protected against deletion and overwrites.`,
			},
			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `User-provided metadata, in key/value pairs.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `A url reference to this object.`,
			},

			// https://github.com/hashicorp/terraform/issues/19052
			"output_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the object. Use this field in interpolations with google_storage_object_acl to recreate google_storage_object_acl resources when your google_storage_bucket_object is recreated.`,
			},

			"media_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `A url reference to download this object.`,
			},
		},
		UseJSONNumber: true,
	}
}

func objectGetID(attrs *storage.ObjectAttrs) string {
	return attrs.Bucket + "-" + attrs.Name
}

func compareCryptoKeyVersions(_, old, new string, _ *schema.ResourceData) bool {
	// The API can return cryptoKeyVersions even though it wasn't specified.
	// format: projects/<project>/locations/<region>/keyRings/<keyring>/cryptoKeys/<key>/cryptoKeyVersions/1

	kmsKeyWithoutVersions := strings.Split(old, "/cryptoKeyVersions")[0]
	if kmsKeyWithoutVersions == new {
		return true
	}

	return false
}

func resourceStorageBucketObjectCreate(d *schema.ResourceData, meta interface{}) error {

	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)

	var media io.Reader

	if v, ok := d.GetOk("source"); ok {
		var err error
		media, err = os.Open(v.(string))
		if err != nil {
			return err
		}
	} else if v, ok := d.GetOk("content"); ok {
		media = bytes.NewReader([]byte(v.(string)))
	} else {
		return fmt.Errorf("Error, either \"content\" or \"source\" must be specified")
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	wc := client.Bucket(bucket).Object(name).NewWriter(ctx)

	if err != nil {
		log.Fatal(err)
	}

	if v, ok := d.GetOk("cache_control"); ok {
		wc.CacheControl = v.(string)
	}

	if v, ok := d.GetOk("content_disposition"); ok {
		wc.ContentDisposition = v.(string)
	}

	if v, ok := d.GetOk("content_encoding"); ok {
		wc.ContentEncoding = v.(string)
	}

	if v, ok := d.GetOk("content_language"); ok {
		wc.ContentLanguage = v.(string)
	}

	if v, ok := d.GetOk("content_type"); ok {
		wc.ContentType = v.(string)
	}

	if v, ok := d.GetOk("metadata"); ok {
		wc.Metadata = convertStringMap(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("storage_class"); ok {
		wc.StorageClass = v.(string)
	}

	if v, ok := d.GetOk("kms_key_name"); ok {
		wc.KMSKeyName = v.(string)
	}

	if v, ok := d.GetOk("event_based_hold"); ok {
		wc.EventBasedHold = v.(bool)
	}

	if v, ok := d.GetOk("temporary_hold"); ok {
		wc.TemporaryHold = v.(bool)
	}

	if _, err = io.Copy(wc, media); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return resourceStorageBucketObjectRead(d, meta)
}

func resourceStorageBucketObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)

	object := client.Bucket(bucket).Object(name)
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{}

	if d.HasChange("event_based_hold") {
		v := d.Get("event_based_hold")
		objectAttrsToUpdate.EventBasedHold = v.(bool)
	}

	if d.HasChange("temporary_hold") {
		v := d.Get("temporary_hold")
		objectAttrsToUpdate.TemporaryHold = v.(bool)
	}

	if _, err := object.Update(ctx, objectAttrsToUpdate); err != nil {
		return fmt.Errorf("Object(%q).Update: %v", name, err)
	}

	return nil
}

func resourceStorageBucketObjectRead(d *schema.ResourceData, meta interface{}) error {

	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	object := client.Bucket(bucket).Object(name)
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).Attrs: %v", name, err)
	}

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Storage Bucket Object %q", d.Get("name").(string)))
	}

	if err := d.Set("md5hash", base64.StdEncoding.EncodeToString(attrs.MD5)); err != nil {
		return fmt.Errorf("Error setting md5hash: %s", err)
	}
	if err := d.Set("detect_md5hash", base64.StdEncoding.EncodeToString(attrs.MD5)); err != nil {
		return fmt.Errorf("Error setting detect_md5hash: %s", err)
	}
	if err := d.Set("crc32c", fmt.Sprint(attrs.CRC32C)); err != nil {
		return fmt.Errorf("Error setting crc32c: %s", err)
	}
	if err := d.Set("cache_control", attrs.CacheControl); err != nil {
		return fmt.Errorf("Error setting cache_control: %s", err)
	}
	if err := d.Set("content_disposition", attrs.ContentDisposition); err != nil {
		return fmt.Errorf("Error setting content_disposition: %s", err)
	}
	if err := d.Set("content_encoding", attrs.ContentEncoding); err != nil {
		return fmt.Errorf("Error setting content_encoding: %s", err)
	}
	if err := d.Set("content_language", attrs.ContentLanguage); err != nil {
		return fmt.Errorf("Error setting content_language: %s", err)
	}
	if err := d.Set("content_type", attrs.ContentType); err != nil {
		return fmt.Errorf("Error setting content_type: %s", err)
	}
	if err := d.Set("storage_class", attrs.StorageClass); err != nil {
		return fmt.Errorf("Error setting storage_class: %s", err)
	}
	if err := d.Set("kms_key_name", attrs.KMSKeyName); err != nil {
		return fmt.Errorf("Error setting kms_key_name: %s", err)
	}
	if err := d.Set("output_name", attrs.Name); err != nil {
		return fmt.Errorf("Error setting output_name: %s", err)
	}
	if err := d.Set("metadata", attrs.Metadata); err != nil {
		return fmt.Errorf("Error setting metadata: %s", err)
	}
	if err := d.Set("media_link", attrs.MediaLink); err != nil {
		return fmt.Errorf("Error setting media_link: %s", err)
	}
	if err := d.Set("event_based_hold", attrs.EventBasedHold); err != nil {
		return fmt.Errorf("Error setting event_based_hold: %s", err)
	}
	if err := d.Set("temporary_hold", attrs.TemporaryHold); err != nil {
		return fmt.Errorf("Error setting temporary_hold: %s", err)
	}

	d.SetId(objectGetID(attrs))

	return nil
}

func resourceStorageBucketObjectDelete(d *schema.ResourceData, meta interface{}) error {

	bucket := d.Get("bucket").(string)
	name := d.Get("name").(string)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	object := client.Bucket(bucket).Object(name)
	if err := object.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", name, err)
	}

	return nil
}

func getFileMd5Hash(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("[WARN] Failed to read source file %q. Cannot compute md5 hash for it.", filename)
		return ""
	}

	return getContentMd5Hash(data)
}

func getContentMd5Hash(content []byte) string {
	h := md5.New()
	if _, err := h.Write(content); err != nil {
		log.Printf("[WARN] Failed to compute md5 hash for content: %v", err)
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
