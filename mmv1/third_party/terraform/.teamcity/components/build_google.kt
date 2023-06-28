import jetbrains.buildServer.configs.kotlin.ParametrizedWithType

class ClientConfiguration(var custId: String,
                          var org: String,
                          val org2 : String,
                          val billingAccount : String,
                          val billingAccount2: String,
                          val masterBillingAccount: String,
                          val credentials : String,) {
}

fun ParametrizedWithType.ConfigureGoogleSpecificTestParameters(environment: String, config: ClientConfiguration) {
    hiddenVariable("env.GOOGLE_CUST_ID", config.custId, "The ID of the Google Identity Customer")
    hiddenVariable("env.GOOGLE_ORG", config.org, "The Google Organization Id")
    hiddenVariable("env.GOOGLE_ORG_2", config.org2, "The second Google Organization Id")
    hiddenVariable("env.GOOGLE_BILLING_ACCOUNT", config.billingAccount, "The billing account associated with the first google organization")
    hiddenVariable("env.GOOGLE_BILLING_ACCOUNT_2", config.billingAccount2, "The billing account associated with the second google organization")
    hiddenVariable("env.GOOGLE_MASTER_BILLING_ACCOUNT", config.masterBillingAccount, "The master billing account")
    hiddenPasswordVariable("env.GOOGLE_CREDENTIALS", config.credentials, "The Google credentials for this test runner")
}