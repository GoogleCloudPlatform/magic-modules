/*
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
import generated.ServicesListGa
import generated.ServicesListBeta

// This file is maintained in the GoogleCloudPlatform/magic-modules repository and copied into the downstream provider repositories. Any changes to this file in the downstream will be overwritten.

fun getServicesList(Services: Array<String>, version: String): Map<String,Map<String,String>> {
    if (Services.isEmpty()) {
        throw Exception("No services found for version $version")
    }

    var servicesList = mutableMapOf<String,Map<String,String>>()
    for (service in Services) {
        if (version == "GA" || version == "GA-MM") {
            servicesList[service] = ServicesListGa.getOrElse(service) { throw Exception("Service $service not found") }
        } else if (version == "Beta" || version == "Beta-MM") {
            servicesList[service] = ServicesListBeta.getOrElse(service) { throw Exception("Service $service not found") }
        } else {
            throw Exception("Invalid version $version")
        }
    }

    when (version) {
        "GA" -> return servicesList
        "Beta" -> {
            servicesList = servicesList.mapValues { (_, value) ->
                value + mapOf(
                    "displayName" to "${value["displayName"]} - Beta",
                    "path" to (value["path"]?.replace("./google/", "./google-beta/") ?: "")
                )
            }.toMutableMap()
        }
        "GA-MM" -> {
            servicesList = servicesList.mapValues { (_, value) ->
                value + mapOf(
                    "displayName" to "${value["displayName"]} - MM"
                )
            }.toMutableMap()
        }
        "Beta-MM" -> {
            servicesList = servicesList.mapValues { (_, value) ->
                value + mapOf(
                    "displayName" to "${value["displayName"]} - Beta - MM",
                    "path" to (value["path"]?.replace("./google-beta/", "./google-beta/services/") ?: "")
                )
            }.toMutableMap()
        }
        else -> throw Exception("Invalid version $version")
    }
    
    return servicesList
}