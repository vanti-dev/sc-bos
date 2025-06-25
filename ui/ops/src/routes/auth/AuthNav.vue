<template>
  <v-list class="pa-0" density="compact" nav>
    <v-list-item
        v-for="(item, key) in enabledMenuItems"
        :to="item.link"
        :key="key"
        class="my-2"
        :disabled="hasNoAccess(item.link.path)"
        :prepend-icon="item.icon">
      <v-list-item-title class="text-truncate">{{ item.title }}</v-list-item-title>
    </v-list-item>
  </v-list>
</template>

<script setup>
import useAuthSetup from '@/composables/useAuthSetup';
import {useUiConfigStore} from '@/stores/uiConfig.js';
import {computed} from 'vue';

const {hasNoAccess} = useAuthSetup();

const menuItems = [
  {
    title: 'Users',
    icon: 'mdi-account-group',
    link: {path: '/auth/users'}
  },
  {
    title: 'Accounts',
    icon: 'mdi-account-box-outline',
    link: {path: '/auth/accounts', name: 'accounts'}
  },
  {
    title: 'Roles',
    icon: 'mdi-account-details',
    link: {path: '/auth/roles', name: 'roles'}
  },
  {
    title: 'Third-party Accounts',
    icon: 'mdi-key',
    link: {path: '/auth/third-party'}
  }
];

const uiConfig = useUiConfigStore();

const enabledMenuItems = computed(() => {
  return menuItems.filter((item) => uiConfig.pathEnabled(item.link.path));
});
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: rgb(var(--v-theme-primary));
}
</style>
