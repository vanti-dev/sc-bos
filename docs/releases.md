# Creating Releases

Releases are created automatically by GitHub Actions when a tag is pushed to the repository. The tag should be of the
form `v0.YYYY.R` (e.g. `v0.2024.1`). The 0 prefix is used to indicate the project is still under development and 
private to Vanti - i.e. not yet publicly available. The `YYYY` is the year, and `R` is the release number for that year.

## Testing releases

You can test the release process locally using [act](https://github.com/nektos/act). This will run mostly the same steps 
as the GitHub Actions workflow, but locally without actually releasing/publishing to GitHub. You will need to have 
Docker installed, the `act` binary, and a file containing the secrets required for the build (the example below uses 
`.secret/act.secrets`):

```
GITHUB_TOKEN=<your github personal access token, should start `ghp_`>
GO_MOD_TOKEN=<as above>
NEXUS_NPM_TOKEN=<your nexus npm token - copy from ~/.npmrc, should start `NpmToken.`>
```

You can then run the following commands to perform the build of the binaries, UIs, and Docker images:
```shell
act --secret-file .secret/act.secrets --artifact-server-path ./.downloads/sc-bos --pull false -j build-sc-bos
act --secret-file .secret/act.secrets --artifact-server-path ./.downloads/ops-ui --pull false -j build-ops-ui
act --secret-file .secret/act.secrets --artifact-server-path ./.downloads --pull false -j build-docker
``` 

