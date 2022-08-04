import {defineStore} from 'pinia';

export const useAccountStore = defineStore('accountStore', {
  state: () => ({
    account: null
  }),
  getters: {
    loggedIn: state => Boolean(state.account)
  }
})
