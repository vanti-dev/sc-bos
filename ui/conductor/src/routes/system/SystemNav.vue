<template>
  <v-list class="pa-0" dense nav>
    <v-list-item v-for="(item, key) in enabledMenuItems" :to="item.link" :key="key">
      <v-list-item-icon>
        <v-icon>{{ item.icon }}</v-icon>
      </v-list-item-icon>
      <v-list-item-content>{{ item.title }}</v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {computed} from 'vue';
import {useAppConfigStore} from '@/stores/app-config';

const menuItems = [
  {
    title: 'Drivers',
    icon: 'mdi-memory',
    link: {path: '/system/drivers'}
  },
  {
    title: 'Features',
    icon: 'mdi-tools',
    link: {path: '/system/features'}
  }
];

const appConfig = useAppConfigStore();

const enabledMenuItems = computed(() => {
  return menuItems.filter(item => appConfig.pathEnabled(item.link.path));
});

</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
