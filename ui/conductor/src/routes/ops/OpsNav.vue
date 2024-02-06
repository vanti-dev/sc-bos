<template>
  <v-list class="pa-0" dense nav>
    <v-list-item-group class="mt-2 mb-n1">
      <span v-if="uiConfig.pathEnabled('/ops/overview')" class="d-flex flex-row align-center ma-0">
        <v-list-item class="mb-0" :disabled="hasNoAccess('/ops/overview/building')" to="/ops/overview/building">
          <v-list-item-icon>
            <v-icon>mdi-domain</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Building Overview</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-tooltip v-if="!miniVariant" bottom>
          <template #activator="{ on }">
            <v-btn
                class="ma-0 pa-0 ml-2"
                :disabled="!buildingChildren.length"
                icon
                small
                v-on="on"
                @click="displayList = !displayList">
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
        <OpsNavList
            v-if="displayList"
            :display-list="displayList"
            :items="buildingChildren"
            :mini-variant="miniVariant"/>
      </v-slide-y-transition>

      <v-divider v-if="displayList" class="mb-3 mt-n1"/>
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
              :color="item.countType && counts[item.countType] ? 'primary' : 'transparent'"
              :content="counts[item.countType]"
              overlap
              :value="counts[item.countType]">
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
import useAuthSetup from '@/composables/useAuthSetup';
import {useAlertMetadata} from '@/routes/ops/notifications/alertMetadata';
import OpsNavList from '@/routes/ops/overview/OpsNavList.vue';
import {usePageStore} from '@/stores/page';
import {useUiConfigStore} from '@/stores/ui-config';
import {storeToRefs} from 'pinia';
import {computed, onMounted, reactive, ref} from 'vue';

const pageStore = usePageStore();
const {miniVariant} = storeToRefs(pageStore);

const {hasNoAccess} = useAuthSetup();
const alertMetadata = useAlertMetadata();
const uiConfig = useUiConfigStore();
const {config} = storeToRefs(uiConfig);


const displayList = ref(false);

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


// --------- Helpers --------- //


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
 *  import('vue').ComputedRef<{title: string, icon: string, link: {path: string}, countType?: string}[]>
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
const enabledMenuItems = computed(() => menuItems.value.filter((item) => uiConfig.pathEnabled(item.link.path)));

/**
 * Check if the notification area is enabled
 *
 * @type {import('vue').ComputedRef<boolean>}
 */
const notificationEnabled = computed(() => uiConfig.pathEnabled('/ops/notifications'));

onMounted(() => {
  if (notificationEnabled.value) alertMetadata.init();
  displayList.value = false;
});
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
