import {closeResource, newActionTracker} from '@/api/resource';
import {getMetadata} from '@/api/sc/traits/metadata';
import {computed, reactive, toValue, watch} from 'vue';

/**
 * Collect metadata for a collection of device names.
 *
 * @param {MaybeRefOrGetter<string|string[]>} devices
 * @return {{
 *   loading: import('vue').Ref<boolean>,
 *   trackers: {string: ActionTracker<Metadata.AsObject>}
 * }}
 */
export default function useMetadata(devices) {
  const deviceArr = computed(() => {
    const ds = toValue(devices);
    if (Array.isArray(ds)) return ds;
    return [ds];
  });

  const trackers = reactive(
      /** @type {{string:ActionTracker<Metadata.AsObject>}} */
      {});
  /**
   * Return items in a that aren't in b, and items in b that aren't in a.
   *
   * @param {string[]} a
   * @param {string[]} b
   * @return {{inA: Set<string>, inB: Set<string>}}
   */
  const diff = (a, b) => {
    const inA = new Set(a);
    const inB = new Set(b);
    if (b) {
      for (const s of b) {
        inA.delete(s);
      }
    }
    if (a) {
      for (const s of a) {
        inB.delete(s);
      }
    }
    return {inA, inB};
  };
  watch(deviceArr, (devices, oldDevices) => {
    const {inA: added, inB: removed} = diff(devices, oldDevices);
    // handle remove first to avoid excessive resource utilization
    for (const n of removed) {
      const t = trackers[n];
      closeResource(t);
      delete(trackers[n]);
    }
    for (const n of added) {
      const t = reactive(newActionTracker());
      getMetadata({name: n}, t)
          // ignore error, captured as part of the tracker.
          .catch(() => {});
      trackers[n] = t;
    }
  }, {immediate: true, deep: true});

  /**
   * Report whether any of the configured devices are loading - aka fetching data from the server.
   *
   * @type {import('vue').ComputedRef<boolean>}
   */
  const loading = computed(() => {
    return Object.values(trackers).some(t => t.loading);
  });

  return {
    trackers,
    loading
  };
}
