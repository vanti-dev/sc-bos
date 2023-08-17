import {defineStore} from 'pinia';
import {ref} from 'vue';

import {StatusLog} from '@sc-bos/ui-gen/proto/status_pb';

export const useStatusBarStore = defineStore('statusBarStore', () => {
  const showClose = ref(false);

  const setBarColor = (level, grant) => {
    let color = 'transparent';

    if (level !== null && level >= 0) {
      if (level <= StatusLog.Level.NOMINAL) color = 'red';
      if (level <= StatusLog.Level.NOTICE) color = 'info';
      if (level <= StatusLog.Level.REDUCED_FUNCTION) color = 'warning';
      if (level <= StatusLog.Level.NON_FUNCTIONAL) color = 'error';
      if (level <= StatusLog.Level.OFFLINE) color = 'grey';
    }

    if (grant !== null) {
      if (grant === 'granted') {
        color = 'granted';

        (function(currentGrant) {
          setTimeout(() => {
            // Check if the grant was 'granted' at the time the timeout was set
            if (currentGrant === 'granted') {
              color = 'transparent';
            }
          }, 3000);
        })(grant);
      } else color = grant;
    }

    return color;
  };

  return {
    showClose,

    setBarColor
  };
});
