<template>
  <v-list class="pa-0" dense nav>
    <v-list-item-group>
      <v-list-item to="/ops/overview" v-if="overviewEnabled">
        <v-list-item-icon>
          <v-icon>mdi-domain</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>Building Overview</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list-item-group>

    <v-list-item v-for="(item, key) in enabledMenuItems" :to="item.link" :key="key">
      <v-list-item-icon>
        <v-icon>{{ item.icon }}</v-icon>
      </v-list-item-icon>
      <v-list-item-content>
        <v-list-item-title>{{ item.title }}</v-list-item-title>
      </v-list-item-content>

      <v-chip class="font-weight-bold text primary" v-if="item.count">
        {{ item.count }}
      </v-chip>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {useAlertMetadata} from '@/routes/ops/notifications/alertMetadata';
import {onMounted, ref} from 'vue';
import {featureEnabled} from '@/routes/config';

const alertMetadata = useAlertMetadata();

const menuItems = [
  {
    title: 'Notifications',
    icon: 'mdi-bell-outline',
    link: {path: '/ops/notifications'},
    count: ref(alertMetadata.unacknowledgedAlertCount)
  },
  {
    title: 'Emergency Lighting',
    icon: 'mdi-alarm-light-outline',
    link: {path: '/ops/emergency-lighting'}
  }
];

// computed props shouldn't return a promise, so instead we're setting this ref based on mounted
const enabledMenuItems = ref([]);
const overviewEnabled = ref(false);

onMounted(async () => {
  // create array of true/false vals for whether each menu item is enabled
  const isEnabled = await Promise.all(menuItems.map(item => featureEnabled(item.link.path)));
  // filter menu items based on above list
  enabledMenuItems.value = menuItems.filter((item, index) => isEnabled[index]);

  overviewEnabled.value = await featureEnabled('/ops/overview');
});

</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
