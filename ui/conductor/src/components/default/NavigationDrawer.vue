<template>
  <v-navigation-drawer
      v-if="hasNav && isLoggedIn"
      v-model="drawer"
      app
      class="siteNavigation pt-4"
      clipped
      color="transparent"
      :expand-on-hover="!pinDrawer"
      floating
      :mini-variant.sync="miniVariant"
      :mini-variant-width="drawerWidth"
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
            text
            block
            @click="pinDrawer = !pinDrawer">
          <v-icon left>mdi-pin-outline</v-icon>
          {{ !pinDrawer ? 'Pin navigation' : 'Unpin navigation' }}
        </v-btn>
      </div>
      <div class="text-caption text-center neutral--text text--lighten-2 text-no-wrap">
        <v-divider/>
        <p class="mt-2 mb-0">
          Smart Core &copy; {{ new Date().getFullYear() }}
        </p>
        <p :title="appVersion" class="mt-0" style="cursor: default">
          {{ appVersion.split('-')[0] }}
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

.nav-btns {
  margin: 0 10px;
}

.v-navigation-drawer--mini-variant .nav-btns {
  max-width: 40px;
}
</style>
