import {acknowledgeAlert, unacknowledgeAlert} from '@/api/ui/alerts.js';
import {useAccountStore} from '@/stores/account';
import {Alert} from '@sc-bos/ui-gen/proto/alerts_pb';
import {acceptHMRUpdate, defineStore} from 'pinia';

export const SeverityStrings = {
  [Alert.Severity.INFO]: 'INFO',
  [Alert.Severity.WARNING]: 'WARN',
  [Alert.Severity.SEVERE]: 'ALERT',
  [Alert.Severity.LIFE_SAFETY]: 'DANGER'
};
export const SeverityColor = {
  [Alert.Severity.INFO]: 'info',
  [Alert.Severity.WARNING]: 'warning',
  [Alert.Severity.SEVERE]: 'error',
  [Alert.Severity.LIFE_SAFETY]: 'error'
};

/**
 *
 * @param {Severity} severity
 * @return {{text: string, color: string, background: string}}
 */
export function severityData(severity) {
  for (let i = severity; i > 0; i--) {
    if (SeverityStrings[i]) {
      let str = SeverityStrings[i];
      if (i < severity) {
        str += '+' + (severity - i);
      }
      return {text: str, color: `${SeverityColor[i]}--text`, background: `${SeverityColor[i]}`};
    }
  }
  return {text: 'unspecified', color: 'text-gray', background: 'gray'};
}

export const useNotifications = defineStore('notifications', () => {
  const account = useAccountStore();

  /**
   *
   * @param {boolean} e
   * @param {Alert.AsObject} alert
   * @param {string} name
   */
  function setAcknowledged(e, alert, name = '') {
    if (e) {
      let author = undefined;
      if (account.email || account.fullName) {
        author = {
          email: account.email,
          displayName: account.fullName
        };
      }
      acknowledgeAlert({
        name,
        id: alert.id,
        allowAcknowledged: false,
        allowMissing: false,
        author
      }).catch((err) => console.error(err));
    } else {
      unacknowledgeAlert({name, id: alert.id, allowAcknowledged: false, allowMissing: false})
          .catch((err) => console.error(err));
    }
  }

  return {
    severityData,
    setAcknowledged
  };
});

// enable hot reload for this store
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useNotifications, import.meta.hot));
}
