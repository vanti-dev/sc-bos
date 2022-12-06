import {closeResource, newActionTracker, newResourceCollection} from '@/api/resource.js';
import {acknowledgeAlert, listAlerts, pullAlerts, unacknowledgeAlert} from '@/api/ui/alerts.js';
import {Alert, ListAlertsResponse} from '@bsp-ew/ui-gen/proto/alerts_pb';
import {acceptHMRUpdate, defineStore} from 'pinia';
import {computed, reactive, set, watch} from 'vue';


const SeverityStrings = {};
for (const [name, val] of Object.entries(Alert.Severity)) {
  SeverityStrings[val] = name;
}

export const useNotifications = defineStore('notifications', () => {
  // todo: get the name from somewhere
  const name = computed(() => 'test-ac');

  const alerts = reactive(/** @type {ResourceCollection<Alert.AsObject, Alert>} */newResourceCollection()); // holds all the alerts we can show
  const fetchingPage = reactive(/** @type {ActionTracker<ListAlertsResponse.AsObject>} */ newActionTracker()); // tracks the fetching of a single page

  watch(name, async name => {
    closeResource(alerts);
    pullAlerts({name}, alerts);
    try {
      const firstPage = await listAlerts({name, pageSize: 100, pageToken: undefined}, fetchingPage);
      for (let alert of firstPage.alertsList) {
        set(alerts.value, alert.id, alert);
      }
      fetchingPage.response = null;
    } catch (e) {
      console.warn('Error fetching first page', e);
    }
  }, {immediate: true});

  function severityData(severity) {
    for (let i = severity; i > 0; i--) {
      if (SeverityStrings[i]) {
        let str = SeverityStrings[i];
        if (i < severity) {
          str += '+' + (severity - i);
        }
        return str;
      }
    }
    return 'unspecified';
  }

  function setAcknowledged(e, alert) {
    if (e) {
      acknowledgeAlert({name: name.value, id: alert.id, allowAcknowledged: false, allowMissing: false})
          .catch(err => console.error(err));
    } else {
      unacknowledgeAlert({name: name.value, id: alert.id, allowAcknowledged: false, allowMissing: false})
          .catch(err => console.error(err));
    }
  }

  function isAcknowledged(alert) {
    return Boolean(alert.acknowledgement);
  }

  return {
    name,
    alerts,
    loading: computed(() => alerts.loading || fetchingPage.loading),
    severityData,
    setAcknowledged,
    isAcknowledged
  }
});

// enable hot reload for this store
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useNotifications, import.meta.hot))
}
