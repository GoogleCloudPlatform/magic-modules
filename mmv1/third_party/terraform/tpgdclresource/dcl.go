package tpgdclresource

var (
	// CreateDirective restricts Apply to creating resources for Create
	CreateDirective = []ApplyOption{
		WithLifecycleParam(BlockAcquire),
		WithLifecycleParam(BlockDestruction),
		WithLifecycleParam(BlockModification),
	}

	// UpdateDirective restricts Apply to modifying resources for Update
	UpdateDirective = []ApplyOption{
		WithLifecycleParam(BlockCreation),
		WithLifecycleParam(BlockDestruction),
	}
)
