import {timestampToDate} from '@/api/convpb.js';
import {getDownloadDevicesUrl} from '@/api/ui/devices.js';
import {MINUTE} from '@/components/now.js';
import {computed, onScopeDispose, ref, toValue, watch} from 'vue';

/**
 * @param {import('vue').MaybeRefOrGetter<Partial<Device.Query.AsObject>>} query
 * @param {import('vue').MaybeRefOrGetter<Partial<Period.AsObject|undefined|null>>} [history]
 * @return {{
 *   downloadBtnProps: ComputedRef<Record<string,any>>
 * }}
 */
export function useDownloadLink(query, history = null) {
  const tableDownloadUrl = ref(/** @type {DownloadDevicesUrl.AsObject | null} */ null);
  const getDownloadDevicesUrlRequest = computed(() => {
    return {query: toValue(query), history: toValue(history)};
  });

  const leeway = 1 * MINUTE;
  let refreshHandle = 0;
  onScopeDispose(() => clearTimeout(refreshHandle));

  const fetchDownloadUrl = async (request) => {
    clearTimeout(refreshHandle);
    if (!request) return;
    try {
      const url = await getDownloadDevicesUrl(request);
      tableDownloadUrl.value = url;
      if (url.expireAfterTime) {
        const expireAtDate = timestampToDate(url.expireAfterTime);
        const expireIn = expireAtDate - Date.now() - leeway;
        if (expireIn > 0) {
          refreshHandle = setTimeout(() => fetchDownloadUrl(request), expireIn);
        }
      }
    } catch (e) {
      console.warn('Failed to get download devices URL', e);
    }
  };

  watch(getDownloadDevicesUrlRequest, async (request) => {
    tableDownloadUrl.value = null;
    await fetchDownloadUrl(request);
  }, {immediate: true});

  const downloadBtnProps = computed(() => {
    const props = {};
    const url = tableDownloadUrl.value;
    if (url) {
      props.href = url.url;
      props.download = url.filename;
    }
    props.disabled = !url;

    return props;
  });

  return {
    downloadBtnProps
  };
}
