import {defineStore} from 'pinia';
import {computed, ref} from 'vue';
import {getMetadata} from '@/api/sc/traits/metadata';

export const useConfigStore = defineStore('config', () => {
  const zoneId = ref('');
  const zoneMeta = ref({});

  const zoneName = computed(() => zoneMeta.value?.appearance?.title ?? zoneId.value ?? '');

  const isConfigured = computed(() => {
    return Boolean(zoneId.value);
  });

  /**
   * @param {string} zone
   * @param {Metadata.AsObject} [meta]
   */
  async function setZone(zone, meta = null) {
    if (zone) {
      zoneId.value = zone;
      if (meta) {
        zoneMeta.value = meta;
      } else {
        zoneMeta.value = await getMetadata({name: zone});
      }
    }
  }

  /**
   *
   */
  function reset() {
    zoneId.value = '';
    zoneMeta.value = {};
  }

  return {
    zoneId,
    zoneMeta,
    zoneName,
    isConfigured,
    setZone,
    reset
  };
},
{persist: true});
