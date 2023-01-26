<template>
  <v-card>
    <v-card-title>{{ type }}</v-card-title>
  </v-card>
</template>

<script setup>

defineProps({
  type: {
    type: String,
    default: ''
  }
});

</script>

<style scoped>

</style>
