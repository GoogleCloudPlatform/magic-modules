/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// this file is copied from mmv1, any changes made here will be overwritten

import jetbrains.buildServer.configs.kotlin.*
import jetbrains.buildServer.configs.kotlin.AbsoluteId

class packageDetails(packageName: String, displayName: String, providerName: String, environment: String) {
    val packageName = packageName
    val displayName = displayName
    val providerName = providerName
    val environment = environment

    // buildConfiguration returns a BuildType for a service package
    // For BuildType docs, see https://teamcity.jetbrains.com/app/dsl-documentation/root/build-type/index.html
    fun buildConfiguration(path: String, manualVcsRoot: AbsoluteId, parallelism: Int, triggerConfig: NightlyTriggerConfiguration) : BuildType {
        return BuildType {
            // TC needs a consistent ID for dynamically generated packages
            id(uniqueID(providerName))

            name = "%s - Acceptance Tests".format(displayName)

            vcs {
                root(rootId = manualVcsRoot)
                cleanCheckout = true
            }

            steps {
                ConfigureGoEnv()
                DownloadTerraformBinary()
                RunAcceptanceTests()
            }

            failureConditions {
                errorMessage = true
            }

            features {
                Golang()
            }

            params {
                TerraformAcceptanceTestParameters(parallelism, "TestAcc", "12", "us-central1", "")
                TerraformAcceptanceTestsFlag()
                TerraformCoreBinaryTesting()
                TerraformShouldPanicForSchemaErrors()
                ReadOnlySettings()
                WorkingDirectory(path, packageName)
            }

            triggers {
                RunNightly(triggerConfig)
            }
        }
    }

    fun uniqueID(provider: String) : String {
        // Replacing chars can be necessary, due to limitations on IDs
        // "ID should start with a latin letter and contain only latin letters, digits and underscores (at most 225 characters)." 
        var pv = provider.replace("-", "").toUpperCase()
        var env = environment.toUpperCase().replace("-", "").replace(".", "").toUpperCase()
        var pkg = packageName.toUpperCase()

        return "%s_SERVICE_%s_%s".format(pv, env, pkg)
    }
}
