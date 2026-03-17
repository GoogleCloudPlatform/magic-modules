package tpgdclresource

// RequiresRecreate is for Operations that require recreating.
func RequiresRecreate() func(d *FieldDiff) []string {
	return func(d *FieldDiff) []string { return []string{"Recreate"} }
}

// TriggersOperation is used to tell the diff checker to trigger an operation.
func TriggersOperation(op string) func(d *FieldDiff) []string {
	return func(d *FieldDiff) []string { return []string{op} }
}
