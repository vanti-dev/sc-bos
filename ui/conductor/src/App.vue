<template>
  <v-app class="app-root">
    <v-app-bar app height="60">
      <app-menu btn-class="full-btn"/>
      <sc-logo :fill="logoColor" :outline="logoColor ? 'white' : undefined" style="height: 35px; margin-left: 16px"/>
      <span class="heading">Smart Core</span>
      <page-title/>

      <v-spacer/>
      <account-btn btn-class="full-btn mr-0"/>
    </v-app-bar>


    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>

<script setup>
import AccountBtn from './components/AccountBtn.vue';
import AppMenu from './components/AppMenu.vue';
import PageTitle from './components/PageTitle.vue';
import ScLogo from './components/ScLogo.vue';
import {useTheme} from './components/theme.js';

const theme = useTheme();
const logoColor = theme.logoColor;
</script>

<style scoped>
.v-application {
  background: #101820;
}

.v-app-bar {
  background: #283037;
}

.v-app-bar ::v-deep(.v-toolbar__content) {
  padding-right: 30px;
}

.v-app-bar ::v-deep(.v-toolbar__content > .v-btn.full-btn) {
  margin-top: -4px;
  margin-bottom: -4px;
  height: 60px;
  min-width: 60px;
  width: auto;
}

.v-app-bar ::v-deep(.v-toolbar__content > .v-btn.full-btn:first-child) {
  margin-left: -16px;
}

.v-app-bar ::v-deep(.v-toolbar__content > .v-btn.full-btn:last-child) {
  margin-right: -16px;
}

.v-app-bar ::v-deep(.v-toolbar__content > .v-btn.full-btn .v-btn__content) {
  height: 100%
}

.heading {
  font-size: 22px;
  font-weight: 300;
}
</style>
