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
              :key="zone">
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
import {useServicesStore} from '@/stores/services';
import {computed, onMounted, onUnmounted, ref} from 'vue';
import {ServiceNames} from '@/api/ui/services';

const serviceStore = useServicesStore();
const zonesCollection = ref(serviceStore.getService(ServiceNames.Zones).servicesCollection);

// todo: this causes us to load all pages, connect with paging logic instead
zonesCollection.value.needsMorePages = true;

onMounted(() => zonesCollection.value.query(ServiceNames.Zones));
onUnmounted(() => zonesCollection.value.reset());

const zoneList = computed(() => {
  console.log(zonesCollection.value);
  return Object.values(zonesCollection.value.resources.value).map(zone => {
    return zone.id;
  });
});

</script>

<style scoped>
</style>
