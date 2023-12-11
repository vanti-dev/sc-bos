<template>
  <v-list class="pa-0" dense nav>
    <v-list-item-group class="mt-2 mb-n1">
      <span class="d-flex flex-row align-center ma-0">
        <v-list-item class="mb-0" :disabled="hasNoAccess('/ops/overview')" to="/ops/overview/building">
          <v-list-item-icon>
            <v-icon>mdi-domain</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Building Overview</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-tooltip v-if="!miniVariant" bottom>
          <template #activator="{ on }">
            <v-btn class="ma-0 pa-0 ml-2" :disabled="false" icon small v-on="on" @click="displayList = !displayList">
              <v-icon>
                {{ displayList ? 'mdi-chevron-down' : 'mdi-chevron-left' }}
              </v-icon>
            </v-btn>
          </template>
          <span>{{ displayList ? 'Hide' : 'Show' }} Lists</span>
        </v-tooltip>
      </span>

      <!-- Overview Sub-lists (Areas and Floors) -->
      <v-slide-y-transition>
        <v-list-item-group
            v-if="displayList"
            :class="[
              'pr-10',
              {'ml-4 mr-9': !miniVariant, 'mx-auto': miniVariant, 'mb-n2': !displayList, 'mb-1': displayList}
            ]"
            style="width: 95.5%; height: 100%; max-height: 60vh; overflow-y: auto; overflow-x: visible;">
          <template v-for="(listItem, index) in overviewSubLists">
            <!-- Main List Item (Areas or Floors) -->
            <span class="d-flex flex-row align-center" :key="index">
              <v-list-item
                  :class="['my-2', {'primary darken-4 primary--text': whichRouteActive === listItem.title}]"
                  :input-value="false"
                  @click="toggleLists('sub', listItem.title)">
                <v-list-item-icon>
                  <v-icon :class="{'primary--text': whichRouteActive === listItem.title}">
                    {{ listItem.icon }}
                  </v-icon>
                </v-list-item-icon>
                <v-list-item-title class="text-capitalize">{{ listItem.title }}</v-list-item-title>
              </v-list-item>
              <v-tooltip bottom>
                <template #activator="{ on }">
                  <v-btn
                      class="ma-0 pa-0 ml-2"
                      :disabled="false"
                      icon
                      small
                      v-on="on"
                      @click="displaySubList[listItem.title] = !displaySubList[listItem.title]">
                    <v-icon>
                      {{ displaySubList[listItem.title] ? 'mdi-chevron-down' : 'mdi-chevron-left' }}
                    </v-icon>
                  </v-btn>
                </template>
                <span>{{ displaySubList[listItem.title] ? 'Hide' : 'Show' }} Lists</span>
              </v-tooltip>
            </span>

            <!-- Sub-List Items for Areas or Floors -->
            <v-slide-y-transition v-if="displaySubList[listItem.title]" :key="index">
              <v-list-item-group
                  :class="['mt-n1 mb-2', {'ml-4': !miniVariant, 'mx-auto': miniVariant}]"
                  style="width: 85.5%">
                <!-- Conditionally render Areas or Floors -->

                <!-- Areas List -->
                <template v-if="listItem.title === 'areas'">
                  <template v-for="(subItem, subIndex) in listItem.subListItems">
                    <v-list-item
                        :key="subIndex"
                        :to="createToLinks.toArea(subItem.title)"
                        :class="['my-0 mb-1', {'ml-1': !miniVariant, 'ml-3': miniVariant}]">
                      <v-list-item-icon>
                        <v-icon v-if="!miniVariant || !subItem.shortTitle">mdi-select-group</v-icon>
                        <v-list-item-title v-else class="text-center">{{ subItem.shortTitle }}</v-list-item-title>
                      </v-list-item-icon>
                      <v-list-item-title class="text-capitalize">
                        {{ subItem.title }}
                      </v-list-item-title>
                    </v-list-item>
                  </template>
                </template>

                <!-- Floors List -->
                <template v-else>
                  <template v-for="(floor, floorIndex) in listOfFloors">
                    <span class="d-flex flex-row align-center" :key="floorIndex">
                      <v-list-item :class="['my-0 mb-1', {'ml-1': !miniVariant }]" :to="createToLinks.toFloor(floor)">
                        <v-list-item-icon>
                          <v-icon v-if="!miniVariant">mdi-floor-plan</v-icon>
                          <v-list-item-title v-else class="text-center text-truncate" style="max-width: 25px;">{{
                            floor
                          }}</v-list-item-title>
                        </v-list-item-icon>
                        <v-list-item-title class="text-capitalize">
                          {{ floor }}
                        </v-list-item-title>
                      </v-list-item>
                      <v-tooltip v-if="floorDropdownCollection[floor]" bottom>
                        <template #activator="{ on }">
                          <v-btn
                              class="ma-0 pa-0 ml-2"
                              :disabled="false"
                              icon
                              small
                              v-on="on"
                              @click="toggleLists('floorZone', floor)">
                            <v-icon>
                              {{ displayFloorZoneList[floor] ? 'mdi-chevron-down' : 'mdi-chevron-left' }}
                            </v-icon>
                          </v-btn>
                        </template>
                        <span>{{ displayFloorZoneList[floor] ? 'Hide' : 'Show' }} Lists</span>
                      </v-tooltip>
                    </span>

                    <!-- Nested Zones for each Floor -->
                    <v-slide-y-transition
                        v-if="displayFloorZoneList[floor] && floorDropdownCollection[floor]"
                        :key="floorIndex">
                      <v-list-item-group>
                        <template v-for="(zone, zoneIndex) in floorAreas.find(f => f.floor === floor)?.zones || []">
                          <v-list-item
                              :key="zoneIndex"
                              :to="createToLinks.toZone(floor, zone.title)"
                              :class="[{'ml-3': miniVariant, 'ml-4': !miniVariant}]">
                            <v-list-item-icon>
                              <v-icon v-if="!miniVariant || !zone.shortTitle">mdi-select-group</v-icon>
                              <v-list-item-title v-else class="text-center text-truncate" style="max-width: 25px;">
                                {{
                                  zone.shortTitle
                                }}
                              </v-list-item-title>
                            </v-list-item-icon>
                            <v-list-item-title class="text-capitalize">
                              {{ zone.title }}
                            </v-list-item-title>
                          </v-list-item>
                        </template>
                        <v-divider class="my-1"/>
                      </v-list-item-group>
                    </v-slide-y-transition>
                  </template>
                </template>
              </v-list-item-group>
            </v-slide-y-transition>
            <v-divider v-if="displayList && index < overviewSubLists.length - 1" class="my-0" :key="'divider' + index"/>
          </template>
        </v-list-item-group>
      </v-slide-y-transition>

      <v-divider v-if="displayList" class="mt-n3 mb-3"/>
      <!-- Main List -->
      <v-list-item
          v-for="(item, key) in enabledMenuItems"
          :to="item.link"
          :key="key"
          class="my-2"
          :disabled="hasNoAccess(item.link.path)">
        <v-list-item-icon>
          <v-badge
              class="font-weight-bold"
              :color="item.count && counts[item.count] ? 'primary' : 'transparent'"
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
    </v-list-item-group>
  </v-list>
