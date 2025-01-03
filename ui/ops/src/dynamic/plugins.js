import {newActionTracker, trackAction} from '@/api/resource.js';
import deepEqual from 'fast-deep-equal';
import {RpcError, StatusCode} from 'grpc-web';
import {computed, reactive, ref, toValue, watch} from 'vue';

/**
 * @param {MaybeRefOrGetter<PluginDesc>} desc
 * @return {UsePluginResponse}
 */
export function usePlugin(desc) {
  const tracker =
      /** @type {import('@/api/resource').ActionTracker<Plugin>} */
      reactive(newActionTracker());

  const loaded = ref(false);
  let runCount = 0;

  watch(() => /** @type {PluginDesc|undefined} */ toValue(desc), (n, o) => {
    if (!n || deepEqual(n, o)) return;
    loaded.value = false;
    const run = ++runCount;
    trackAction('PluginApi.GetPlugin', tracker, async (_) => {
      const plugin = await getPlugin(desc.name, n.category, n.type);
      // make the response look like it came from a gRPC call
      return {
        toObject() {
          return plugin;
        }
      };
    })
        .catch(() => {}) // the tracker handles this, avoid console warnings
        .finally(() => {
          if (run === runCount) {
            loaded.value = true;
          }
        });
  }, {immediate: true, deep: true});

  return {
    loaded,
    plugin: computed(() => tracker.response),
    loading: computed(() => tracker.loading),
    error: computed(() => tracker.error?.error ?? tracker.error)
  };
}

/**
 * @param {MaybeRefOrGetter<PluginDesc>} desc
 * @return {UseServicePluginResponse}
 */
export function useServicePlugin(desc) {
  const {plugin, ...rest} = usePlugin(desc);
  return {
    plugin: computed(() => plugin.value?.service),
    ...rest
  };
}

/**
 * @typedef {Object} PluginDesc
 * @property {string} type - the registered service type, e.g. 'bacnet' or 'area' or 'metermail'
 * @property {('automation'|'driver'|'zone'|'system')} category - the type of service this is
 * @property {string} name - the name of the device that hosts the service
 */

/**
 * @typedef {Object} UsePluginResponse
 * @property {Ref<Plugin | undefined>} plugin
 * @property {Ref<boolean>} loaded
 * @property {Ref<boolean>} loading
 * @property {Ref<Error | undefined>} error
 */

/**
 * @typedef {UsePluginResponse} UseServicePluginResponse
 * @property {Ref<ServicePlugin | undefined>} plugin
 */

/**
 * @typedef {Object} CategoryPlugin
 * @property {Partial<Metadata.AsObject>?} metadata
 * @property {Record<string, import('vue').Component>?} slots
 * @property {string?} defaultRoute
 * @property {import('vue-router').RouteRecordRaw[]?} routes
 */

/**
 * @typedef {CategoryPlugin} ServicePlugin
 */

/**
 * @typedef {Object} Plugin
 * @property {ServicePlugin?} service
 */

// todo: we should load this information from the server, once we figure out how to do that
/**
 * @type {Record<string, Record<string, Record<string, () => Promise<Plugin>>>>}
 */
const plugins = {
  '*': {
    'drivers': {
      'mock': () => import('@/dynamic/plugins/drivers/mock/mock.js')
    }
  }
};

// error thrown when a plugin can't be found.
const notFound = (reason) => new RpcError(StatusCode.NOT_FOUND, 'plugin not found: ' + reason, {});

/**
 * @param {string} name
 * @param {string} category
 * @param {string} type
 * @return {Promise<Plugin>}
 */
async function getPlugin(name, category, type) {
  const categories = plugins[name] ?? plugins['*'];
  if (!categories) throw notFound(`name ${name}`);
  const types = categories[category] ?? categories['*'];
  if (!types) throw notFound(`category ${category}`);
  const plugin = types[type] ?? types['*'];
  if (!plugin) throw notFound(`type ${type}`);
  const mod = await plugin();
  return mod.plugin ?? {
    service: mod.service
  };
}
