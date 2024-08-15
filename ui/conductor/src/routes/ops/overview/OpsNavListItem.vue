<template>
  <div>
    <div class="d-flex flex-row align-center" :style="setLeftMargin">
      <v-list-item
          :active-class="!isDeepestActiveItem ? 'primary--text black' : ''"
          :class="[
            'd-flex flex-row align-center mb-0',
            {
              'mr-11': !hasChildren && !props.miniVariant,
              'my-1': props.depth > 0
            }
          ]"
          :to="toAreaLink">
        <template #prepend>
          <v-icon v-if="!props.miniVariant || !props.item.shortTitle">{{ props.item.icon }}</v-icon>
          <v-list-item-title v-else class="text-center text-truncate" style="max-width: 24px;">
            {{ props.item.shortTitle }}
          </v-list-item-title>
        </template>
        <v-list-item-title>{{ props.item.title }}</v-list-item-title>
      </v-list-item>
      <v-btn
          v-show="hasChildren && !props.miniVariant"
          class="ml-2"
          rounded="circle"
          @click="toggle">
        <v-icon>{{ isOpen ? 'mdi-chevron-down' : 'mdi-chevron-left' }}</v-icon>
      </v-btn>
    </div>
    <v-slide-y-transition hide-on-leave>
      <ops-nav-list
          v-if="isOpen && hasChildren"
          :items="props.item.children"
          :depth="props.depth + 1"
          :mini-variant="props.miniVariant"
          :parent-path="currentPath"/>
    </v-slide-y-transition>
  </div>
</template>

<script setup>
import OpsNavList from '@/routes/ops/overview/OpsNavList.vue';
import {computed, ref} from 'vue';
import {useRoute} from 'vue-router';

const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  items: {
    type: Array,
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
  pathSegments.push(encodeURIComponent(props.item.path ?? props.item.title));
  return pathSegments.join('/');
});

/**
 * Computed checker to return the link to the area
 *
 * @type {import('vue').ComputedRef<string>}
 */
const toAreaLink = computed(() => `/ops/overview/${currentPath.value}`);

/**
 * Computed checker to see if the current path is the deepest active item
 *
 * @type {import('vue').ComputedRef<boolean>}
 */
const isDeepestActiveItem = computed(() => {
  const currentPathSegments = route.path.split('/').filter(segment => segment);
  const lastSegment = currentPathSegments[currentPathSegments.length - 1];
  return lastSegment === encodeURIComponent(props.item.path ?? props.item.title);
});

/**
 * Computed value to set the left margin
 *
 * @type {import('vue').ComputedRef<{marginLeft?: string}>}
 */
const setLeftMargin = computed(() => {
  if (props.miniVariant) {
    return {};
  }
  const baseMargin = 0;
  const increment = 8;
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
</script>
