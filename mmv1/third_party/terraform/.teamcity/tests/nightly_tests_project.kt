/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import AllNightlyTestsName
import ServiceSweeperName
import jetbrains.buildServer.configs.kotlin.BuildTypeSettings
import jetbrains.buildServer.configs.kotlin.triggers.ScheduleTrigger
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import projects.googleCloudRootProject

class NightlyTestProjectsTests {
    @Test
    fun allBuildsShouldHaveTrigger() {
        val root = googleCloudRootProject(testContextParameters())

        // Find GA nightly test project
        var gaNightlyTestProject = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)

        // Find Beta nightly test project
        var betaNightlyTestProject = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)

        // Package tests and the Service Sweeper keep their upstream CRON triggers.
        // The composite All Nightly Tests build is started via snapshot dependency only.
        (gaNightlyTestProject.buildTypes + betaNightlyTestProject.buildTypes).forEach{bt ->
            if (bt.name == AllNightlyTestsName) {
                assertTrue("Build configuration `${bt.name}` should not have a trigger; it is started via snapshot dependency", bt.triggers.items.isEmpty())
                return@forEach
            }

            assertTrue("Build configuration `${bt.name}` should contain at least one trigger", bt.triggers.items.isNotEmpty())
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

            assertTrue("Build configuration `${bt.name}` should contain a CRON/'schedulingTrigger' trigger", found)

            // Check that nightly test is being ran on the nightly-test branch
            var isNightlyTestBranch: Boolean = false
            if (schedulingTrigger.branchFilter == "+:refs/heads/nightly-test"){
                isNightlyTestBranch = true
            }
            assertTrue("Build configuration `${bt.name}` is using the nightly-test branch filter;", isNightlyTestBranch)
        }
    }

    @Test
    fun nightlyTestsShouldHaveCompositeAllTestsBuild() {
        val root = googleCloudRootProject(testContextParameters())

        var gaNightlyTestProject = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)
        var betaNightlyTestProject = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)

        listOf(gaNightlyTestProject, betaNightlyTestProject).forEach { project ->
            val composite = getBuildFromProject(project, AllNightlyTestsName)
            assertEquals("Build configuration `${composite.name}` should be a COMPOSITE build", BuildTypeSettings.Type.COMPOSITE, composite.type)

            val packageBuilds = project.buildTypes.filter { bt ->
                bt.name != ServiceSweeperName && bt.name != AllNightlyTestsName
            }
            assertTrue("Nightly test project `${project.name}` should have package test builds", packageBuilds.isNotEmpty())
            assertEquals(
                "Composite `${composite.name}` should snapshot-depend on every package test build",
                packageBuilds.size,
                composite.dependencies.items.size
            )

            val sweeper = getBuildFromProject(project, ServiceSweeperName)
            assertEquals("Service sweeper should snapshot-depend on the composite All Nightly Tests build", 1, sweeper.dependencies.items.size)
        }
    }
}
