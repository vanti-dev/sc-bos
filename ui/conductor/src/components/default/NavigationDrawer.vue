<template>
  <v-navigation-drawer
      v-if="hasNav && isLoggedIn"
      v-model="drawer"
      class="siteNavigation pt-4"
      color="transparent"
      :expand-on-hover="!pinDrawer"
      floating
      v-model:rail="miniVariant"
      :rail-width="drawerWidth"
      width="275"
      permanent>
    <router-view
        v-if="hasNav"
        name="nav"
        class="nav-btns"/>
    <template v-if="!miniVariant" #append>
      <div class="pa-2">
        <v-btn
            v-if="hasNav && !miniVariant"
            variant="text"
            block
            @click="pinDrawer = !pinDrawer">
          <v-icon start>mdi-pin-outline</v-icon>
          {{ !pinDrawer ? 'Pin navigation' : 'Unpin navigation' }}
        </v-btn>
      </div>
      <div class="text-caption text-center text-neutral-lighten-2 text-no-wrap">
        <v-divider/>
        <p class="mt-2 mb-0">
          Smart Core &copy; {{ new Date().getFullYear() }}
        </p>
        <p :title="appVersion" class="mt-0" style="cursor: default">
          {{ displayVersion }}
        </p>
      </div>
    </template>
  </v-navigation-drawer>
</template>

<script setup>
import {usePage} from '@/components/page';
import useAuthSetup from '@/composables/useAuthSetup';
import {useNavStore} from '@/stores/nav';
import {storeToRefs} from 'pinia';
import {computed} from 'vue';

const {isLoggedIn} = useAuthSetup();
const {hasNav} = usePage();

const {drawer, miniVariant, drawerWidth, pinDrawer} = storeToRefs(useNavStore());

const appVersion = computed(() => {
  if (GIT_VERSION.startsWith('ui/')) {
    return GIT_VERSION.substring(3);
  }
  return GIT_VERSION;
});
const displayVersion = computed(() => {
  if (import.meta.env.MODE !== 'production') {
    // fake version for screen shots, etc
    const now = new Date(); // not reactive, but we only care about year and month, and we're local
    return `${now.getFullYear()}.${now.getMonth() + 1}`;
  }
  return appVersion.value.split('-')[0];
});
</script>

<style lang="scss" scoped>
.section-divider {
  border-color: currentColor;
}

.pin-sidebar-btn {
  width: 100%;
}

/** This helps displaying the notification counter badge, while keeping the right sidebar scrollable */
.siteNavigation.v-navigation-drawer--rail :deep(.v-list) {
  overflow: visible;
}

.siteNavigation.v-navigation-drawer--rail :deep(.v-navigation-drawer__content) {
  overflow: hidden;
}

/** this fixes indentation of list items in groups when showing rail variant */
.siteNavigation :deep(.v-list-group__items .v-list-item) {
  transition: padding-inline-start 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.siteNavigation.v-navigation-drawer--rail :deep(.v-list-group__items .v-list-item) {
  padding-inline-start: 8px !important;
}

.nav-btns {
  margin: 0 10px;
}

.v-navigation-drawer--mini-variant .nav-btns {
  max-width: 40px;
}
</style>
