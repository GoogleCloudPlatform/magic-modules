---
subcategory: "Cloud Storage"
description: |-
  A Google Cloud Storage Bucket Objects Lifecycle.
---

# google_storage_bucket_life_cycle_config

A Google Cloud Storage Bucket Objects Lifecycle.

## Example Usage - Storage Bucket Lifecycle Config


```hcl
resource "google_storage_bucket" "bucket" {
  name     = "my-bucket"
  location = "US"
}

resource "google_storage_bucket_life_cycle_config" "bucketlifecycle" {
    depends_on = [google_storage_bucket.bucket]
    bucket = google_storage_bucket.bucket.name
    lifecycle_rule {
        action {
          type = "Delete"
        }
        condition {
          age        = 10
        }
    }
}
```

## Argument Reference

The following arguments are supported:


* `bucket` -
  (Required)
  The name of the bucket that contains the folder.


- - -


* `lifecycle_rules` -
  (Optional)
  A nested object resource.
  Structure is [documented below](#nested_lifecycle_rules).


<a name="nested_lifecycle_rules"></a>The `lifecycle_rules` block supports:

* `rule` -
  (Required)
  A lifecycle management rule, which is made of an action to take
  and the condition(s) under which the action will be taken.
  Structure is [documented below](#nested_lifecycle_rules_rule).


<a name="nested_lifecycle_rules_rule"></a>The `rule` block supports:

* `action` -
  (Required)
  The action to take.
  Structure is [documented below](#nested_lifecycle_rules_rule_rule_action).

* `condition` -
  (Required)
  The condition(s) under which the action will be taken.
  Structure is [documented below](#nested_lifecycle_rules_rule_rule_condition).


<a name="nested_lifecycle_rules_rule_rule_action"></a>The `action` block supports:

* `storage_class` -
  (Optional)
  Target storage class. Required if the type of the
  action is SetStorageClass.

* `type` -
  (Required)
  Type of the action. Currently, only Delete and
  SetStorageClass are supported.
  Possible values are: `Delete`, `SetStorageClass`.

<a name="nested_lifecycle_rules_rule_rule_condition"></a>The `condition` block supports:

* `age` -
  (Optional)
  Age of an object (in days). This condition is satisfied
  when an object reaches the specified age.

* `created_before` -
  (Optional)
  A date in RFC 3339 format with only the date part (for
  instance, "2013-01-15"). This condition is satisfied
  when an object is created before midnight of the
  specified date in UTC.

* `custom_time_before` -
  (Optional)
  A date in the RFC 3339 format YYYY-MM-DD. This condition
  is satisfied when the customTime metadata for the object
  is set to an earlier date than the date used in
  this lifecycle condition.

* `days_since_custom_time` -
  (Optional)
  Days since the date set in the customTime metadata for the
  object. This condition is satisfied when the current date
  and time is at least the specified number of days after
  the customTime.

* `days_since_noncurrent_time` -
  (Optional)
  Relevant only for versioned objects. This condition is
  satisfied when an object has been noncurrent for more than
  the specified number of days.

* `with_state` -
  (Optional)
  Match to live and/or archived objects. 
  Unversioned buckets have only live objects. 
  Supported values include: "LIVE", "ARCHIVED", "ANY"..

* `matches_storage_class` -
  (Optional)
  Objects having any of the storage classes specified by
  this condition will be matched. Values include
  MULTI_REGIONAL, REGIONAL, NEARLINE, COLDLINE, ARCHIVE,
  STANDARD, and DURABLE_REDUCED_AVAILABILITY.

* `matches_suffix` -
  (Optional)
  The suffix of an object. This condition is
  satisfied when the end of an object's
  name is an exact case-sensitive match with the suffix.

* `matches_prefix` -
  (Optional)
  The prefix of an object. This condition
  is satisfied when the beginning of an object's name
  is an exact case-sensitive match with the prefix.

* `noncurrent_time_before` -
  (Optional)
  Relevant only for versioned objects. A date in the
  RFC 3339 format YYYY-MM-DD. This condition is satisfied
  for objects that became noncurrent on a date prior to the
  one specified in this condition.

* `num_newer_versions` -
  (Optional)
  Relevant only for versioned objects. If the value is N,
  this condition is satisfied when there are at least N
  versions (including the live version) newer than this
  version of the object.

* `send_age_if_zero` -
  (Optional)
  While set true, age value will be sent in the request 
  even for zero value of the field. This field is only useful 
  for setting 0 value to the age field. 
  It can be used alone or together with age.

* `send_days_since_noncurrent_time_if_zero` -
  (Optional)
  While set true, days_since_noncurrent_time value will be sent 
  in the request even for zero value of the field. 
  This field is only useful for setting 0 value to 
  the days_since_noncurrent_time field. 
  It can be used alone or together with days_since_noncurrent_time.

* `send_days_since_custom_time_if_zero` -
  (Optional)
  While set true, days_since_custom_time value will be sent 
  in the request even for zero value of the field. 
  This field is only useful for setting 0 value to 
  the days_since_custom_time field. 
  It can be used alone or together with days_since_custom_time.

* `send_num_newer_versions_if_zero` -
  (Optional)
  While set true, num_newer_versions value will be sent 
  in the request even for zero value of the field. 
  This field is only useful for setting 0 value to 
  the num_newer_versions field. 
  It can be used alone or together with num_newer_versions.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{bucket}}`

* `create_time` -
  The timestamp at which this folder was created.

* `update_time` -
  The timestamp at which this folder was most recently updated.
* `self_link` - The URI of the created resource.


## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


BucketLifeCycleConfig can be imported using any of these accepted formats:

* `{{name}}`


In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import BucketLifeCycleConfig using one of the formats above. For example:

```tf
import {
  id = "{{name}}"
  to = google_storage_bucket_life_cycle_config.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), BucketLifeCycleConfig can be imported using one of the formats above. For example:

```
$ terraform import google_storage_bucket_life_cycle_config.default {{name}}
```
