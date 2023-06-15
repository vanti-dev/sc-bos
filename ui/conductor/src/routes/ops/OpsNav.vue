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

    <v-list-item v-for="(item, key) in enabledMenuItems" :to="item.link" :key="key" class="my-2">
      <v-list-item-icon>
        <v-badge
            class="font-weight-bold"
            :color="item.count ? 'primary' : 'transparent'"
            :content="counts[item.count]"
            overlap
            :value="counts[item.count]">
          <v-icon>
            {{ item.icon }}
          </v-icon>
        </v-badge>
      </v-list-item-icon>
      <v-list-item-content>
        <v-list-item-title class="text-truncate">{{ item.title }}</v-list-item-title>
      </v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {computed, onMounted, reactive} from 'vue';

import {useAlertMetadata} from '@/routes/ops/notifications/alertMetadata';
import {useAppConfigStore} from '@/stores/app-config';

const alertMetadata = useAlertMetadata();
const appConfig = useAppConfigStore();

const counts = reactive({
  unacknowledgedAlertCount: computed(() => alertMetadata.unacknowledgedAlertCount)
});

const menuItems = [
  {
    title: 'Notifications',
    icon: 'mdi-bell-outline',
    link: {path: '/ops/notifications'},
    count: 'unacknowledgedAlertCount'
  },
  {
    title: 'Emergency Lighting',
    icon: 'mdi-alarm-light-outline',
    link: {path: '/ops/emergency-lighting'}
  },
  {
    title: 'Security',
    icon: 'mdi-shield-key',
    link: {path: '/ops/security'}
  }
];


const enabledMenuItems = computed(() => {
  return menuItems.filter(item => appConfig.pathEnabled(item.link.path));
});
const overviewEnabled = computed(() => appConfig.pathEnabled('/ops/overview'));

onMounted(() => {
  alertMetadata.init();
});
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
