<template>
  <v-list class="pa-0" dense nav>
    <v-list-item-group>
      <v-list-group group="zone">
        <template #activator>
          <v-list-item>
            <template #prepend>
              <v-icon>mdi-group</v-icon>
            </template>
            <v-list-item-title class="text-truncate">Zone Editor</v-list-item-title>
          </v-list-item>
        </template>
        <v-list-item-group>
          <v-list-item
              v-for="zone of zoneList"
              :key="zone"
              :to="'/site/zone/' + zone"
              class="my-2 text-truncate"
              :disabled="hasNoAccess('/site/zone/' + zone)">
            <template #prepend>
              <v-icon>mdi-select-all</v-icon>
            </template>
            {{ zone }}
          </v-list-item>
        </v-list-item-group>
      </v-list-group>
    </v-list-item-group>
  </v-list>
</template>

<script setup>
import useAuthSetup from '@/composables/useAuthSetup';
import useZones from '@/composables/useZones';

const {hasNoAccess} = useAuthSetup();
const {zoneList} = useZones();

</script>

<style scoped></style>
