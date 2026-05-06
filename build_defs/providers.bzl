"""
Provider definitions for Magic Modules custom build rules.
"""

ProductInfo = provider(
    "Provider for Magic Modules products.",
    fields = {
        "name": "name of the product",
        "version": "version of the provider to be generated",
        "yaml": "the product.yaml file",
    },
)

ResourceInfo = provider(
    "Provider for Magic Modules resources.",
    fields = {
        "name": "name of the resource",
        "yaml": "the [resource].yaml input file",
        "has_sweeper": "whether a sweeper must be generated",
        "has_operation": "whether an operation must be generated",
    },
)

TpgResourceInfo = provider(
    "Provider for TPG resources.",
    fields = {
        "src": "the generated go file",
        "metadata": "the generated metadata.yaml file",
        "sweeper": "the optional generated sweeper file",
    },
)
