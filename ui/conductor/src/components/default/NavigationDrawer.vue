<template>
  <v-navigation-drawer
      v-if="hasNav && isLoggedIn"
      v-model="drawer"
      app
      class="siteNavigation"
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
</template>

<script setup>
import {usePage} from '@/components/page';
import useAuthSetup from '@/composables/useAuthSetup';
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import {computed, watch} from 'vue';

const {isLoggedIn} = useAuthSetup();
const {hasNav} = usePage();

const {drawer, miniVariant, drawerWidth, pinDrawer} = storeToRefs(usePageStore());

const appVersion = computed(() => {
  if (GIT_VERSION.startsWith('ui/')) {
    return GIT_VERSION.substring(3);
  }
  return GIT_VERSION;
});

watch(
    miniVariant,
    (expanded) => {
      if (expanded) {
        drawerWidth.value = 45;
      } else {
        drawerWidth.value = 275;
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);
</script>

<style lang="scss" scoped>
.section-divider {
  border-color: currentColor;
}

.pin-sidebar-btn {
  width: 100%;
}

/** This helps displaying the notification counter badge, while keeping the right sidebar scrollable */
.siteNavigation,
.siteNavigation ::v-deep .v-navigation-drawer__content {
  overflow: visible;
}
</style>
