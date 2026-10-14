/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package generated

// specifies the default level of parallelism for per-service-package test builds, unless overridden in ServiceParallelism below
const val DefaultServiceParallelism = 12

var ServiceParallelism = mapOf(
    // Kept at the previous default of 6 while failures at that level are investigated before raising them to DefaultServiceParallelism.
    // alloydb: Error 409 "unable to queue the operation"
    "alloydb" to 6,
    // compute: Quota 'SECURITY_POLICIES' exceeded (limit 30 globally)
    "compute" to 6,
    // container: GKE cluster create/delete failures
    "container" to 6,
    // gkebackup: GKE cluster create/delete failures; shares GKE capacity with container
    "gkebackup" to 6,
    // netapp: BackupVaultsPerRegion and ReplicatedVolumesPerRegion quotas exceeded
    "netapp" to 6,
    // networkservices: operations failing with "deadline exceeded"
    "networkservices" to 6,
    // redis: tests conflict deleting shared ServiceConnectionPolicies
    "redis" to 6,
    // sql: "zone or region does not have sufficient resources"
    "sql" to 6,

    "looker" to 1
)
