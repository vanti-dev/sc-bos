<template>
  <div>
    <div class="d-flex flex-row align-center" :style="setLeftMargin">
      <v-list-item
          :active-class="!isDeepestActiveItem ? 'primary--text black' : ''"
          :class="[
            'd-flex flex-row align-center my-0 mr-2',
            {'mb-2': !isOpen && !isDeepestActiveItem, 'mr-9': !hasChildren && !props.miniVariant}
          ]"
          :to="toAreaLink"
          @click="setActiveOverview(item)">
        <v-list-item-icon>
          <v-icon v-if="!props.miniVariant || !props.item.shortTitle">{{ props.item.icon }}</v-icon>
          <v-list-item-title v-else class="text-center text-truncate" style="max-width: 24px;">
            {{ props.item.shortTitle }}
          </v-list-item-title>
        </v-list-item-icon>
        <v-list-item-title>{{ props.item.title }}</v-list-item-title>
      </v-list-item>
      <v-btn
          v-show="hasChildren && !props.miniVariant"
          class="ma-0 pa-0"
          icon
          small
          @click="toggle">
        <v-icon>{{ isOpen ? 'mdi-chevron-down' : 'mdi-chevron-left' }}</v-icon>
      </v-btn>
    </div>
    <v-slide-y-transition hide-on-leave>
      <OpsNavList
          v-if="isOpen && hasChildren"
          :class="isOpen ? 'mb-n1' : ''"
          :items="props.item.children"
          :depth="props.depth + 0.5"
          :mini-variant="props.miniVariant"
          :parent-path="currentPath"/>
    </v-slide-y-transition>
  </div>
</template>

<script setup>
import {computed, ref} from 'vue';
import {storeToRefs} from 'pinia';
import {useRoute} from 'vue-router/composables';
import {useOverviewStore} from '@/routes/ops/overview/overviewStore';

import OpsNavList from '@/routes/ops/overview/OpsNavList.vue';

const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  depth: {
    type: Number,
    default: 0
  },
  miniVariant: {
    type: Boolean,
    default: false
  },
  parentPath: {
    type: String,
    default: ''
  }
});
const {activeOverview} = storeToRefs(useOverviewStore());
const route = useRoute();
const isOpen = ref(false);

/**
 * Computed checker to see if the item has children
 *
 * @type {import('vue').ComputedRef<boolean>}
 */
const hasChildren = computed(() => props.item.children && props.item.children.length > 0);

/**
 * Computed checker to return the current path
 *
 * @type {import('vue').ComputedRef<string>}
 */
const currentPath = computed(() => {
  const pathSegments = props.parentPath ? [props.parentPath] : [];
  pathSegments.push(encodeURIComponent(props.item.title));
  return pathSegments.join('/');
});

/**
 * Computed checker to return the link to the area
 *
 * @type {import('vue').ComputedRef<string>}
 */
const toAreaLink = computed(() => `/ops/overview/building/${currentPath.value}`);

/**
 * Computed checker to see if the current path is the deepest active item
 *
 * @type {import('vue').ComputedRef<boolean>}
 */
const isDeepestActiveItem = computed(() => {
  const currentPathSegments = route.path.split('/').filter(segment => segment);
  const lastSegment = currentPathSegments[currentPathSegments.length - 1];
  return lastSegment === encodeURIComponent(props.item.title);
});

/**
 * Computed value to set the left margin
 *
 * @type {import('vue').ComputedRef<{marginLeft: string}>}
 */
const setLeftMargin = computed(() => {
  const baseMargin = 1;
  const increment = 7;
  const marginLeft = `${baseMargin + (props.depth * increment)}px`;
  return {marginLeft};
});

/**
 * Toggle the open state
 *
 * @return {void}
 */
const toggle = () => {
  isOpen.value = !isOpen.value;
};

/**
 * Set the active overview
 *
 * @param {Object} item
 * @return {void}
 */
const setActiveOverview = (item) => {
  // Destructure the item to separate the 'children' property and the rest of the properties
  // eslint-disable-next-line no-unused-vars
  const {children, ...rest} = item;

  // Set activeOverview with the rest of the properties, excluding 'children'
  activeOverview.value = rest;
};
</script>
