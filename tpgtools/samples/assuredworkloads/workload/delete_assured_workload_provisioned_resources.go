// deleteAssuredWorkloadProvisionedResources deletes the resources provisioned by
// assured workloads.. this is needed in order to delete the parent resource
func deleteAssuredWorkloadProvisionedResources(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		timeout := *schema.DefaultTimeout(4 * time.Minute)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_assured_workloads_workload" {
				continue
			}
			resourceAttributes := rs.Primary.Attributes
			n, err := strconv.Atoi(resourceAttributes["resources.#"])
			log.Printf("[DEBUG]found %v resources\n", n)
			log.Println(resourceAttributes)
			if err != nil {
				return err
			}

			// first delete the projects
			for i := 0; i < n; i++ {
				typee := resourceAttributes[fmt.Sprintf("resources.%d.resource_type", i)]
				if !strings.Contains(typee, "PROJECT") {
					continue
				}
				resource_id := resourceAttributes[fmt.Sprintf("resources.%d.resource_id", i)]
				log.Printf("[DEBUG] searching for project %s\n", resource_id)
				err := retryTimeDuration(func() (reqErr error) {
					_, reqErr = config.NewResourceManagerClient(config.userAgent).Projects.Get(resource_id).Do()
					return reqErr
				}, timeout)
				if err != nil {
					log.Printf("[DEBUG] did not find project %sn", resource_id)
					continue
				}
				log.Printf("[DEBUG] found project %s\n", resource_id)

				err = retryTimeDuration(func() error {
					_, delErr := config.NewResourceManagerClient(config.userAgent).Projects.Delete(resource_id).Do()
					return delErr
				}, timeout)
				if err != nil {
					log.Printf("Error deleting project '%s': %s\n ", resource_id, err)
					continue
				}
				log.Printf("[DEBUG] deleted project %s\n", resource_id)
			}

			// Then delete the folders
			for i := 0; i < n; i++ {
				typee := resourceAttributes[fmt.Sprintf("resources.%d.resource_type", i)]
				if typee != "CONSUMER_FOLDER" {
					continue
				}
				resource_id := "folders/" + resourceAttributes[fmt.Sprintf("resources.%d.resource_id", i)]
				err := retryTimeDuration(func() error {
					var reqErr error
					_, reqErr = config.NewResourceManagerV2Client(config.userAgent).Folders.Get(resource_id).Do()
					return reqErr
				}, timeout)
				log.Printf("[DEBUG] searching for folder %s\n", resource_id)
				if err != nil {
					log.Printf("[DEBUG] did not find folder %sn", resource_id)
					continue
				}
				log.Printf("[DEBUG] found folder %s\n", resource_id)
				err = retryTimeDuration(func() error {
					_, reqErr := config.NewResourceManagerV2Client(config.userAgent).Folders.Delete(resource_id).Do()
					return reqErr
				}, timeout)
				if err != nil {
					return fmt.Errorf("Error deleting folder '%s': %s\n ", resource_id, err)
				}
				log.Printf("[DEBUG] deleted folder %s\n", resource_id)
			}
		}
		return nil
	}
}