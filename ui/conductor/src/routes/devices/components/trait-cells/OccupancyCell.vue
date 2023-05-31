<template>
  <div>
    <p :class="[props.value.occupancyState.toLowerCase(), 'ma-0 text-body-2']">
      {{ props.value.occupancyState }}
    </p>
    <v-progress-linear color="primary" indeterminate :active="props.value.occupancyValue.loading"/>
  </div>
</template>

<script setup>
const props = defineProps({
  value: {
    type: Object,
    default: () => {}
  }
});
</script>
