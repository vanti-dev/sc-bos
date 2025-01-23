<template>
  <div>
    <h1 class="text-h1 py-12 mt-n12">Panel Setup</h1>
    <v-card class="px-4 py-4" color="rgba(255,255,255,0.3)" :loading="zoneMetadataLoading">
      <v-list v-if="zoneList.length > 0" max-height="63vh" class="overflow-auto bg-transparent">
        <v-list-item
            v-for="item in zoneList"
            :key="item.id"
            rounded="xl"
            class="bg-primary mb-2"
            @click="submit(item)">
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item>
        <div v-intersect="handleIntersect"/>
      </v-list>
      <template v-else>
        <v-card-title class="justify-center">No zone available</v-card-title>
        <div v-intersect="handleIntersect"/>
      </template>
    </v-card>
    <v-btn v-if="!disableAuthentication" block class="mt-12" variant="text" @click="logout">Logout</v-btn>
  </div>
</template>

<script setup>
import useMetadata from '@/composables/useMetadata';
import useZoneCollection from '@/composables/useZoneCollection';
import {useAccountStore} from '@/stores/account';
import {useConfigStore} from '@/stores/config';
import {useUiConfigStore} from '@/stores/ui-config';
import {storeToRefs} from 'pinia';
import {computed, ref, watch, watchEffect} from 'vue';
import {useRouter} from 'vue-router';

const emits = defineEmits(['shouldAutoLogout']);

const router = useRouter();
const {zoneCollection, loadNextPage} = useZoneCollection();
const {zones, logout, isInitialized} = storeToRefs(useAccountStore());
const uiConfig = useUiConfigStore();
const disableAuthentication = computed(() => uiConfig.auth.disabled);
const configStore = useConfigStore();

const zoneIds = computed(() => {
  // prefer using the account zones over fetching from the server
  if (zones.value.length > 0) return zones.value;
  if (!zoneCollection?.response?.servicesList) {
    return [];
  }
  const zoneIds = zoneCollection.response.servicesList
      .filter(s => s.active)
      .map(s => s.id);
  return [...new Set(zoneIds)];
});
const {loading: zoneMetadataLoading, trackers: zoneMetadata} = useMetadata(zoneIds);

// Combine zone ids with metadata to make the list of zones easier to read - i.e. use titles.
const zoneList = computed(() => {
  return zoneIds.value.map(id => ({
    id,
    metadata: zoneMetadata[id].response,
    title: zoneMetadata[id]?.response?.appearance?.title ?? id
  }));
});

/**
 *
 * @param {{id: string, metadata: Metadata.AsObject}} zone
 */
async function submit(zone) {
  await configStore.setZone(zone.id, zone.metadata);
  emits('shouldAutoLogout', false);
  await router.push({name: 'home'});
}

// Handle loading zones from the server.
// This is based on
// 1. Whether we want to load zones from the server, as opposed to getting from the account
// 2. Whether the list of zones isn't full enough - aka infinite scroll needs another page
const scrolledToBottom = ref(false);
const handleIntersect = (isIntersecting) => {
  scrolledToBottom.value = isIntersecting;
};
const loadMoreServerZones = computed(() => {
  if (!isInitialized.value) return false; // don't do anything until we know about account zones
  if (zones.value.length > 0) return false; // don't load zones from server if we have account zones to use
  return scrolledToBottom.value;
});
watchEffect(() => {
  loadNextPage.value = loadMoreServerZones.value;
});

// Watch for changes to the zoneList, zoneName, and zoneId
// If zoneName or zoneId is not set, emit an event to auto logout the user
// If a matching zone exists, emit an event to auto logout the user
// If no matching zone exists, do not emit an event to auto logout the user
watch(
    [() => zoneList.value, () => configStore.zoneName, () => configStore.zoneId],
    ([zones, zoneName, zoneId]) => {
      // Check if zoneName or zoneId is not set
      if (!zoneName || !zoneId) {
        emits('shouldAutoLogout', true);
        return;
      }

      // Check if a matching zone exists
      const match = zones.some(zone => zone.text === zoneId);

      emits('shouldAutoLogout', match);
    },
    {immediate: true}
);
</script>
