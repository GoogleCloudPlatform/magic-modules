# Discovery Doc Generator
This script will generate a first-pass version of an api.yaml

## Usage
```
  ruby tools/builder/run.rb -p <name of the products/ directory> -u <discovery doc URL> -o <comma-separated list of objects to generate>
```
This will output a new api.yaml file at the root of the Magic Modules directory. If there's already an api.yaml that exists, this new api.yaml will be the same as the existing api.yaml, but with the new object / properties added.

## What this is not
* This will not generate perfect api.yamls! This is meant to help get through the boilerplate. You will absolutely need to change descriptions and add fields!

## FAQs

* **Why doesn't this overwrite the original api.yaml?** Machine-outputted YAML has different string formatting and whitespace choices. This will create an unintentionally large diff.
