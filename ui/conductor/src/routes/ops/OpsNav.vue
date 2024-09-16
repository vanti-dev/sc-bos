<template>
  <v-list class="px-0" density="compact" nav>
    <ops-nav-list-items
        :items="overviewChildren"
        :mini-variant="miniVariant"/>

    <v-divider v-if="overviewChildren.length > 1" class="mb-3 mt-n1"/>
    <!-- Main List -->
    <v-list-item
        v-for="(item, key) in enabledMenuItems"
        :to="item.link"
        :key="key"
        :disabled="hasNoAccess(item.link.path)">
      <template #prepend>
        <v-badge
            class="font-weight-bold"
            :color="item.badgeType ? badges[item.badgeType].color : 'transparent'"
            :content="item.badgeType ? badges[item.badgeType].value : ''"
            :model-value="Boolean(item.badgeType ? badges[item.badgeType].value : null)">
          <v-icon>
            {{ item.icon }}
          </v-icon>
        </v-badge>
      </template>
      <v-list-item-title class="text-truncate">{{ item.title }}</v-list-item-title>
    </v-list-item>
  </v-list>
</template>


<script setup>
import {closeResource} from '@/api/resource.js';
import useAuthSetup from '@/composables/useAuthSetup';
import {useAlertMetadataStore} from '@/routes/ops/notifications/alertMetadata';
import OpsNavListItems from '@/routes/ops/overview/OpsNavListItems.vue';
import {useNavStore} from '@/stores/nav';
import {useUiConfigStore} from '@/stores/ui-config';
import {storeToRefs} from 'pinia';
import {computed, onMounted, onUnmounted, reactive} from 'vue';

const {miniVariant} = storeToRefs(useNavStore());

const {hasNoAccess} = useAuthSetup();
const alertMetadata = useAlertMetadataStore();
const uiConfig = useUiConfigStore();

/**
 * Collect the overview children
 * This is used to create the sub-lists
 * Each child has a list of children
 *
 * @type {
 *  import('vue').ComputedRef<{title: string, icon: string, children: {title: string, shortTitle: string}[]}[]>
 * } overviewChildren
 */
const overviewChildren = computed(function() {
  const pages = uiConfig.getOrDefault('ops.pages');
  if (!pages) {
    const root = {
      title: 'Building Overview',
      icon: 'mdi-domain',
      path: 'building',
      ...uiConfig.getOrDefault('ops.overview')
    };
    return [root];
  }
  return pages;
});


// --------- Helpers --------- //


/**
 * Reactive object containing icon badge types and their values
 *
 * @typedef {import('vue').ComputedRef<{color: string, value: string|number|null}>} unacknowledgedAlertCount
 * @type {{
 * unacknowledgedAlertCount: unacknowledgedAlertCount
 * }} badges
 */
const badges = reactive({
  unacknowledgedAlertCount: computed(() => {
    if (notificationEnabled.value) {
      if (alertMetadata.alertError?.error) {
        return {
          color: 'error',
          value: '!'
        };
      } else {
        return {
          color: 'primary',
          value: alertMetadata.badgeCount
        };
      }
    } else {
      return {
        color: 'transparent',
        value: null
      };
    }
  })
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
    badgeType: 'unacknowledgedAlertCount'
  },
  {
    title: 'Air Quality',
    icon: 'mdi-air-filter',
    link: {path: '/ops/air-quality'},
    badgeType: null
  },
  {
    title: 'Emergency Lighting',
    icon: 'mdi-alarm-light-outline',
    link: {path: '/ops/emergency-lighting'},
    badgeType: null
  },
  {
    title: 'Security',
    icon: 'mdi-shield-key',
    link: {path: '/ops/security'},
    badgeType: null
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
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: rgb(var(--v-theme-primary));
}
</style>
