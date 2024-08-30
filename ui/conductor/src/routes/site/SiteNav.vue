<template>
  <v-list class="pa-0" density="compact" nav>
    <v-list-group>
      <template #activator>
        <v-list-item prepend-icon="mdi-group">
          <v-list-item-title class="text-truncate">Zone Editor</v-list-item-title>
        </v-list-item>
      </template>
      <v-list-item
          v-for="zone of zoneList"
          :key="zone"
          :to="'/site/zone/' + zone"
          class="my-2 text-truncate"
          :disabled="hasNoAccess('/site/zone/' + zone)"
          prepend-icon="mdi-select-all">
        {{ zone }}
      </v-list-item>
    </v-list-group>
  </v-list>
</template>

<script setup>
import useAuthSetup from '@/composables/useAuthSetup';
import useZones from '@/composables/useZones';

const {hasNoAccess} = useAuthSetup();
const {zoneList} = useZones();

</script>

<style scoped></style>
