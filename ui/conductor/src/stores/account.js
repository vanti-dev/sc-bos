import {events, keycloak} from '@/api/keycloak.js';
import localLogin from '@/api/localLogin.js';
import jwtDecode from 'jwt-decode';
import {defineStore} from 'pinia';
import {computed, ref} from 'vue';

export const useAccountStore = defineStore('accountStore', () => {
  const kcp = keycloak();
  const kcEvents = events;

  const loggedIn = ref(false);
  const token = ref('');
  const claims = ref({});
  const loginForm = ref(false);
  const loginDialog = ref(false);
  const snackbar = ref(false);

  const updateRefs = () => {
    kcp.then((kc) => {
      loggedIn.value = kc.authenticated;
      token.value = kc.token;
      claims.value = kc.idTokenParsed;
      localStorage.setItem('keyclock', true);
      saveLocalStorage();
    });
  };

  const toggleSnackbar = () => {
    snackbar.value = !snackbar.value;
  };

  const toggleLoginForm = () => {
    loginForm.value = !loginForm.value;
  };

  const toggleLoginDialog = () => {
    loginDialog.value = !loginDialog.value;
    loginForm.value = false;
  };

  const loginLocal = async (username, password) => {
    try {
      const res = await localLogin(username, password);
      if (res.status === 200) {
        const payload = await res.json();
        token.value = payload.access_token;
        loggedIn.value = true;
        claims.value = {
          email: username,
          ...jwtDecode(payload.access_token)
        };
        toggleLoginDialog();
        saveLocalStorage();
        localStorage.setItem('keyclock', false);
      } else {
        snackbar.value = true;
      }
    } catch (err) {
      snackbar.value = true;
    }
  };

  const saveLocalStorage = () => {
    localStorage.setItem('loggedIn', loggedIn.value);
    localStorage.setItem('token', token.value);
    localStorage.setItem('loggedIn', loggedIn.value);
    localStorage.setItem('claims', JSON.stringify(claims.value));
  };

  const loadLocalStorage = () => {
    token.value = localStorage.getItem('token');
    loggedIn.value = JSON.parse(localStorage.getItem('loggedIn'));
    claims.value = JSON.parse(localStorage.getItem('claims'));
  };

  const logout = async () => {
    localStorage.getItem('keyclock') === 'true' &&
    kcp.then((kc) => kc.logout());
    loggedIn.value = false;
    token.value = '';
    claims.value = {};
    saveLocalStorage();
  };

  kcEvents.addEventListener('authSuccess', updateRefs);

  return {
    loggedIn,
    token,
    claims,
    loginForm,
    loginDialog,
    toggleSnackbar,
    toggleLoginForm,
    toggleLoginDialog,
    loginLocal,
    loadLocalStorage,
    snackbar,

    fullName: computed(() => claims.value?.name || ''),
    email: computed(() => claims.value?.email || ''),

    login: (scopes) => {
      return kcp.then((kc) => kc.login({scope: scopes.join(' ')}));
    },
    logout
  };
});
