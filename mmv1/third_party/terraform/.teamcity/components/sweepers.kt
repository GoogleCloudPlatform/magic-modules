/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.AbsoluteId

class sweeperBuildConfigs(packageName: String) {
    val packageName = packageName

    fun preSweeperBuildConfig(path: String, manualVcsRoot: AbsoluteId, parallelism: Int, triggerConfig: NightlyTriggerConfiguration, environmentVariables: ClientConfiguration) : BuildType {
        val testPrefix = "TestAcc"
        val testTimeout = "12"
        val sweeperRegions = "us-central1"
        val sweeperRun = "" // Empty string means all sweepers run

        val configName = "Pre-Sweeper"
        val sweeperStepName = "Pre-Sweeper"

        return createBuildConfig(manualVcsRoot, preSweeperBuildConfigId, configName, sweeperStepName, parallelism, testPrefix, testTimeout, sweeperRegions, sweeperRun, path, packageName, triggerConfig, environmentVariables)
   }

    fun postSweeperBuildConfig(path: String, manualVcsRoot: AbsoluteId, parallelism: Int, triggerConfig: NightlyTriggerConfiguration, environmentVariables: ClientConfiguration, dependencies: ArrayList<String>) : BuildType {
        val testPrefix = "TestAcc"
        val testTimeout = "12"
        val sweeperRegions = "us-central1"
        val sweeperRun = "" // Empty string means all sweepers run

        val configName = "Post-Sweeper"
        val sweeperStepName = "Post-Sweeper"

        val build = createBuildConfig(manualVcsRoot, postSweeperBuildConfigId, configName, sweeperStepName, parallelism, testPrefix, testTimeout, sweeperRegions, sweeperRun, path, packageName, triggerConfig, environmentVariables)
        build.addDependencies(dependencies)

        return build
    }

    fun createBuildConfig(
        manualVcsRoot: AbsoluteId,
        configId: String,
        configName: String,
        sweeperStepName: String,
        parallelism: Int,
        testPrefix: String,
        testTimeout: String,
        sweeperRegions: String,
        sweeperRun: String,
        path: String,
        packageName: String,
        triggerConfig: NightlyTriggerConfiguration,
        environmentVariables: ClientConfiguration,
        buildTimeout: Int = defaultBuildTimeoutDuration
        ) : BuildType {
        return BuildType {

            id(configId)

            name = configName

            vcs {
                root(rootId = manualVcsRoot)
                cleanCheckout = true
            }

            steps {
                SetGitCommitBuildId()
                ConfigureGoEnv()
                DownloadTerraformBinary()
                RunSweepers(sweeperStepName)
            }

            failureConditions {
                errorMessage = true
            }

            features {
                Golang()
            }

            params {
                ConfigureGoogleSpecificTestParameters(environmentVariables)
                TerraformAcceptanceTestParameters(parallelism, testPrefix, testTimeout, sweeperRegions, sweeperRun)
                TerraformAcceptanceTestsFlag()
                TerraformCoreBinaryTesting()
                TerraformShouldPanicForSchemaErrors()
                ReadOnlySettings()
                WorkingDirectory(path)
            }

            triggers {
                RunNightly(triggerConfig)
            }

            failureConditions {
                executionTimeoutMin = buildTimeout
            }
        }
    }
}
