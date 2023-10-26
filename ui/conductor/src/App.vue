<template>
  <v-app class="app-root">
    <v-app-bar
        app
        height="60"
        :clipped-left="hasNav"
        :clipped-right="hasSidebar"
        elevation="0"
        class="pr-7">
      <app-menu v-if="isLoggedIn"/>
      <sc-logo outline="white" style="height: 35px; margin-left: 16px"/>
      <span class="heading">Smart Core {{ isLoggedIn ? ' | ' + pageTitle : '' }}</span>

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
        v-if="hasNav && isLoggedIn"
        v-model="drawer"
        app
        :class="[!miniVariant ? 'pt-0' : 'pt-3', 'pb-8 ml-2']"
        clipped
        color="transparent"
        :expand-on-hover="!pinDrawer"
        floating
        :mini-variant.sync="miniVariant"
        :mini-variant-width="drawerWidth"
        width="275"
        permanent>
      <v-btn
          v-if="hasNav && !miniVariant"
          x-small
          text
          class="d-block neutral--text text--lighten-4 text-caption text-center ma-0 pa-0 mb-n3 ml-1 mt-1"
          width="100%"
          @click="pinDrawer = !pinDrawer">
        {{ !pinDrawer ? 'Pin navigation' : 'Unpin navigation' }}
      </v-btn>
      <router-view
          v-if="hasNav"
          name="nav"
          class="ml-1 mt-4"
          :style="miniVariant ? 'width: 40px;' : 'width: auto;'"/>
      <template v-if="!miniVariant" #append>
        <v-footer class="pa-0" style="background:transparent">
          <v-col class="pa-0">
            <v-divider/>
            <p class="mt-2 mb-n4 text-caption text-center neutral--text text--lighten-2">
              Smart Core {{ appVersion }}
            </p>
          </v-col>
        </v-footer>
      </template>
    </v-navigation-drawer>
    <router-view v-if="isLoggedIn"/>

    <error-view/>
  </v-app>
</template>

<script setup>
import {computed, watch} from 'vue';
import {storeToRefs} from 'pinia';

import AccountBtn from '@/components/AccountBtn.vue';
import AppMenu from '@/components/AppMenu.vue';
import {usePage} from '@/components/page.js';
import ScLogo from '@/components/ScLogo.vue';
import ErrorView from '@/components/ui-error/ErrorView.vue';
import {useAccountStore} from '@/stores/account.js';
import {usePageStore} from '@/stores/page';

import useAuthSetup from '@/composables/useAuthSetup';
const {isLoggedIn} = useAuthSetup();

const {pageTitle, hasSections, hasNav, hasSidebar} = usePage();

const {drawer, miniVariant, drawerWidth, pinDrawer} = storeToRefs(usePageStore());

const store = useAccountStore();
store.loadLocalStorage();

const appVersion = computed(() => {
  if (GITVERSION.startsWith('ui/')) {
    return GITVERSION.substring(3);
  }
  return GITVERSION;
});

watch(miniVariant, expanded => {
  if (expanded) {
    drawerWidth.value = 45;
  } else {
    drawerWidth.value = 275;
  }
}, {immediate: true, deep: true, flush: 'sync'});
</script>

<style scoped>
.v-application {
  background: var(--v-neutral-darken1);
}

.v-app-bar.v-toolbar.v-sheet {
  background: var(--v-neutral-base);
}

.v-app-bar :deep(.v-toolbar__content) {
  padding-right: 0px;
}

.heading {
  font-size: 22px;
  font-weight: 300;
}

.section-divider {
  border-color: currentColor;
}

.pin-sidebar-btn {
  width: 100%;
}

::v-deep .v-navigation-drawer,
::v-deep .v-navigation-drawer--mini-variant,
::v-deep .v-navigation-drawer__content {
  overflow: visible;
}
</style>
