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
import {computed, onMounted, ref} from 'vue';
import {featureEnabled} from '@/routes/config';

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
    title: 'Site Settings',
    subtitle: 'Configure site-specific settings and edit zones',
    icon: 'mdi-cog',
    link: {name: 'site'}
  },
  {
    title: 'Workflows & Automations',
    subtitle: 'View automation status and update settings',
    icon: 'mdi-priority-low',
    link: {name: 'automations'}
  }
];

// computed props shouldn't return a promise, so instead we're setting this ref based on mounted
const enabledMenuItems = ref([]);
onMounted(async () => {
  // create array of true/false vals for whether each menu item is enabled
  const isEnabled = await Promise.all(menuItems.map(item => featureEnabled(item.link.name)));
  // filter menu items based on above list
  enabledMenuItems.value = menuItems.filter((item, index) => isEnabled[index]);
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
