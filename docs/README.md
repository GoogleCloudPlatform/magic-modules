The docsite is generated using [Hugo](https://gohugo.io/) and hosted using Github Pages. The theme is [Hugo Book](https://themes.gohugo.io/themes/hugo-book/) by [Alex Shpak](https://github.com/alex-shpak/). Magic Modules documentation should adhere to the [Google developer documentation style guide](https://developers.google.com/style/).

To view locally:

1. Ensure you've installed [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git),
   [Go](https://go.dev/doc/install), and [Dart Sass](https://gohugo.io/hugo-pipes/transpile-sass-to-css/#dart-sass).
   You require these prerequisites for installing Hugo.

1. Install Hugo v0.136.5:
   ```bash
   CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@v0.136.5
   ```

1. Restart your terminal.

1. Clone the `magic-modules` GitHub repository:
   ```bash
   git clone https://github.com/GoogleCloudPlatform/magic-modules.git
   ```
1. Navigate to the `docs` directory inside the `magic-modules` repository:
   ```bash
   cd magic-modules/docs/
   ```

1. Start Hugo's development server to view the Magic Modules site:
   ```bash
   hugo server
   ```

1. View the docs by visiting the following URL in a browser window:
   ```bash
   http://localhost:1313/magic-modules/
   ```

You can press `Ctrl+C` to stop Hugo's development server.

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