</template>


<script setup>
import {computed, onMounted, reactive, ref, watch} from 'vue';
import {storeToRefs} from 'pinia';
import {useRouter, useRoute} from 'vue-router/composables';

import {useAlertMetadata} from '@/routes/ops/notifications/alertMetadata';
import {useAppConfigStore} from '@/stores/app-config';
import {usePageStore} from '@/stores/page';
import useAuthSetup from '@/composables/useAuthSetup';
import useFloors from '@/composables/useFloors';

const router = useRouter();
const route = useRoute();

const pageStore = usePageStore();
const {miniVariant} = storeToRefs(pageStore);

const {hasNoAccess} = useAuthSetup();
const alertMetadata = useAlertMetadata();
const appConfig = useAppConfigStore();
const {config} = storeToRefs(appConfig);

const {listOfFloors} = useFloors();


const displayList = ref(false);
const displaySubList = ref({
  areas: false,
  floors: false
});
const displayFloorZoneList = ref({});


/**
 * Collect the building children
 * This is used to create the sub-lists
 * Each child has a list of children
 *
 * @type {
 *  import('vue').ComputedRef<{title: string, icon: string, children: {title: string, shortTitle: string}[]}[]>
 * } buildingChildren
 */
const buildingChildren = computed(() => config.value?.building?.children || []);


/**
 * Collect the floor areas
 * Each floor has a list of zones
 *
 * @type {
 *  import('vue').ComputedRef<{floor: string, zones: {title: string, shortTitle: string}[]}[]>
 * } floorAreas
 */
const floorAreas = computed(() => {
  const floors = buildingChildren.value.find(child => child.title === 'Floors')?.children || [];
  const areas = buildingChildren.value.find(child => child.title === 'Areas')?.children || [];

  return floors.map(floor => ({
    floor: floor.floor,
    zones: areas.filter(area => area.floor === floor.floor)
  }));
});


/**
 * Collect the sub-lists for the overview section
 * Each sub-list has a title, icon, and list of sub-list items
 *
 * @type {
 * import('vue').ComputedRef<{
 *  title: string, icon: string, subListItems: {title: string, shortTitle: string}[], shortTitle: string}[]
 * >} overviewSubLists
 */
