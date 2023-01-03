<template>
  <v-card elevation="0" class="pa-3 rounded-lg" v-bind="$attrs">
    <slot/>
  </v-card>
</template>

<script>
export default {
  inheritAttrs: false
};
</script>
