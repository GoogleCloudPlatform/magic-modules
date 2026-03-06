"""
Resource-related custom build rules for Magic Modules.
"""

load("//build_defs:providers.bzl", "ProductInfo", "ResourceInfo", "TpgResourceInfo")

def _mm_resource_impl(ctx):
    return [
        ctx.attr.product[ProductInfo],
        ResourceInfo(
            name = ctx.label.name,
            yaml = ctx.file.src,
            has_sweeper = ctx.attr.has_sweeper,
            has_operation = ctx.attr.has_operation,
        ),
        DefaultInfo(files = depset([ctx.file.src])),
    ]

mm_resource = rule(
    implementation = _mm_resource_impl,
    attrs = {
        "src": attr.label(
            allow_single_file = [".yaml"],
            mandatory = True,
        ),
        "product": attr.label(
            providers = [ProductInfo],
            mandatory = True,
        ),
        "has_sweeper": attr.bool(default = False, mandatory = False),
        "has_operation": attr.bool(default = False, mandatory = False),
    },
)

def _tpg_resource_impl(ctx):
    product = ctx.attr.resource[ProductInfo]
    resource = ctx.attr.resource[ResourceInfo]
    resource_name = "resource_{}_{}".format(product.name, resource.name)
    inputs = [product.yaml, resource.yaml] + [f for f in ctx.files._templates]
    outputs = [
        ctx.actions.declare_file("{}/{}.go".format(product.version, resource_name)),
        ctx.actions.declare_file("{}/{}_generated_meta.yaml".format(product.version, resource_name)),
    ]
    ctx.actions.run(
        executable = ctx.executable._compiler,
        arguments = [
            "--product",
            product.yaml.path,
            "--resource",
            resource.yaml.path,
            "--version",
            product.version,
            "--product_name",
            product.name,
            "--type",
            "resource",
            "--provider",
            "tpg",
            "--output",
            outputs[0].path,
        ],
        inputs = depset([i for i in inputs]),
        outputs = [outputs[0]],
        mnemonic = "TpgGenerateResource",
    )
    ctx.actions.run(
        executable = ctx.executable._compiler,
        arguments = [
            "--product",
            product.yaml.path,
            "--resource",
            resource.yaml.path,
            "--version",
            product.version,
            "--product_name",
            product.name,
            "--type",
            "metadata",
            "--provider",
            "tpg",
            "--output",
            outputs[1].path,
        ],
        inputs = depset([i for i in inputs]),
        outputs = [outputs[1]],
        mnemonic = "TpgGenerateResourceMetadata",
    )
    if resource.has_sweeper:
        sweeper_go = ctx.actions.declare_file("{}/{}_sweeper.go".format(product.version, resource_name))
        ctx.actions.run(
            executable = ctx.executable._compiler,
            arguments = [
                "--product",
                product.yaml.path,
                "--resource",
                resource.yaml.path,
                "--version",
                product.version,
                "--product_name",
                product.name,
                "--type",
                "sweeper",
                "--provider",
                "tpg",
                "--output",
                sweeper_go.path,
            ],
            inputs = depset([i for i in inputs]),
            outputs = [sweeper_go],
            mnemonic = "TpgGenerateResourceSweeper",
        )
        outputs.append(sweeper_go)
    return [
        ctx.attr.resource[ProductInfo],
        ctx.attr.resource[ResourceInfo],
        TpgResourceInfo(
            src = outputs[0],
            metadata = outputs[1],
        ),
        DefaultInfo(files = depset([out for out in outputs])),
    ]

tpg_resource = rule(
    implementation = _tpg_resource_impl,
    attrs = {
        "resource": attr.label(
            providers = [ResourceInfo, ProductInfo],
            allow_single_file = [".yaml"],
            mandatory = True,
        ),
        "version": attr.string(default = "ga", mandatory = False),
        "_compiler": attr.label(
            default = Label("//mmv1/cmd"),
            allow_single_file = True,
            executable = True,
            cfg = "exec",
        ),
        "_templates": attr.label(
            default = Label("//mmv1/templates"),
        ),
    },
)
