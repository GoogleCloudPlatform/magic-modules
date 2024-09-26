/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import jetbrains.buildServer.configs.kotlin.triggers.ScheduleTrigger
import org.junit.Assert
import org.junit.Test
import projects.feature_branches.featureBranchEphemeralResources
import projects.googleCloudRootProject

class FeatureBranchEphemeralResourcesSubProject {
    @Test
    fun checkProjectSetup() {
        val root = googleCloudRootProject(testContextParameters())

        // Find feature branch project
        var project = getSubProject(root, featureBranchEphemeralResources)

        // Make assertions about builds in the feature branch testing project
        project.buildTypes.forEach{bt ->
            Assert.assertTrue(
                "Build configuration `${bt.name}` should contain at least one trigger",
                bt.triggers.items.isNotEmpty()
            )
            // Look for at least one CRON trigger
            var found: Boolean = false
            lateinit var schedulingTrigger: ScheduleTrigger
            for (item in bt.triggers.items){
                if (item.type == "schedulingTrigger") {
                    schedulingTrigger = item as ScheduleTrigger
                    found = true
                    break
                }
            }

            Assert.assertTrue(
                "Build configuration `${bt.name}` should contain a CRON/'schedulingTrigger' trigger",
                found
            )

            // Check that triggered builds are being run on the feature branch
            var isCorrectBranch: Boolean = schedulingTrigger.branchFilter == "+:refs/heads/${featureBranchEphemeralResources}"

            Assert.assertTrue(
                "Build configuration `${bt.name}` is using the ${featureBranchEphemeralResources} branch filter;",
                isCorrectBranch
            )
        }
    }
}
