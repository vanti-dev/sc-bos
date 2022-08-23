import {events, keycloak} from '@/api/keycloak.js';
import {defineStore} from 'pinia';
import {computed, ref} from 'vue';

export const useAccountStore = defineStore('accountStore', () => {
  const kcp = keycloak();
  const kcEvents = events;

  const loggedIn = ref(false);
  const token = ref('');
  const claims = ref({});

  const updateRefs = () => {
    kcp.then(kc => {
      loggedIn.value = kc.authenticated;
      token.value = kc.token;
      claims.value = kc.idTokenParsed;
    })
  }
  kcEvents.addEventListener('ready', updateRefs);
  kcEvents.addEventListener('authError', updateRefs);
  kcEvents.addEventListener('authSuccess', updateRefs);
  kcEvents.addEventListener('authRefreshError', updateRefs);
  kcEvents.addEventListener('authRefreshSuccess', updateRefs);
  kcEvents.addEventListener('authLogout', updateRefs);
  kcEvents.addEventListener('tokenExpired', updateRefs);

  return {
    loggedIn,
    token,
    claims,

    fullName: computed(() => claims.value?.name || ''),
    email: computed(() => claims.value?.email || ''),

    login: (scopes) => {
      return kcp.then(kc => kc.login({scope: scopes.join(' ')}))
    },
    logout: () => {
      return kcp.then(kc => kc.logout())
    }
  }
})
