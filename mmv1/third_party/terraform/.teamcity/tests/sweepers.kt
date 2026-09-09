/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package tests

import ServiceSweeperCronName
import ServiceSweeperManualName
import ServiceSweeperName
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Project
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import projects.googleCloudRootProject

class SweeperTests {
    @Test
    fun globalSweepersConfig() {
        val root = googleCloudRootProject(testContextParameters())

        // Find Global sweepers project
        val globalSweepersProject = getSubProject(root, globalSweepersProjectName)

        globalSweepersProject.buildTypes.forEach{bt ->
            val skipProjectSweeper = bt.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
            assertTrue("env.SKIP_PROJECT_SWEEPER should be set to an empty value in the ${globalSweepersProject.name} project. Value = `${skipProjectSweeper}` ", skipProjectSweeper == "")

            val skipFolderSweeper = bt.params.findRawParam("env.SKIP_FOLDER_SWEEPER")!!.value
            assertTrue("env.SKIP_FOLDER_SWEEPER should be set to an empty value in the ${globalSweepersProject.name} project. Value = `${skipFolderSweeper}` ", skipFolderSweeper == "")
        }
    }

    @Test
    fun gaNightlyTestsServiceSweeperConfig() {
        val root = googleCloudRootProject(testContextParameters())

        // Find GA nightly test project
        val project = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)

        // Find sweeper inside
        val sweeper = getBuildFromProject(project, ServiceSweeperName)

        // Check PACKAGE_PATH is in google (not google-beta)
        val value = sweeper.params.findRawParam("PACKAGE_PATH")!!.value
        assertEquals("./google/sweeper", value)

        // SKIP_PROJECT_SWEEPER and SKIP_FOLDER_SWEEPER should have values so they will be skipped
        val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
        assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")

        val skipFolderSweeper = sweeper.params.findRawParam("env.SKIP_FOLDER_SWEEPER")!!.value
        assertTrue("env.SKIP_FOLDER_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipFolderSweeper}` ", skipFolderSweeper != "")
    }

    @Test
    fun betaNightlyTestsServiceSweeperConfig() {
        val root = googleCloudRootProject(testContextParameters())

        // Find Beta nightly test project
        val project = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)

        // Find sweeper inside
        val sweeper: BuildType = getBuildFromProject(project, ServiceSweeperName)

        // Check PACKAGE_PATH is in google-beta
        val value = sweeper.params.findRawParam("PACKAGE_PATH")!!.value
        assertEquals("./google-beta/sweeper", value)

        // SKIP_PROJECT_SWEEPER and SKIP_FOLDER_SWEEPER should have values so they will be skipped
        val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
        assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")

        val skipFolderSweeper = sweeper.params.findRawParam("env.SKIP_FOLDER_SWEEPER")!!.value
        assertTrue("env.SKIP_FOLDER_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipFolderSweeper}` ", skipFolderSweeper != "")
    }

    @Test
    fun gaMmUpstreamServiceSweeperConfig() {
        val root = googleCloudRootProject(testContextParameters())

        // Find Beta nightly test project
        val project = getNestedProjectFromRoot(root, gaProjectName, mmUpstreamProjectName)

        // Find sweepers inside
        val cronSweeper = getBuildFromProject(project, ServiceSweeperCronName)
        val manualSweeper = getBuildFromProject(project, ServiceSweeperManualName)
        val allSweepers: ArrayList<BuildType> = arrayListOf(cronSweeper, manualSweeper)
        allSweepers.forEach{ sweeper ->
            // Check PACKAGE_PATH is in google-beta
            val value = sweeper.params.findRawParam("PACKAGE_PATH")!!.value
            assertEquals("./google/sweeper", value)

            // SKIP_PROJECT_SWEEPER and SKIP_FOLDER_SWEEPER should have values so they will be skipped
            val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
            assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")

            val skipFolderSweeper = sweeper.params.findRawParam("env.SKIP_FOLDER_SWEEPER")!!.value
            assertTrue("env.SKIP_FOLDER_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipFolderSweeper}` ", skipFolderSweeper != "")
        }
    }

    @Test
    fun betaMmUpstreamServiceSweeperConfig() {
        val root = googleCloudRootProject(testContextParameters())

        // Find Beta nightly test project
        val project = getNestedProjectFromRoot(root, betaProjectName, mmUpstreamProjectName)

        // Find sweepers inside
        val cronSweeper = getBuildFromProject(project, ServiceSweeperCronName)
        val manualSweeper = getBuildFromProject(project, ServiceSweeperManualName)
        val allSweepers: ArrayList<BuildType> = arrayListOf(cronSweeper, manualSweeper)
        allSweepers.forEach{ sweeper ->
            // Check PACKAGE_PATH is in google-beta
            val value = sweeper.params.findRawParam("PACKAGE_PATH")!!.value
            assertEquals("./google-beta/sweeper", value)

            // SKIP_PROJECT_SWEEPER and SKIP_FOLDER_SWEEPER should have values so they will be skipped
            val skipProjectSweeper = sweeper.params.findRawParam("env.SKIP_PROJECT_SWEEPER")!!.value
            assertTrue("env.SKIP_PROJECT_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipProjectSweeper}` ", skipProjectSweeper != "")

            val skipFolderSweeper = sweeper.params.findRawParam("env.SKIP_FOLDER_SWEEPER")!!.value
            assertTrue("env.SKIP_FOLDER_SWEEPER should be set to a non-empty string in the ${project.name} project (${sweeper.name}). Value = `${skipFolderSweeper}` ", skipFolderSweeper != "")
        }
    }

    @Test
    fun globalSweepersDependOnAllNightlyTests() {
        val root = googleCloudRootProject(testContextParameters())

        // Find GA nightly test project's service sweeper
        val gaNightlyTests: Project = getNestedProjectFromRoot(root, gaProjectName, nightlyTestsProjectName)
        val sweeperGa: BuildType = getBuildFromProject(gaNightlyTests, ServiceSweeperName)

        // Find Beta nightly test project's service sweeper
        val betaNightlyTests : Project = getNestedProjectFromRoot(root, betaProjectName, nightlyTestsProjectName)
        val sweeperBeta: BuildType = getBuildFromProject(betaNightlyTests, ServiceSweeperName)

        // Find Global sweepers project's builds
        val globalSweepersProject = getSubProject(root, globalSweepersProjectName)
        val projectSweeper: BuildType = getBuildFromProject(globalSweepersProject, "Project Sweeper")
        val folderSweeper: BuildType = getBuildFromProject(globalSweepersProject, "Folder Sweeper")

        // Each nightly service sweeper follows its composite, while global sweepers
        // follow the GA service sweeper after both GA and Beta sweepers complete.
        assertTrue(sweeperGa.triggers.items.size == 1)
        assertTrue(sweeperBeta.triggers.items.size == 1)
        assertTrue(projectSweeper.triggers.items.size == 1)
        assertTrue(folderSweeper.triggers.items.size == 1)

        assertEquals("Project sweeper should snapshot-depend on GA and Beta composites and service sweepers", 4, projectSweeper.dependencies.items.size)
        assertEquals("Folder sweeper should snapshot-depend on GA and Beta composites and service sweepers", 4, folderSweeper.dependencies.items.size)
    }
}
