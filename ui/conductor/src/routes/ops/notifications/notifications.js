import {closeResource, newActionTracker, newResourceCollection} from '@/api/resource';
import {acknowledgeAlert, listAlerts, pullAlerts, unacknowledgeAlert} from '@/api/ui/alerts.js';
import {useAccountStore} from '@/stores/account';
import {useControllerStore} from '@/stores/controller';
import {Collection} from '@/util/query.js';
import {Alert} from '@sc-bos/ui-gen/proto/alerts_pb';
import {acceptHMRUpdate, defineStore} from 'pinia';
import {computed, onMounted, onUnmounted, reactive, set, watch} from 'vue';
import {useErrorStore} from '@/components/ui-error/error';


const SeverityStrings = {
  [Alert.Severity.INFO]: 'INFO',
  [Alert.Severity.WARNING]: 'WARN',
  [Alert.Severity.SEVERE]: 'ALERT',
  [Alert.Severity.LIFE_SAFETY]: 'DANGER'
};
const SeverityColor = {
  [Alert.Severity.INFO]: 'info',
  [Alert.Severity.WARNING]: 'warning',
  [Alert.Severity.SEVERE]: 'error',
  [Alert.Severity.LIFE_SAFETY]: 'error'
};

export const useNotifications = defineStore('notifications', () => {
  const controller = useControllerStore();

  // holds all the alerts we can show
  const alerts = reactive(/** @type {ResourceCollection<Alert.AsObject, Alert>} */newResourceCollection());
  // tracks the fetching of a single page
  const fetchingPage = reactive(/** @type {ActionTracker<ListAlertsResponse.AsObject>} */ newActionTracker());

  watch(() => controller.controllerName, async name => {
    closeResource(alerts);
    pullAlerts({name}, alerts);
    try {
      const firstPage = await listAlerts({name, pageSize: 100, pageToken: undefined}, fetchingPage);
      for (const alert of firstPage.alertsList) {
        set(alerts.value, alert.id, alert);
      }
      fetchingPage.response = null;
    } catch (e) {
      console.warn('Error fetching first page', e);
    }
  }, {immediate: true});

  // UI Error Handling
  const errorStore = useErrorStore();
  let unwatchAlertErrors; let unwatchPageErrors;
  onMounted(() => {
    unwatchAlertErrors = errorStore.registerCollection(alerts);
    unwatchPageErrors = errorStore.registerValue(fetchingPage);
  });
  onUnmounted(() => {
    if (unwatchAlertErrors) unwatchAlertErrors();
    if (unwatchPageErrors) unwatchPageErrors();
  });

  /**
   *
   * @param {Severity} severity
   * @return {{color: string, text: *}|{color: string, text: string}}
   */
  function severityData(severity) {
    for (let i = severity; i > 0; i--) {
      if (SeverityStrings[i]) {
        let str = SeverityStrings[i];
        if (i < severity) {
          str += '+' + (severity - i);
        }
        return {text: str, color: `${SeverityColor[i]}--text`};
      }
    }
    return {text: 'unspecified', color: 'gray--text'};
  }

  const account = useAccountStore();

  /**
   *
   * @param {boolean} e
   * @param {Alert.AsObject} alert
   */
  function setAcknowledged(e, alert) {
    if (e) {
      let author = undefined;
      if (account.email || account.fullName) {
        author = {
          email: account.email,
          displayName: account.fullName
        };
      }
      acknowledgeAlert({
        name: controller.controllerName, id: alert.id, allowAcknowledged: false, allowMissing: false, author
      })
          .catch(err => console.error(err));
    } else {
      unacknowledgeAlert({name: controller.controllerName, id: alert.id, allowAcknowledged: false, allowMissing: false})
          .catch(err => console.error(err));
    }
  }

  /**
   *
   * @param {Alert.AsObject} alert
   * @return {boolean}
   */
  function isAcknowledged(alert) {
    return Boolean(alert.acknowledgement);
  }

  /**
   *
   * @return {Collection}
   */
  function newCollection() {
    const listFn = async (query, tracker, pageToken, recordFn) => {
      const page = await listAlerts({name: controller.controllerName, pageToken, query, pageSize: 100}, tracker);
      for (const alert of page.alertsList) {
        recordFn(alert, alert.id);
      }
      return page.nextPageToken;
    };
    const pullFn = (query, resources) => {
      pullAlerts({name: controller.controllerName, query}, resources);
    };
    return new Collection(listFn, pullFn);
  }

  return {
    name,
    alerts,
    loading: computed(() => alerts.loading || fetchingPage.loading),
    newCollection,
    severityData,
    setAcknowledged,
    isAcknowledged
  };
});

// enable hot reload for this store
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useNotifications, import.meta.hot));
}
