<template>
  <v-list class="d-flex flex-column pa-0" dense :nav="!miniVariant">
    <ops-nav-list-item
        v-for="(item, index) in props.items"
        :key="index"
        :item="item"
        :items="props.items"
        :depth="depth"
        :mini-variant="props.miniVariant"
        :parent-path="props.parentPath"/>
  </v-list>
</template>

<script setup>
import OpsNavListItem from '@/routes/ops/overview/OpsNavListItem.vue';

// These props are going to be passed in from the child component too
const props = defineProps({
  depth: {
    type: Number,
    required: false,
    default: 0
  },
  items: {
    type: Array,
    required: true
  },
  parentPath: {
    type: String,
    required: false,
    default: ''
  },
  miniVariant: {
    type: Boolean,
    required: false,
    default: true
  }
});
</script>
