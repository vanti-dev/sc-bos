<template>
  <svg :viewBox="`0 0 ${iconSize} ${iconSize}`" class="icon">
    <template v-for="(_,y) in 3">
      <template v-for="(_,x) in 3">
        <circle :r="iconRadius" :cx="iconPadding + x*iconGap" :cy="iconPadding + y*iconGap"/>
      </template>
    </template>
  </svg>
</template>
<script setup>
// icon props
const iconSize = 60;
const iconRadius = 2;
const iconPadding = 15 + iconRadius;
const iconGap = (iconSize - 2 * iconPadding) / 2;

defineProps({
  color: {
    type: String,
    default: '#fff'
  }
});

</script>
<style scoped>
.icon circle {
  stroke: none;
  fill: v-bind(color);
}
</style>
