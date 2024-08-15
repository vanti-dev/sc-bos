<template>
  <div class="d-flex align-center">
    <v-combobox
        v-model="inputModel"
        @update:model-value="sendZoneEvent($event)"
        :items="inputItems"
        item-title="title"
        item-value="name"
        :loading="findZonesLoading"
        :search.sync="searchText"
        :message="findZonesError || []"
        no-filter
        chips
        closable-chips
        density="compact"
        hide-details="auto"
        :no-data-text="`No ${zonesOnly ? 'zones' : 'devices'} found match your query`"
        variant="outlined"
        multiple
        auto-select-first
        return-object
        :disabled="blockActions">
      <template #prepend-item v-if="findZonesTracker.response?.nextPageToken">
        <v-list-subheader class="mx-2">
          <template v-if="findZonesTracker.response?.totalSize > 0">
            Showing {{ inputItems.length }} of {{ findZonesTracker.response.totalSize }}
            {{ zonesOnly ? 'zones' : 'devices' }}
            search to refine your results.
          </template>
          <template v-else>
            Showing up to {{ inputItems.length }} {{ zonesOnly ? 'zones' : 'devices' }}, search to refine your results.
          </template>
        </v-list-subheader>
        <v-divider class="my-2"/>
      </template>

      <template #item="{ item }">
        <div class="d-flex flex-row flex-wrap">
          <v-list-item-title>{{ item.title }}</v-list-item-title>
          <v-list-item-subtitle v-if="item.title !== item.name">{{ item.name }}</v-list-item-subtitle>
        </div>
      </template>

      <template #append-item v-if="findZonesTracker.response?.nextPageToken">
        <!-- New element to handle intersection -->
        <v-divider class="my-2"/>
        <v-btn
            block
            class="text-center mt-1 rounded-0"
            color="transparent"
            :disabled="findZonesLoading"
            elevation="0"
            :loading="findZonesLoading"
            @click="fetchNextPage">
          Load more
        </v-btn>
      </template>
    </v-combobox>
    <v-menu location="bottom" v-if="!blockActions">
      <template #activator="{props}">
        <v-btn rounded="circle" v-bind="props" class="ml-2 mr-n2">
          <v-icon>mdi-dots-vertical</v-icon>
        </v-btn>
      </template>
      <v-card min-width="300">
        <v-list>
          <v-list-item>
            <v-list-item-title>Only show zones</v-list-item-title>
            <v-list-item-action>
              <v-switch v-model="zonesOnly"/>
            </v-list-item-action>
          </v-list-item>
        </v-list>
      </v-card>
    </v-menu>
  </div>
</template>

<script setup>
import {newActionTracker} from '@/api/resource';
import {listDevices} from '@/api/ui/devices';
import useAuthSetup from '@/composables/useAuthSetup';
import debounce from 'debounce';
import {computed, reactive, ref, watch} from 'vue';

const props = defineProps({
  zones: {
    type: Array, // string[]
    default: () => ([])
  },
  maxResultSize: {
    type: Number,
    default: 5
  }
});
const emit = defineEmits(['update:zones', 'update:maxResultSize']);

const propZones = computed(() => deviceArrToItems(props.zones));
const inputZones = ref([]);
const inputItems = computed(() => {
  return deviceArrToItems(findZonesTracker.response?.devicesList ?? []);
});
const inputModel = computed({
  get() {
    if (inputZones.value.length > 0) {
      return inputZones.value;
    } else {
      return propZones.value;
    }
  },
  set(v) {
    inputZones.value = v;
  }
});
const sendZoneEvent = (event) => {
  const zones = [];
  for (const zone of event) {
    if (typeof zone === 'string') {
      zones.push({name: zone}); // manually entered zone
    } else {
      zones.push(zone); // zone selected from search results
    }
  }
  emit('update:zones', zones);
};

const deviceArrToItems = (devices) => devices.map(device => {
  if (typeof device === 'string') {
    device = {name: device};
  }
  return {
    name: device.name,
    title: deviceTitle(device),
    src: device
  };
});
const deviceTitle = (device) => {
  let title = device.metadata?.appearance?.title;
  if (!title) {
    title = device.name.split('/').pop();
  }
  return title;
};

const searchText = ref('');
const zonesOnly = ref(true);

// server size zone vars

// tracks the request to fetch zones; loading, error, etc
const findZonesTracker = reactive(
    /** @type {ActionTracker<ListDevicesResponse.AsObject>} */
    newActionTracker()
);
const findZonesNextPageTracker = reactive(
    /** @type {ActionTracker<ListDevicesResponse>} */
    newActionTracker()
);
const findZonesLoading = computed(() => findZonesTracker.loading || findZonesNextPageTracker.loading);
const findZonesError = computed(() => findZonesTracker.error ?? findZonesNextPageTracker.error);

// the query we use to filter the zones returned by the server
const findZonesQuery = computed(() => {
  const q = {conditionsList: []};
  if (zonesOnly.value) {
    q.conditionsList.push({field: 'metadata.membership.subsystem', stringEqualFold: 'zones'});
  }
  if (searchText.value) {
    const words = searchText.value.split(/\s+/);
    q.conditionsList.push(...words.map(word => ({stringContainsFold: word})));
  }
  return q;
});
// do the fetch of zones, debounced to avoid spamming the server
const fetchZones = debounce((query) => {
  listDevices({query, pageSize: props.maxResultSize}, findZonesTracker)
      .catch(() => {
      }); // errors are recorded in findZonesTracker
}, 500);
// watch for changes in the query and fetch zones when it changes
watch(findZonesQuery, (query) => fetchZones(query), {immediate: true});

const nextPageToken = computed(() => findZonesTracker.response?.nextPageToken);
const fetchNextPage = async () => {
  const pageToken = nextPageToken.value;
  if (pageToken) {
    const req = {query: findZonesQuery.value, pageSize: props.maxResultSize, pageToken};
    try {
      const res = await listDevices(req, findZonesNextPageTracker);
      findZonesTracker.response.devicesList.push(...res.devicesList);
      findZonesTracker.response.nextPageToken = res.nextPageToken;
    } catch (error) {
      console.error('An error occurred while loading more items:', error);
    }
  }
};

const {blockActions} = useAuthSetup();
</script>

<style scoped>

</style>
