<template>
  <v-app class="app-root">
    <app-bar/>
    <navigation-drawer/>

    <router-view/>
    <error-view/>
  </v-app>
</template>

<script setup>
import AppBar from '@/components/default/AppBar.vue';
import NavigationDrawer from '@/components/default/NavigationDrawer.vue';
import ErrorView from '@/components/ui-error/ErrorView.vue';

import {onMounted} from 'vue';
import {useControllerStore} from '@/stores/controller';

const controller = useControllerStore();
onMounted(() => {
  controller.sync();
});

</script>

<style lang="scss" scoped>
.v-application {
  background: var(--v-neutral-darken1);
}
</style>
