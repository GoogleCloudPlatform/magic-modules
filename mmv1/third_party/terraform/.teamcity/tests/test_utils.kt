/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is copied from mmv1, any changes made here will be overwritten

package tests

import jetbrains.buildServer.configs.kotlin.AbsoluteId

import ClientConfiguration

fun testConfiguration() : ClientConfiguration {
    return ClientConfiguration("custId", "org", "org2", "billingAccount", "billingAccount2", "masterBillingAccount", "credentials", "project", "orgDomain", "projectNumber", "region", "serviceAccount", "zone", "firestoreProject", "identityUser")
}

fun testVcsRootId() : AbsoluteId {
    return AbsoluteId("TerraformProviderFoobar")
}