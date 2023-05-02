<template>
  <v-list class="pa-0" dense nav>
    <v-list-item-group>
      <v-list-group group="zone">
        <template #activator>
          <v-list-item-icon>
            <v-icon>mdi-group</v-icon>
          </v-list-item-icon>
          <v-list-item-content>Zone Editor</v-list-item-content>
        </template>
        <v-list-item-group>
          <v-list-item
              v-for="zone of zoneList"
              :key="zone"
              :to="'/site/zone/'+zone">
            <v-list-item-icon>
              <v-icon>mdi-select-all</v-icon>
            </v-list-item-icon>
            <v-list-item-content>{{ zone }}</v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-list-group>
    </v-list-item-group>
  </v-list>
</template>

<script setup>
import {ServiceNames} from '@/api/ui/services';
import {usePageStore} from '@/stores/page';
import {useServicesStore} from '@/stores/services';
import {storeToRefs} from 'pinia';
import {computed, onUnmounted, ref, watch} from 'vue';

const servicesStore = useServicesStore();
const pageStore = usePageStore();
const {sidebarNode} = storeToRefs(pageStore);
const zoneCollection = ref({});

watch(sidebarNode, async () => {
  zoneCollection.value = servicesStore.getService(
      ServiceNames.Zones,
      await sidebarNode.value.commsAddress,
      await sidebarNode.value.commsName).servicesCollection;

  // todo: this causes us to load all pages, connect with paging logic instead - although we might want it in this case
  zoneCollection.value.needsMorePages = true;
}, {immediate: true});

watch(zoneCollection, () => {
  zoneCollection.value.query(ServiceNames.Zones);
});

onUnmounted(() => zoneCollection.value.reset());

const zoneList = computed(() => {
  return Object.values(zoneCollection.value?.resources?.value ?? []).map(zone => {
    return zone.id;
  }).sort();
});

</script>

<style scoped>
</style>
