<template>
  <v-list class="pa-0" dense nav>
    <v-list-item-group>
      <v-list-item to="/ops/overview">
        <v-list-item-icon>
          <v-icon>mdi-domain</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>Building Overview</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list-item-group>

    <v-list-item to="/ops/notifications">
      <v-list-item-icon>
        <v-icon>mdi-bell-outline</v-icon>
      </v-list-item-icon>
      <v-list-item-content>
        <v-list-item-title>Notifications</v-list-item-title>
      </v-list-item-content>

      <v-chip class="font-weight-bold text primary" v-if="alertMetadata.unacknowledgedAlertCount">
        {{ alertMetadata.unacknowledgedAlertCount }}
      </v-chip>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {useAlertMetadata} from '@/routes/ops/notifications/alertMetadata';

const alertMetadata = useAlertMetadata();
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
