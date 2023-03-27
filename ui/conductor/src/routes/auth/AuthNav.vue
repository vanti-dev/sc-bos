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
import {featureEnabled} from '@/routes/config';
import {onMounted, ref} from 'vue';

const menuItems = [
  {
    title: 'Users',
    icon: 'mdi-account-group',
    link: {path: '/auth/users'}
  },
  {
    title: 'Third-party Accounts',
    icon: 'mdi-key',
    link: {path: '/auth/third-party'}
  }
];

// computed props shouldn't return a promise, so instead we're setting this ref based on mounted
const enabledMenuItems = ref([]);
onMounted(async () => {
  // create array of true/false vals for whether each menu item is enabled
  const isEnabled = await Promise.all(menuItems.map(item => featureEnabled(item.link.path)));
  // filter menu items based on above list
  enabledMenuItems.value = menuItems.filter((item, index) => isEnabled[index]);
});

</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