const overviewSubLists = computed(() => {
  const floorTitles = new Set(floorAreas.value.flatMap(area => area.zones.map(zone => zone.title)));

  return buildingChildren.value.map(child => {
    const filteredChildren = child.title.toLowerCase() === 'areas' ?
        child.children.filter(area => !floorTitles.has(area.title)) :
        child.children;

    return {
      icon: child.icon,
      subListItems: filteredChildren.map(c => ({title: c.title, shortTitle: c.shortTitle})),
      title: child.title.toLowerCase()
    };
  });
});


/**
 * Determines if a dropdown should be displayed for each floor based on the presence of areas.
 * Returns an object where each key is a floor identifier and the value is a boolean
 * indicating whether the floor has associated areas.
 *
 * @type {
 *  import('vue').ComputedRef<Record<string, boolean>>
 * } floorDropdownCollection
 */
const floorDropdownCollection = computed(() => listOfFloors.value.reduce((acc, floor) => {
  acc[floor] = floorAreas.value.some(area => area.floor === floor);
  return acc;
}, {}));


// --------- Helpers --------- //
/**
 * Create the links for the areas, floors, and zones
 * Making sure the vue-router gets the correct encoded URI
 *
 * @type {import('vue').ComputedRef<{
 *  toArea: (area: string) => string,
 *  toFloor: (floor: string) => string,
 *  toZone: (floor: string, zone: string) => string}
 * >}
 */
const createToLinks = computed(() => ({
  toArea: area => `/ops/overview/building/areas/${encodeURIComponent(area)}`,
  toFloor: floor => `/ops/overview/building/floors/${encodeURIComponent(floor)}`,
  toZone: (floor, zone) =>
    `/ops/overview/building/floors/${encodeURIComponent(floor)}/zones/${encodeURIComponent(zone)}`
}));


/**
 * Function to toggle the lists
 * listType: root, sub, floorZone
 * key: the key of the list to toggle
 *
 * @param {string} listType
 * @param {string} key
 */
const toggleLists = (listType, key) => {
  if (typeof key === 'object') return;

  let listObject = {};

  if (listType === 'sub') {
    listObject = {...displaySubList.value};
  } else if (listType === 'floorZone') {
    listObject = {...displayFloorZoneList.value};
  }

  // Set the value false
  listObject[key] = !listObject[key];

  if (listType === 'sub') {
    displaySubList.value = listObject;
  } else if (listType === 'floorZone') {
    displayFloorZoneList.value = listObject;
  }
};


/**
 * Determine which route is active
 *
 * @type {
 *  import('vue').ComputedRef<string>
 * } whichRouteActive
 */
const whichRouteActive = computed(() => {
  if (route.path.includes('areas')) return 'areas';
  if (route.path.includes('floors')) return 'floors';
  return 'building';
});


/**
 * Notification badge count
 *
 * @type {
 *  import('vue').ComputedRef<number>
 * } counts
 */
const counts = reactive({
  unacknowledgedAlertCount: computed(() => alertMetadata.badgeCount)
});

/**
 * Menu Items
 * This is the main list of items
 *
 * @type {
 *  import('vue').ComputedRef<{title: string, icon: string, link: {path: string}, countType: string}[]>
 * } menuItems
 */
const menuItems = computed(() => [
  {
    title: 'Notifications',
    icon: 'mdi-bell-outline',
    link: {path: '/ops/notifications'},
    countType: 'unacknowledgedAlertCount'
  },
  {
    title: 'Air Quality',
    icon: 'mdi-air-filter',
    link: {path: '/ops/air-quality'}
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
]);

/**
 * Filter the menu items based on the app config (enabled/disabled)
 *
 * @type {
 *  import('vue').ComputedRef<{title: string, icon: string, link: {path: string}, countType: string}[]>
 * } enabledMenuItems
 */
const enabledMenuItems = computed(() => menuItems.value.filter((item) => appConfig.pathEnabled(item.link.path)));

/**
 * If the overview list is open and an area/floor is selected, when the user closes the main list,
 * we redirect them to the building overview page.
 */
watch(displayList, (newVal, oldVal) => {
  if (!newVal && oldVal) {
    displaySubList.value = {areas: false, floors: false};
    Object.keys(displayFloorZoneList.value).forEach(key => displayFloorZoneList.value[key] = false);
    if (route.fullPath !== '/ops/overview/building') {
      router.push({path: '/ops/overview/building'});
    }
  }
}, {deep: true});

onMounted(() => {
  alertMetadata.init();
  displayList.value = false;
  displaySubList.value = {
    areas: false,
    floors: false
  };
});
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
