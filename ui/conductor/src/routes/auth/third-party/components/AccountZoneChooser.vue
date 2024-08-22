<template>
  <div class="d-flex align-center">
    <v-combobox
        v-model="inputModel"
        @update:model-value="sendZoneEvent($event)"
        :items="inputItems"
        item-title="title"
        item-value="name"
        :loading="findZonesLoading"
        @update:search="searchText"
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
        :disabled="blockActions"
        :menu="true">
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

      <template #item="{ item, props: _props }">
        <v-list-item v-bind="_props" :subtitle="item.props.title !== item.props.value ? item.props.value : undefined"/>
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
    <v-menu location="bottom" v-if="!blockActions" :close-on-content-click="false">
      <template #activator="{props}">
        <v-btn icon="true" variant="text" v-bind="props" size="small" class="ml-2 mr-n2">
          <v-icon size="24">mdi-dots-vertical</v-icon>
        </v-btn>
      </template>
      <v-card min-width="300">
        <v-list density="compact">
          <v-list-item @click="zonesOnly = !zonesOnly">
            <v-list-item-title>Only show zones</v-list-item-title>
            <template #append>
              <v-list-item-action end>
                <v-switch v-model="zonesOnly" hide-details/>
              </v-list-item-action>
            </template>
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

const _zones = defineModel('zones', {
  type: Array,
  default: () => []
});
const _maxResultSize = defineModel('maxResultSize', {
  type: Number,
  default: 5
});

const propZones = computed(() => deviceArrToItems(_zones.value));
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
  _zones.value = zones;
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
  listDevices({query, pageSize: _maxResultSize.value}, findZonesTracker)
      .catch(() => {
      }); // errors are recorded in findZonesTracker
}, 500);
// watch for changes in the query and fetch zones when it changes
watch(findZonesQuery, (query) => fetchZones(query), {immediate: true});

const nextPageToken = computed(() => findZonesTracker.response?.nextPageToken);
const fetchNextPage = async () => {
  const pageToken = nextPageToken.value;
  if (pageToken) {
    const req = {query: findZonesQuery.value, pageSize: _maxResultSize.value, pageToken};
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
