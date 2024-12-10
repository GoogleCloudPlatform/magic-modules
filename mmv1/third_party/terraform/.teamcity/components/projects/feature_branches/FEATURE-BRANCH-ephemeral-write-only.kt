/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

package projects.feature_branches

import ProviderNameBeta
import ProviderNameGa
import SharedResourceNameBeta
import SharedResourceNameGa
import SharedResourceNameVcr
import builds.*
import generated.ServicesListBeta
import generated.ServicesListGa
import jetbrains.buildServer.configs.kotlin.Project
import replaceCharsId
import vcs_roots.HashiCorpVCSRootBeta
import vcs_roots.HashiCorpVCSRootGa
import vcs_roots.ModularMagicianVCSRootBeta
import vcs_roots.ModularMagicianVCSRootGa

const val featureBranchEphemeralWriteOnly = "FEATURE-BRANCH-ephemeral-write-only"
const val EphemeralWriteOnlyTfCoreVersion = "1.10.0" // will be changed to 1.11.0 when the new ephemeral values feature is released in release candidates

// featureBranchEphemeralWriteOnlySubProject creates a project just for testing ephemeral write-only attributes.
// We know that all ephemeral write-only attributes we're adding are part of the Resource Manager service, so we only include those builds.
// We create builds for testing the resourcemanager service:
//    - Against the GA hashicorp repo
//    - Against the GA modular-magician repo
//    - Against the Beta hashicorp repo
//    - Against the Beta modular-magician repo
// These resemble existing projects present in TeamCity, but these all use a more recent version of Terraform including
// the new ephemeral values feature.
fun featureBranchEphemeralWriteOnlySubProject(allConfig: AllContextParameters): Project {

    val projectId = replaceCharsId(featureBranchEphemeralWriteOnly)

    val vcrConfig = getVcrAcceptanceTestConfig(allConfig) // Reused below for both MM testing build configs
    val trigger  = NightlyTriggerConfiguration(
        branch = "refs/heads/$featureBranchEphemeralWriteOnly" // triggered builds must test the feature branch
    )

    // All Ephemeral Write-Only attributes are in the following packages
    var PackagesListWriteOnly = mapOf(
        "compute" to mapOf(
            "name" to "compute",
            "displayName" to "Compute",
            "path" to "./google/services/compute"
        ),
        "secretmanager" to mapOf(
            "name" to "secretmanager",
            "displayName" to "Secretmanager",
            "path" to "./google/services/secretmanager"
        ),
        "bigquerydatatransfer" to mapOf(
            "name" to "bigquerydatatransfer",
            "displayName" to "Bigquerydatatransfer",
            "path" to "./google/services/bigquerydatatransfer"
        ),
        "sql" to mapOf(
            "name" to "sql",
            "displayName" to "Sql",
            "path" to "./google/services/sql"
        ),
    )

    // GA
    var parentId = "${projectId}_HC_GA"
    val buildConfigHashiCorpGa = BuildConfigurationsForPackages(PackagesListWriteOnly, ProviderNameGa, parentId, vcsRoot, listOf(SharedResourceNameGa), config)
    packageBuildConfigs.forEach { buildConfiguration ->
        buildConfiguration.addTrigger(cron)
    }

    var parentId = "${projectId}_MM_GA"
    val buildConfigModularMagicianGa = BuildConfigurationsForPackages(PackagesListWriteOnly, ProviderNameGa, parentId, vcsRoot, listOf(SharedResourceNameGa), config)
    // No trigger added here (MM upstream is manual only)

    // Beta
    parentId = "${projectId}_HC_Beta"
    val buildConfigHashiCorpBeta = BuildConfigurationsForPackages(PackagesListWriteOnly, ProviderNameBeta, parentId, vcsRoot, listOf(SharedResourceNameBeta), config)
    buildConfigHashiCorpBeta.forEach { buildConfiguration ->
        buildConfiguration.addTrigger(cron)
    }

    parentId = "${projectId}_MM_Beta"
    val buildConfigModularMagicianBeta = BuildConfigurationsForPackages(PackagesListWriteOnly, ProviderNameBeta, parentId, vcsRoot, listOf(SharedResourceNameBeta), config)
    // No trigger added here (MM upstream is manual only)

    // Create build config for sweeping the ephemeral write-only project
    var sweepersList: Map<String,Map<String,String>>
    when(providerName) {
        ProviderNameGa -> sweepersList = SweepersListGa
        ProviderNameBeta -> sweepersList = SweepersListBeta
        else -> throw Exception("Provider name not supplied when generating a nightly test subproject")
    }
    val serviceSweeperConfig = BuildConfigurationForServiceSweeper(providerName, ServiceSweeperName, sweepersList, projectId, vcsRoot, sharedResources, config)
    val sweeperCron = cron.clone()
    sweeperCron.startHour += 5  // Ensure triggered after the package test builds are triggered
    serviceSweeperConfig.addTrigger(sweeperCron)

    // ------

    // Make all builds use a 1.11.0 version of TF core
    val allBuildConfigs = listOf(buildConfigHashiCorpGa, buildConfigModularMagicianGa, buildConfigHashiCorpBeta, buildConfigModularMagicianBeta)
    allBuildConfigs.forEach{ b ->
        b.overrideTerraformCoreVersion(EphemeralWriteOnlyTfCoreVersion)
    }

    // ------

    return Project{
        id(projectId)
        name = featureBranchEphemeralWriteOnly
        description = "Subproject for testing feature branch $featureBranchEphemeralWriteOnly"

        // Register build configs in the project
        packageBuildConfigs.forEach { buildConfiguration ->
            buildType(buildConfiguration)
        }
        buildType(serviceSweeperConfig)

        params {
            readOnlySettings()
        }
    }
}
