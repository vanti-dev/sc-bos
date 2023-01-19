<template>
  <v-app class="app-root">
    <v-app-bar
        app
        height="60"
        :clipped-left="hasNav"
        :clipped-right="hasSidebar"
        elevation="0"
        class="pr-7">
      <app-menu/>
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
      <h1 class="pl-1 text-h1">{{ pageTitle }}</h1>
      <v-divider class="my-5"/>
      <router-view name="nav" v-if="hasNav"/>
      <template #append>
        <v-footer class="pa-0" style="background:transparent">
          <v-col class="pa-0">
            <v-divider/>
            <p class="my-2 text-caption text-center neutral--text text--lighten-2">Smart Core v2022.11</p>
          </v-col>
        </v-footer>
      </template>
    </v-navigation-drawer>

    <router-view/>
  </v-app>
</template>

<script setup>
import AccountBtn from '@/components/AccountBtn.vue';
import AppMenu from '@/components/AppMenu.vue';
import {usePage} from '@/components/page.js';
import ScLogo from '@/components/ScLogo.vue';
import {useAccountStore} from '@/stores/account.js';
// import {onMounted} from 'vue';

const {pageTitle, hasSections, hasNav, hasSidebar} = usePage();

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

.heading {
  font-size: 22px;
  font-weight: 300;
}

.section-divider {
  border-color: currentColor;
}
</style>
