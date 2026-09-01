/*
 * Copyright IBM Corp. 2014, 2026
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects

import GlobalSweepersProjectName
import NightlyTestsProjectId
import SharedResourceNameBeta
import SharedResourceNameGa
import SharedResourceNameVcr
import builds.*
import generated.SweepersListGa
import jetbrains.buildServer.configs.kotlin.AbsoluteId
import jetbrains.buildServer.configs.kotlin.FailureAction
import jetbrains.buildServer.configs.kotlin.Project
import replaceCharsId
import vcs_roots.HashiCorpVCSRootGa

// globalSweepersSubProject returns a subproject that contains sweepers for global resources (projects, folders)
// Sweeping projects is an edge case because it doesn't respect boundaries between different testing projects GA/Beta/PR
fun globalSweepersSubProject(allConfig: AllContextParameters): Project {

    val sweeperId = replaceCharsId("GLOBAL_SWEEPER")

    // Get config for using the GA identity (arbitrary choice as sweeper isn't confined by GA/Beta etc.)
    val gaConfig = getGaAcceptanceTestConfig(allConfig)

    // List of ALL shared resources; avoid clashing with any other running build
    val sharedResources: List<String> = listOf(SharedResourceNameGa, SharedResourceNameBeta, SharedResourceNameVcr)

    // Compute IDs of the "All Tests" composite builds in GA and Beta nightly test projects
    // These IDs must mirror how googleSubProjectGa/Beta and nightlyTests() compute their project IDs
    val gaProjectId = replaceCharsId("GOOGLE")
    val betaProjectId = replaceCharsId("GOOGLE_BETA")
    val gaAllTestsId = AbsoluteId(replaceCharsId("${gaProjectId}_${NightlyTestsProjectId}_ALL_TESTS"))
    val betaAllTestsId = AbsoluteId(replaceCharsId("${betaProjectId}_${NightlyTestsProjectId}_ALL_TESTS"))

    // Create build config for sweeping project resources
    // Uses the HashiCorpVCSRootGa VCS Root so that the latest sweepers in hashicorp/terraform-provider-google are used
    val serviceSweeperConfig = BuildConfigurationForGlobalSweeper("N/A", "Project Sweeper", "GoogleProject", SweepersListGa, sweeperId, HashiCorpVCSRootGa, sharedResources, gaConfig)
    serviceSweeperConfig.addTrigger(NightlyTriggerConfiguration(startHour=12))
    serviceSweeperConfig.dependencies {
        snapshot(gaAllTestsId) {
            onDependencyFailure = FailureAction.IGNORE
            onDependencyCancel = FailureAction.IGNORE
        }
        snapshot(betaAllTestsId) {
            onDependencyFailure = FailureAction.IGNORE
            onDependencyCancel = FailureAction.IGNORE
        }
    }

    // Create build config for sweeping folder resources
    val folderSweeperConfig = BuildConfigurationForGlobalSweeper("N/A", "Folder Sweeper", "GoogleFolder", SweepersListGa, sweeperId, HashiCorpVCSRootGa, sharedResources, gaConfig)
    folderSweeperConfig.addTrigger(NightlyTriggerConfiguration(startHour=12))
    folderSweeperConfig.dependencies {
        snapshot(gaAllTestsId) {
            onDependencyFailure = FailureAction.IGNORE
            onDependencyCancel = FailureAction.IGNORE
        }
        snapshot(betaAllTestsId) {
            onDependencyFailure = FailureAction.IGNORE
            onDependencyCancel = FailureAction.IGNORE
        }
    }

    return Project{
        id(sweeperId)
        name = GlobalSweepersProjectName
        description = "Subproject containing build configurations for sweeping global resources like projects and folders"

        // Register build configs in the project
        buildType(serviceSweeperConfig)
        buildType(folderSweeperConfig)

        params {
            readOnlySettings()
        }
    }
}