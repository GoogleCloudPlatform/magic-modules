# Overriding value troubleshooting

1. Validate error

    ```
    NoMethodError: undefined method `validate'
    ```

    Solution: Don't forget `!ruby/object:Overrides::Terraform::PropertyOverride`

1. Not sending a value.

    ```
    ~ guest_attributes = {} -> (known after apply)
    - internal_ip_only = false -> null
    ~ labels           = {} -> (known after apply)
    ```

    Solution: Use `send_empty_value`

1.
