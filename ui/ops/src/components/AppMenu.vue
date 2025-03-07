<template>
  <v-menu
      v-bind="showMenu"
      location="bottom left"
      transition="slide-x-transition"
      :content-props="{style: 'left: 0'}">
    <template #activator="{ props }">
      <v-btn tile variant="flat" v-bind="props" :ripple="false" id="main-nav-button" color="neutral-lighten-1">
        <menu-icon width="60"/>
      </v-btn>
    </template>
    <v-card max-width="512" rounded="0 be">
      <v-list tile lines="three" class="pt-0">
        <v-list-item
            v-for="(item, key) in enabledMenuItems"
            :to="item.link"
            :disabled="hasNoAccess(item.link.name)"
            :key="key">
          <template #prepend>
            <v-icon size="40px">{{ item.icon }}</v-icon>
          </template>
          <v-list-item-title class="text-h4">{{ item.title }}</v-list-item-title>
          <v-list-item-subtitle class="text-body-small">{{ item.subtitle }}</v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </v-card>
  </v-menu>
</template>

<script setup>
import MenuIcon from '@/components/MenuIcon.vue';
import useAuthSetup from '@/composables/useAuthSetup';
import {useUiConfigStore} from '@/stores/uiConfig.js';
import {computed, ref} from 'vue';

const {hasNoAccess} = useAuthSetup();

const showMenu = ref(false);

const menuItems = [
  {
    title: 'Operations',
    subtitle: 'View status dashboards, check notifications and events',
    icon: 'mdi-bell-ring',
    link: {name: 'ops'}
  },
  {
    title: 'Devices',
    subtitle:
        'View all devices, check status, and control settings',
    icon: 'mdi-devices',
    link: {name: 'devices'}
  },
  {
    title: 'Access Management',
    subtitle: 'Manage user accounts, and create API tokens',
    icon: 'mdi-key',
    link: {name: 'auth'}
  },
  {
    title: 'Workflows & Automations',
    subtitle: 'View automation status and update settings',
    icon: 'mdi-priority-low',
    link: {name: 'automations'}
  },
  {
    title: 'System',
    subtitle: 'Status & settings for the underlying Building Operating System',
    icon: 'mdi-cogs',
    link: {name: 'system'}
  }
];

const uiConfig = useUiConfigStore();

const enabledMenuItems = computed(() => {
  return menuItems.filter((item) => uiConfig.pathEnabled('/' + item.link.name));
});
</script>

<style scoped>
#main-nav-button {
  margin-left: 0;
  height: 60px !important;
  width: 60px;
  min-width: auto;
}

#main-nav-button:focus {
  background-color: rgb(var(--v-theme-neutral)) !important;
  background-blend-mode: normal;
}
</style>
