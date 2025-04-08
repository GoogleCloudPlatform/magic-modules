// Prevent ForceNew when upgrading replication.automatic -> replication.auto
func secretManagerSecretAutoCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	oAutomatic, nAutomatic := diff.GetChange("replication.0.automatic")
	_, nAuto := diff.GetChange("replication.0.auto")
	autoLen := len(nAuto.([]interface{}))

	// Do not ForceNew if we are removing "automatic" while adding "auto"
	if oAutomatic == true && nAutomatic == false && autoLen > 0 {
		return nil
	}

	if diff.HasChange("replication.0.automatic") {
		if err := diff.ForceNew("replication.0.automatic"); err != nil {
			return err
		}
	}

	if diff.HasChange("replication.0.auto") {
		if err := diff.ForceNew("replication.0.auto"); err != nil {
			return err
		}
	}

	return nil
}
