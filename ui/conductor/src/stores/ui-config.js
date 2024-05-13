import {toValue} from '@/util/vue.js';
import {defineStore} from 'pinia';
import {computed, ref} from 'vue';

export const useUiConfigStore = defineStore('uiConfig', () => {
  /**
   * @private
   */
  const _config = ref({});
  const _loaded = ref(false);
  let _configResolve;
  const configPromise = new Promise((resolve) => _configResolve = resolve);

  /**
   * Loads the config from the server
   */
  async function loadConfig() {
    if (_loaded.value) {
      return;
    }
    const url = import.meta.env.VITE_UI_CONFIG_URL || '/__/scos/ui-config.json';
    try {
      const res = await fetch(url);
      _config.value = await res.json();
    } catch (e) {
      console.warn('Failed to load config from server, using default config', e);
      _config.value = _defaultConfig;
    }
    _configResolve(config.value);
    _loaded.value = true;
  }

  const config = computed(() => _config.value?.config ?? {});

  /**
   * Gets the value of path from either uiConfig config or defaultConfig, depending on presence.
   *
   * @template T
   * @param {string} path
   * @param {T?} def
   * @return {T}
   */
  const getOrDefault = (path, def) => {
    const parts = path.split('.');
    let a = config.value;
    let b = _defaultConfig?.config;
    for (let i = 0; i < parts.length; i++) {
      a = a?.[parts[i]];
      b = b?.[parts[i]];
    }
    return a ?? b ?? toValue(def);
  };

  return {
    loadConfig,
    config,
    configPromise,
    defaultConfig: _defaultConfig,
    getOrDefault,
    ...useSiteMap(_config),
    ...useTheme(_config)
  };
});

/**
 * Exposes site map relating features of the UI config.
 *
 * @param {MaybeRefOrGetter<Object>} config
 * @return {{
 *   enabledPaths: import('vue').ComputedRef<RegExp[]>,
 *   homePath: ComputedRef<unknown>,
 *   pathEnabled: (function(string): boolean)
 * }}
 */
export function useSiteMap(config) {
  /**
   * A list of regexp paths for all the enabled features.
   * These are provided as RegExp to handle wildcard paths (e.g. /devices/*)
   *
   * @type {import('vue').ComputedRef<RegExp[]>}
   */
  const enabledPaths = computed(() => {
    let features = [];
    const _config = toValue(config);
    if (_config.hasOwnProperty('features')) {
      // generate paths list, then convert each one to regex
      features = _generatePaths('', _config.features).map(path => {
        return new RegExp(`^${path.replace(/\*/g, '.*')}$`);
      });
    }
    return features;
  });

  /**
   * Recursively generates a list of paths for all the enabled features
   *
   * @param {string} prefix
   * @param {Object} obj
   * @return {string[]}
   * @private
   */
  function _generatePaths(prefix, obj) {
    const paths = [];
    for (const [key, value] of Object.entries(obj)) {
      if (value === true) {
        let path = `${prefix}/${key}`;
        if (key !== '*') {
          path += '(/*)?';
        }
        paths.push(path);
      } else if (value === false) {
        // do nothing
      } else {
        paths.push(`${prefix}/${key}`);
        paths.push(..._generatePaths(`${prefix}/${key}`, value));
      }
    }
    return paths;
  }

  /**
   * Checks whether the path matches against any of the RegExp in the enabledPaths array
   *
   * @param {string} path
   * @return {boolean}
   */
  function pathEnabled(path) {
    for (const regex of enabledPaths.value) {
      if (path.match(regex)) {
        return true;
      }
    }
    return false;
  }

  return {
    enabledPaths,
    pathEnabled,
    homePath: computed(() => toValue(config)?.config?.home ?? _defaultConfig.config.home)
  };
}

/**
 * @param {MaybeRefOrGetter<Object>} config
 * @return {{
 *   appBranding: ComputedRef<Object>
 * }}
 */
export function useTheme(config) {
  // Returns the app branding, merging the default and the config values
  // The config values will override the default ones if any are present
  const appBranding = computed(() => {
    const _config = toValue(config);
    return {
      ..._defaultConfig.config.theme.appBranding,
      ..._config?.config?.theme?.appBranding
    };
  });
  return {
    appBranding
  };
}

/**
 * The default config for the UI - this should mostly be targeted as though it was running on an Area Controller, as
 * this will have the most standardised feature set.
 *
 * @private
 */
const _defaultConfig = {
  features: {
    'auth': {
      'users': true,
      'third-party': true
    },
    'devices': {
      '*': true
    },
    'ops': {
      'overview': {
        '*': true
      },
      'emergency-lighting': false,
      'notifications': true
    },
    'automations': {
      '*': true
    },
    'site': false,
    'system': {
      'drivers': true,
      'features': true
    }
  },
  config: {
    'keycloak': false,
    'home': '/devices',
    'ops': {
      'overview': {
        // Specifies all traits to show or hide on the building page
        'widgets': {
          'showEnergyConsumption': {
            'showChart': true,
            'showIntensity': true
          },
          'showEnvironment': true,
          'showOccupancy': true
        },
        'children': [] // Specifies the children of the building (zones/floors)
      }
    },
    'theme': {
      'appBranding': {
        'brandName': 'Smart Core', // The name of the brand
        'brandLogo': {
          'altText': 'Smart Core logo - representing nodes and connections', // Alt text for the logo
          'src': '' // Empty string will use the default logo
        },
        // Specifies brand colors for the app (theming)
        'brandColors': {
          'primary': {
            'base': '#00BED6',
            'darken3': '#338fa1'
          }
        }
      }
    },
    'hub': false, // Specifies if we're talking to a hub or an area controller
    'proxy': false, // Specifies if we're using querying via a proxy (e.g. EdgeGateway) or not
    'disableAuthentication': false // Specifies if we're using authentication or not
  }
};
