import {acknowledgeAlert, listAlerts, pullAlerts, unacknowledgeAlert} from '@/api/ui/alerts.js';
import {useAccountStore} from '@/stores/account';
import {Collection} from '@/util/query.js';
import {Alert} from '@sc-bos/ui-gen/proto/alerts_pb';
import {acceptHMRUpdate, defineStore} from 'pinia';


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
   * @param {string} name
   */
  function setAcknowledged(e, alert, name='') {
    if (e) {
      let author = undefined;
      if (account.email || account.fullName) {
        author = {
          email: account.email,
          displayName: account.fullName
        };
      }
      acknowledgeAlert({
        name, id: alert.id, allowAcknowledged: false, allowMissing: false, author
      })
          .catch(err => console.error(err));
    } else {
      unacknowledgeAlert({name, id: alert.id, allowAcknowledged: false, allowMissing: false})
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
   * @param {string} name
   * @return {Collection}
   */
  function newCollection(name='') {
    const listFn = async (query, tracker, pageToken, recordFn) => {
      const page = await listAlerts({name, pageToken, query, pageSize: 100}, tracker);
      for (const alert of page.alertsList) {
        recordFn(alert, alert.id);
      }
      return page.nextPageToken;
    };
    const pullFn = (query, resources) => {
      pullAlerts({name, query}, resources);
    };
    return new Collection(listFn, pullFn);
  }

  return {
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
