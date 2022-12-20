<template>
  <v-card elevation="0" class="pa-3" v-bind="$attrs">
    <slot/>
  </v-card>
</template>

<script>
export default {
  inheritAttrs: false
};
</script>
