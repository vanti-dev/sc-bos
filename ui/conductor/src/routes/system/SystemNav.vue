<template>
  <v-list class="pa-0" dense nav>
    <v-list-item
        v-for="(item, key) in enabledMenuItems"
        :to="item.link"
        :key="key"
        class="my-2"
        :disabled="hasNoAccess(item.link.path)">
      <v-list-item-icon>
        <v-icon>{{ item.icon }}</v-icon>
      </v-list-item-icon>
      <v-list-item-content class="text-truncate">{{ item.title }}</v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script setup>
import useAuthSetup from '@/composables/useAuthSetup';
import {useUiConfigStore} from '@/stores/ui-config';
import {computed} from 'vue';

const {hasNoAccess} = useAuthSetup();

const menuItems = [
  {
    title: 'Drivers',
    icon: 'mdi-arrow-left-right-bold',
    link: {path: '/system/drivers'}
  },
  {
    title: 'Features',
    icon: 'mdi-tools',
    link: {path: '/system/features'}
  },
  {
    title: 'Components',
    icon: 'mdi-memory',
    link: {path: '/system/components'}
  }
];

const uiConfig = useUiConfigStore();

const enabledMenuItems = computed(() => {
  return menuItems.filter((item) => uiConfig.pathEnabled(item.link.path));
});
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
