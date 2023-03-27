/**
 * @type {Promise<ServerConfig> | null}
 * @private
 */
let _uiConfig = null;

/**
 *
 */
export async function uiConfig() {
  if (_uiConfig === null) {
    /** @type {string} */
    const url = import.meta.env.VITE_UI_CONFIG_URL || '/ui-config.json';
    _uiConfig = await fetch(url)
        .then(res => /** @type {Promise<ServerConfig>} */ res.json())
        .catch(() => (_uiConfig = _defaultConfig));
  }
  // todo: retry on network failure
  return _uiConfig;
}

/**
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
      'overview': false,
      'emergency-lighting': true,
      'notifications': true
    },
    'automations': {
      '*': true
    },
    'site': {
      'zone': true
    }
  }
};

/**
 *
 * @param {string} path
 */
export async function featureEnabled(path) {
  const p = path.split('/');
  // if absolute path, first element will be empty
  if (p[0] === '') {
    p.splice(0, 1);
  }
  const c = await uiConfig();
  return checkFeatureEnabled(c.features, p);
}

/**
 * Checks whether the config has a feature enabled based on the path
 *
 * @param {Object} config
 * @param {string[]} pathSegments
 * @return {boolean}
 */
export function checkFeatureEnabled(config, pathSegments) {
  if (pathSegments.length === 0) {
    return true;
  }
  const seg = pathSegments[0];
  if (config.hasOwnProperty(seg)) {
    if (typeof config[seg] == 'object') {
      return checkFeatureEnabled(config[seg], pathSegments.slice(1));
    }
    return config[seg];
  } else if (config.hasOwnProperty('*')) {
    if (typeof config['*'] == 'object') {
      return checkFeatureEnabled(config['*'], pathSegments.slice(1));
    }
    return config[seg];
  }
  return false;
}
