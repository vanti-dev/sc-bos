<template>
  <v-app class="app-root">
    <app-bar/>
    <navigation-drawer/>

    <router-view v-if="isLoggedIn"/>
    <error-view/>
  </v-app>
</template>

<script setup>
import AppBar from '@/components/default/AppBar.vue';
import NavigationDrawer from '@/components/default/NavigationDrawer.vue';
import ErrorView from '@/components/ui-error/ErrorView.vue';

import useAuthSetup from '@/composables/useAuthSetup';
import {useAccountStore} from '@/stores/account.js';
import {useControllerStore} from '@/stores/controller';
import {onMounted} from 'vue';


const {isLoggedIn} = useAuthSetup();
const controller = useControllerStore();
const store = useAccountStore();

store.loadLocalStorage();

onMounted(() => {
  controller.sync();
});
</script>

<style lang="scss" scoped>
.v-application {
  background: var(--v-neutral-darken1);
}
</style>
