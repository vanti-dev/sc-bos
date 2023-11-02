# KeyCloak Theming

## Overview

This directory contains the theme files for the KeyCloak server.

## Sources

[KeyCloak theme structure](https://www.keycloak.org/docs/latest/server_development/index.html#_themes)
[Original KeyCloak Source code v21.0.1](https://github.com/keycloak/keycloak/releases/tag/21.0.1)
[How to customize Keycloak themes](https://trigodev.com/blog/how-to-customize-keycloak-themes)
[Custom Themes for KeyCloak](https://medium.com/front-end-weekly/custom-themes-for-keycloak-631bdd3e04e5)

## Theme Structure

> _This structure is based on the `../lib/lib/main/org.keycloak.keycloak-themes-21.0.1.jar` within the [Original KeyCloak Source code v21.0.1](https://github.com/keycloak/keycloak/releases/tag/21.0.1). You can handle it as a .zip, using an unarchiver app._
> _If you wish to create a new theme, it is better to copy either the `smartcore` or the `original` theme and rename it to the new theme name._

The theme is structured as follows:

```plaintext
themes
└── smartcore
    ├── account // Not included as not needed
    ├── admin // Not included as not needed
    ├── common
    ├── email // Not included as not needed
    ├── login
    └── welcome // Not included as not needed
```

The `smartcore` directory contains the theme files for the SmartCore theme. The `account`, `admin`, `common`, `email`, `login` and `welcome` directories contain the theme files for the respective KeyCloak pages.

__! _smartcore_ - DO NOT USE other than alphabetical character in the main folder name, otherwise the theme will not be loaded and/or other things will break.__

## Use of Theme Files

The `docker-compose.yml` file contains the following volumes:

```yaml
  volumes:
    - ./manifests/keycloak/themes/smartcore:/opt/keycloak/themes/smartcore
```

This means that the `smartcore` theme will be loaded (mounted) into the KeyCloak server (docker container) on startup.

If there is no docker available, the theme can be loaded into the KeyCloak server by copying the `smartcore` directory into the `themes` directory of the KeyCloak server.
The target `themes` directory is located in the `root/opt/keycloak/` directory of the KeyCloak server.

The end result (path) should look like this: `//root/opt/keycloak/themes/smartcore`.

The version of KeyCloak server know to be supporting the `smartcore` theme is `v21.0.1`. It hasn't been tested with any other version.

## Notes

The theme uses `Patternfly` for styling. The `Patternfly` files are located in the `common/resources/web_modules/react-core/dist/styles` directory.
Most js modules are written with `React`, however it also includes `Angular`. I've removed from the theme the unnecessary `React` modules because those are present in the system.

We are __NOT__ modifying any content within our theme other than styling.
