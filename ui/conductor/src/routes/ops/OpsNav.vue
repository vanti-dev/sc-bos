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

      <v-chip class="font-weight-bold text primary" v-if="item.count?.value">
        {{ item.count.value }}
      </v-chip>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {useAlertMetadata} from '@/routes/ops/notifications/alertMetadata';
import {useAppConfigStore} from '@/stores/app-config';
import {computed, ref} from 'vue';

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

const appConfig = useAppConfigStore();

const enabledMenuItems = computed(() => {
  return menuItems.filter(item => appConfig.pathEnabled(item.link.path));
});
const overviewEnabled = computed(() => appConfig.pathEnabled('/ops/overview'));


</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
