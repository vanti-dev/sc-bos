<template>
  <v-list class="pa-0" dense nav>
    <v-list-item-group>
      <v-list-group group="zone">
        <template #activator>
          <v-list-item-icon>
            <v-icon>mdi-group</v-icon>
          </v-list-item-icon>
          <v-list-item-title class="text-truncate">Zone Editor</v-list-item-title>
        </template>
        <v-list-item-group>
          <v-list-item
              v-for="zone of zoneList"
              :key="zone"
              :to="'/site/zone/' + zone"
              class="my-2 text-truncate"
              :disabled="hasNoAccess('/site/zone/' + zone)">
            <v-list-item-icon>
              <v-icon>mdi-select-all</v-icon>
            </v-list-item-icon>
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
