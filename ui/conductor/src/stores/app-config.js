import {defineStore} from 'pinia';
import {computed, ref} from 'vue';

export const useAppConfigStore = defineStore('appConfig', () => {
  /**
   * @private
   */
  const _config = ref({});
  let _configResolve;
  const configPromise = new Promise((resolve) => _configResolve = resolve);

  /**
   * The default config for the UI - this should mostly be targeted as though it was running on an Area Controller, as
   * this will have the most standardised feature set.
   *
   * @private
   */
  const _defaultConfig = {
    hub: false, // specifies if we're talking to a hub or an area controller
    proxy: false, // specifies if we're using querying via an proxy (e.g. EdgeGateway) or not
    features: {
      'auth': {
        'users': true,
        'third-party': true
      },
      'devices': {
        '*': true
      },
      'ops': {
        'overview': false,
        'emergency-lighting': false,
        'notifications': true
      },
      'automations': {
        '*': true
      },
      'site': {
        'zone': {
          '*': true
        }
      },
      'system': {
        'drivers': true,
        'features': true
      }
    },
    config: {
      home: '/devices',
      hub: false, // specifies if we're talking to a hub or an area controller
      proxy: false // specifies if we're using querying via an proxy (e.g. EdgeGateway) or not
    }
  };


  /**
   * Loads the config from the server
   */
  async function loadConfig() {
    const url = import.meta.env.VITE_UI_CONFIG_URL || '/__/scos/ui-config.json';
    try {
      const res = await fetch(url);
      _config.value = await res.json();
    } catch (e) {
      console.warn('Failed to load config from server, using default config', e);
      _config.value = _defaultConfig;
    }
    _configResolve(config.value);
  }

  /**
   * A list of regexp paths for all the enabled features.
   * These are provided as RegExp to handle wildcard paths (e.g. /devices/*)
   *
   * @type {ComputedRef<RegExp[]>}
   */
  const enabledPaths = computed(() => {
    let features = [];
    if (_config.value.hasOwnProperty('features')) {
      // generate paths list, then convert each one to regex
      features = _generatePaths('', _config.value.features).map(path => {
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
        paths.push(`${prefix}/${key}`);
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

  const config = computed(() => _config.value?.config ?? {});

  return {
    loadConfig,
    enabledPaths,
    pathEnabled,
    config,
    configPromise,
    homePath: computed(() => _config.value?.config?.home ?? _defaultConfig.config.home)
  };
});
