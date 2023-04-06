<template>
  <v-menu v-bind="showMenu" bottom offset-y content-class="main-nav" tile transition="slide-x-transition">
    <template #activator="{on, attrs}">
      <v-btn
          tile
          v-bind="attrs"
          v-on="on"
          :ripple="false"
          id="main-nav-button"
          color="neutral lighten-1">
        <menu-icon width="60"/>
      </v-btn>
    </template>
    <v-card max-width="512">
      <v-list tile three-line subheader class="ma-0" color="neutral lighten-1">
        <v-list-item v-for="(item, key) in enabledMenuItems" :to="item.link" :key="key">
          <v-list-item-icon>
            <v-icon x-large>{{ item.icon }}</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title class="text-h4">{{ item.title }}</v-list-item-title>
            <v-list-item-subtitle class="text-body-small">{{ item.subtitle }}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-card>
  </v-menu>
</template>

<script setup>
import MenuIcon from '@/components/MenuIcon.vue';
import {computed, ref} from 'vue';
import {useAppConfigStore} from '@/stores/app-config';

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
    subtitle: 'Add/update/delete devices from the system, view device\'s status and configuration, ' +
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

const appConfig = useAppConfigStore();

const enabledMenuItems = computed(() => {
  return menuItems.filter(item => appConfig.pathEnabled('/'+item.link.name));
});

</script>

<style scoped>
#main-nav-button {
  margin-top: -12px;
  margin-bottom: -12px;
  margin-left: -16px !important;
  height: 60px !important;
  width: 60px;
}

#main-nav-button:focus {
  background-color: var(--v-neutral-base) !important;
  background-blend-mode: normal;
}

.main-nav {
  margin-left: -12px;
}
</style>
