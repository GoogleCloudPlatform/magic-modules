"""
Product-related custom build rules for Magic Modules.
"""

load("//build_defs:providers.bzl", "ProductInfo", "ResourceInfo", "TpgResourceInfo")

def _mm_product_impl(ctx):
    return [ProductInfo(
        name = ctx.label.name,
        yaml = ctx.file.src,
        version = ctx.attr.version,
    )]

mm_product = rule(
    implementation = _mm_product_impl,
    attrs = {
        "version": attr.string(default = "ga", mandatory = False),
        "src": attr.label(
            allow_single_file = [".yaml"],
            mandatory = True,
        ),
    },
)

def _tpg_product_impl(ctx):
    product = ctx.attr.product[ProductInfo]
    resources = [res[ResourceInfo] for res in ctx.attr.resources]
    tpg_resources = [res[TpgResourceInfo] for res in ctx.attr.resources]
    operations = [res for res in resources if res.has_operation]
    inputs = [product.yaml] + [res.metadata for res in tpg_resources] + [res.src for res in tpg_resources] + [res.yaml for res in resources] + [f for f in ctx.files._templates]

    outputs = [
        ctx.actions.declare_file("{}/product.go".format(product.version)),
    ]
    ctx.actions.run(
        executable = ctx.executable._compiler,
        arguments = [
            "--product",
            product.yaml.path,
            "--version",
            product.version,
            "--product_name",
            product.name,
            "--type",
            "product",
            "--provider",
            "tpg",
            "--output",
            outputs[0].path,
        ],
        inputs = depset([i for i in inputs]),
        outputs = outputs,
        mnemonic = "TpgGenerateProduct",
    )
    if operations:
        operation_go = ctx.actions.declare_file("{}/{}_operation.go".format(product.version, product.name))
        ctx.actions.run(
            executable = ctx.executable._compiler,
            arguments = [
                "--product",
                product.yaml.path,
                "--resource",
                operations[0].yaml.path,
                "--version",
                product.version,
                "--product_name",
                product.name,
                "--type",
                "operation",
                "--provider",
                "tpg",
                "--output",
                operation_go.path,
            ],
            inputs = depset([i for i in inputs]),
            outputs = [operation_go],
            mnemonic = "TpgGenerateProductOperation",
        )
        outputs.append(operation_go)

    return [
        ctx.attr.product[ProductInfo],
        DefaultInfo(files = depset([out for out in outputs])),
    ]

tpg_product = rule(
    implementation = _tpg_product_impl,
    attrs = {
        "product": attr.label(
            providers = [ProductInfo],
            mandatory = True,
        ),
        "resources": attr.label_list(
            providers = [ResourceInfo, TpgResourceInfo],
            mandatory = True,
        ),
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
