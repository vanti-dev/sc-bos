<template>
  <v-alert
      density="compact"
      type="warning"
      rounded="xl"
      style="position: absolute; top: 20%; max-width: 415px;"
      :model-value="props.showAlert"
      width="100%">
    {{ props.message }}
  </v-alert>
</template>

<script setup>
const props = defineProps({
  showAlert: {
    type: Boolean,
    default: false
  },
  message: {
    type: String,
    default: ''
  }
});
</script>
