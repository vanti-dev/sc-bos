<template>
  <v-app>
    <v-app-bar app height="60">
      <app-menu/>
      <sc-logo :fill="logoColor" :outline="logoColor ? 'white' : undefined" style="height: 35px; margin-left: 16px"/>
      <span class="heading">Smart Core</span>
      <page-title/>
      <v-spacer/>
      <v-btn>{{ loginText }}</v-btn>
    </v-app-bar>


    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>

<script setup>
import {computed} from 'vue';
import AppMenu from './components/AppMenu.vue';
import PageTitle from './components/PageTitle.vue';
import ScLogo from './components/ScLogo.vue';
import {useTheme} from './components/theme.js';
import {useAccountStore} from './stores/account.js';

const theme = useTheme();
const logoColor = theme.logoColor;

const accountStore = useAccountStore();
const loginText = computed(() => {
  if (accountStore.loggedIn) {
    return accountStore.account.name;
  } else {
    return 'Log in';
  }
});
</script>

<style scoped>
.v-app-bar ::v-deep(.v-toolbar__content > .v-btn.v-btn--icon:first-child) {
  margin-left: -16px;
  margin-top: -4px;
  margin-bottom: -4px;
  height: 60px;
  width: 60px;
}

.v-app-bar ::v-deep(.v-toolbar__content > .v-btn.v-btn--icon:first-child .v-btn__content) {
  height: 100%
}

.heading {
  font-size: 22px;
  font-weight: 300;
}
</style>
