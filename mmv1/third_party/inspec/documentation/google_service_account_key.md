### Test that a GCP project IAM service account key has the expected key algorithm

    describe google_service_account_key(name: "projects/sample-project/serviceAccounts/test-sa@sample-project.iam.gserviceaccount.com/keys/c6bd986da9fac6d71178db41d1741cbe751a5080" ) do
      its('key_algorithm') { should eq "KEY_ALG_RSA_2048" }
    end