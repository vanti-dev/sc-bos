<template>
  <floor-list class="container" :floors="props.floors">
    <template #floor="{floor}">
      <div class="floor-item">
        <div class="floor-title">
          {{ floor.title }}
        </div>
      </div>
    </template>
  </floor-list>
</template>

<script setup>
/**
 * @typedef {Object} Floor
 * @property {string} title - display name of the floor
 * @property {number} level - 0 for ground, negative for basements (counting down), positive for upper floors (counting up). Defaults to len - i - 1.
 */

import FloorList from '@/components/FloorList.vue';

const props = defineProps({
  floors: {
    type: Array, // of Floor
    required: true
  },
  selectedFloor: {
    type: Number,
    default: null
  }
});
</script>

<style scoped>
.container {
  height: 100%;
}

.floor-item {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: rgb(var(--v-theme-surface));
  border-radius: 4px;
}
</style>