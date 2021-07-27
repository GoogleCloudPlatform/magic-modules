## doc parser

This tool assists in creating mmv1 yaml fields from cloud documentation.
It was designed with generation of subfields in mind but should work for any scenario.

To use this tool edit the path in index.mjs to point to the documentation path to
the resource and or field in question.
```
const host = "https://cloud.google.com"
const path = "/dlp/docs/reference/rest/v2/projects.deidentifyTemplates#DeidentifyTemplate.CryptoReplaceFfxFpeConfig"
```

Ensure `node > v16.0.0` is installed

```
npm install
node ./index.mjs
```

The generated yaml will be located at `./out.yaml`