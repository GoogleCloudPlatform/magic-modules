package projects.reused

import DefaultBuildTimeoutDuration
import DefaultParallelism
import SharedResourceNameVcr
import VcrRecordingProjectId
import builds.*
import jetbrains.buildServer.configs.kotlin.BuildType
import jetbrains.buildServer.configs.kotlin.Project
import jetbrains.buildServer.configs.kotlin.sharedResources
import jetbrains.buildServer.configs.kotlin.vcs.GitVcsRoot
import replaceCharsId

fun vcrRecording(parentProject:String, providerName: String, hashicorpVcsRoot: GitVcsRoot, modularMagicianVcsRoot: GitVcsRoot, config: AccTestConfiguration): Project {

    // Create unique ID for the dynamically-created project
    var projectId = "${parentProject}_${VcrRecordingProjectId}"
    projectId = replaceCharsId(projectId)

    var buildIdHashiCorp = replaceCharsId("${providerName}_HASHICORP_VCR")
    var buildIdModularMagician = replaceCharsId("${providerName}_MODMAGICIAN_VCR")

    // Shared resource allows VCR recording process to not class with acceptance test or sweeper
    var sharedResources: List<String> = listOf(SharedResourceNameVcr)


    val testPrefix = "TestAcc"
    val testTimeout = "12"
    val parallelism = DefaultParallelism
    val buildTimeout = DefaultBuildTimeoutDuration

    val path = "./${providerName}" // Path is just ./google(-beta) here, whereas other builds use a ./google/something/specific path

    return Project {
        id(projectId)
        name = "VCR Recording"
       description = "A project connected to the hashicorp/terraform-provider-${providerName} repository, where users can trigger ad-hoc tests to re-record VCR cassettes"

        buildType(
            // TODO - pull this into a function and re-use to make both VCR build configs in this project
            BuildType {
                id(buildIdHashiCorp)

                name = "VCR Recording - Using hashicorp/terraform-provider-${providerName}"

                vcs {
                    root(hashicorpVcsRoot)
                    cleanCheckout = true
                }

                steps {
                    checkVcrEnvironmentVariables()
                    setGitCommitBuildId()
                    tagBuildToIndicatePurpose()
                    configureGoEnv()
                    downloadTerraformBinary()
                    runVcrTestRecordingSetup()
                    runVcrAcceptanceTests()
                    runVcrTestRecordingSaveCassettes()
                }

                features {
                    golang()
                    if (sharedResources.isNotEmpty()) {
                        sharedResources {
                            // When the build runs, it locks the value(s) below
                            sharedResources.forEach { sr ->
                                lockAllValues(sr)
                            }
                        }
                    }
                }

                params {
                    configureGoogleSpecificTestParameters(config)
                    vcrEnvironmentVariables(config, providerName)
                    acceptanceTestBuildParams(parallelism, testPrefix, testTimeout)
                    terraformLoggingParameters(providerName)
                    terraformCoreBinaryTesting()
                    terraformShouldPanicForSchemaErrors()
                    readOnlySettings()
                    workingDirectory(path)
                }

                artifactRules = "%teamcity.build.checkoutDir%/debug*.txt"

                failureConditions {
                    errorMessage = true
                    executionTimeoutMin = buildTimeout
                }

            }
        )
        params{
            configureGoogleSpecificTestParameters(config)
        }
    }
}