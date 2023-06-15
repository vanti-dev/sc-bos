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
        v-model="drawer"
        app
        class="py-8 ml-4"
        clipped
        color="transparent"
        expand-on-hover
        floating
        :mini-variant.sync="miniVariant"
        :mini-variant-width="drawerWidth"
        width="275"
        permanent>
      <h1
          :class="[miniVariant ? 'text-subtitle-1 pl-2 my-5' : 'text-h1', 'pl-1 text-truncate']"
          style="maxHeight: 41px;">
        {{ pageTitle }}
      </h1>
      <v-divider class="my-5"/>
      <router-view
          v-if="hasNav"
          name="nav"
          class="mx-4"
          :style="miniVariant ? 'width: 40px; margin-top: 12px;' : 'width: auto;'"/>
      <template #append>
        <v-footer class="pa-0" style="background:transparent">
          <v-col class="pa-0">
            <v-divider/>
            <p class="my-2 text-caption text-center neutral--text text--lighten-2">
              Smart Core<br>{{ appVersion }}
            </p>
          </v-col>
        </v-footer>
      </template>
    </v-navigation-drawer>

    <router-view/>

    <error-view/>
  </v-app>
</template>

<script setup>
import {computed, ref, watch} from 'vue';

import AccountBtn from '@/components/AccountBtn.vue';
import AppMenu from '@/components/AppMenu.vue';
import {usePage} from '@/components/page.js';
import ScLogo from '@/components/ScLogo.vue';
import ErrorView from '@/components/ui-error/ErrorView.vue';
import {useAccountStore} from '@/stores/account.js';

const {pageTitle, hasSections, hasNav, hasSidebar} = usePage();

const drawer = ref(true);
const miniVariant = ref(true);
const drawerWidth = ref(70);

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
    drawerWidth.value = 70;
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
</style>
