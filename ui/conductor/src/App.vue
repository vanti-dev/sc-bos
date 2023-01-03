<template>
  <v-app class="app-root">
    <v-app-bar
        app
        height="60"
        :clipped-left="hasNav"
        elevation="0"
        class="pr-7">
      <app-menu btn-class="full-btn"/>
      <sc-logo outline="white" style="height: 35px; margin-left: 16px"/>
      <span class="heading">Smart Core</span>

      <v-divider
          vertical
          v-if="hasSections"
          class="mx-8 section-divider"
          inset/>

      <v-spacer/>

      <router-view name="actions"/>
      <account-btn btn-class="full-btn mr-0"/>
    </v-app-bar>

    <v-navigation-drawer
        v-if="hasNav"
        app
        clipped
        permanent
        color="transparent"
        floating
        width="275"
        class="pa-8 pr-0 mr-8">
      <page-title/>
      <v-divider class="my-5"/>
      <router-view name="nav" v-if="hasNav"/>
      <template #append>
        <NavFooter/>
      </template>
    </v-navigation-drawer>

    <v-main class="mx-10 my-6">
      <router-view/>
    </v-main>
  </v-app>
</template>

<script setup>
import AccountBtn from '@/components/AccountBtn.vue';
import AppMenu from '@/components/AppMenu.vue';
import {usePage} from '@/components/page.js';
import PageTitle from '@/components/PageTitle.vue';
import ScLogo from '@/components/ScLogo.vue';
import {useAccountStore} from '@/stores/account.js';
// import {onMounted} from 'vue';

const {hasSections, hasNav} = usePage();

const store = useAccountStore();

store.loadLocalStorage();
</script>

<style scoped>
.v-application {
  background: var(--v-neutral-darken1);
}

.v-app-bar.v-toolbar.v-sheet {
  background: var(--v-neutral-base);
}

.v-app-bar ::v-deep(.v-toolbar__content) {
  padding-right: 0px;
}

.v-app-bar ::v-deep(.v-toolbar__content > .v-btn.full-btn) {
  margin-top: -4px;
  margin-bottom: -4px;
  height: 60px;
  min-width: 60px;
  width: auto;

  background-color: transparent;
}

.v-app-bar ::v-deep(.v-toolbar__content > .v-btn.full-btn:first-child) {
  margin-left: -16px;
}

.v-app-bar ::v-deep(.v-toolbar__content > .v-btn.full-btn:last-child) {
  margin-right: -16px;
}

.v-app-bar ::v-deep(.v-toolbar__content > .v-btn.full-btn .v-btn__content) {
  height: 100%;
}

.heading {
  font-size: 22px;
  font-weight: 300;
}

.section-divider {
  border-color: currentColor;
}
</style>
