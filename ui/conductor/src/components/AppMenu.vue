<template>
  <v-menu v-bind="showMenu" location="bottom" offset-y content-class="main-nav" tile transition="slide-x-transition">
    <template #activator="{ props }">
      <v-btn tile v-bind="props" :ripple="false" id="main-nav-button" color="neutral-lighten-1">
        <menu-icon width="60"/>
      </v-btn>
    </template>
    <v-card max-width="512">
      <v-list tile lines="three" subheader class="ma-0" color="neutral-lighten-1">
        <v-list-item
            v-for="(item, key) in enabledMenuItems"
            :to="item.link"
            :disabled="hasNoAccess(item.link.name)"
            :key="key">
          <template #prepend>
            <v-icon size="x-large">{{ item.icon }}</v-icon>
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
import {useUiConfigStore} from '@/stores/ui-config';
import {computed, ref} from 'vue';

const {hasNoAccess} = useAuthSetup();

const showMenu = ref(false);

const menuItems = [
  {
    title: 'Auth',
    subtitle: 'Edit user accounts, and create API tokens',
    icon: 'mdi-key',
    link: {name: 'auth'}
  },
  {
    title: 'Devices',
    subtitle:
        'Add/update/delete devices from the system, view device\'s status and configuration, ' +
        'and control device settings',
    icon: 'mdi-devices',
    link: {name: 'devices'}
  },
  {
    title: 'Operations',
    subtitle: 'View status dashboards, check notifications and events',
    icon: 'mdi-bell-ring',
    link: {name: 'ops'}
  },
  {
    title: 'Workflows & Automations',
    subtitle: 'View automation status and update settings',
    icon: 'mdi-priority-low',
    link: {name: 'automations'}
  },
  {
    title: 'Site Configuration',
    subtitle: 'Configure site-specific settings and edit zones',
    icon: 'mdi-office-building-cog',
    link: {name: 'site'}
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
  margin-top: -12px;
  margin-bottom: -12px;
  margin-left: -16px !important;
  height: 60px !important;
  width: 60px;
  min-width: auto;
}

#main-nav-button:focus {
  background-color: var(--v-neutral-base) !important;
  background-blend-mode: normal;
}

.main-nav {
  margin-left: -12px;
}
</style>
