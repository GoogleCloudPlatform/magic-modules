The docsite is generated using [Hugo](https://gohugo.io/) and hosted using Github Pages. The theme is [Hugo Book](https://themes.gohugo.io/themes/hugo-book/) by [Alex Shpak](https://github.com/alex-shpak/). Magic Modules documentation should adhere to the [Google developer documentation style guide](https://developers.google.com/style/).

To view locally:

1. Install Hugo v0.136.5
   ```bash
   CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@v0.136.5
   ```
2. Run `hugo server` inside the `docs` directory
3. Visit http://localhost:1313/magic-modules/ to view the docs


If you are having deployment issues, try to reset your hugo module cache.
* `hugo mod clean`

To upgrade the theme version:
1. find the version you want at https://github.com/alex-shpak/hugo-book/commits/master
2. Run the following
```bash
go get github.com/alex-shpak/hugo-book@{{commit_hash}}
## example
## go get github.com/alex-shpak/hugo-book@d86d5e70c7c0d787675b13d9aee749c1a8b34776
```
