/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import AllNightlyTestsName
import ServiceSweeperName
import jetbrains.buildServer.configs.kotlin.BuildTypeSettings
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

        // The composite starts on its nightly CRON trigger. Package builds have no
        // individual nightly triggers, and the Service Sweeper follows the composite.
        (gaNightlyTestProject.buildTypes + betaNightlyTestProject.buildTypes).forEach{bt ->
            if (bt.name == AllNightlyTestsName) {
                assertTrue("Build configuration `${bt.name}` should have a trigger", bt.triggers.items.isNotEmpty())
                return@forEach
            }

            if (bt.name == ServiceSweeperName) {
                assertEquals("Build configuration `${bt.name}` should have one finish trigger", 1, bt.triggers.items.size)
                return@forEach
            }

            if (bt.name != AllNightlyTestsName) {
                assertTrue("Package build configuration `${bt.name}` should not contain an individual nightly trigger", bt.triggers.items.isEmpty())
                return@forEach
            }
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
            assertEquals("Service sweeper should have one finish trigger from the composite", 1, sweeper.triggers.items.size)
        }
    }
}
