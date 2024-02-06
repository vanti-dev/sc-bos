<template>
  <v-dialog v-model="dialog" max-width="512">
    <v-card class="pa-2">
      <v-card-title class="text-h4 error--text text--lighten px-4">{{ title }}</v-card-title>
      <v-card-text class="px-4">
        <div class="pb-4">
          <slot/>
        </div>
        <v-alert type="error" class="mb-0">
          <slot name="alert-content"/>
        </v-alert>
      </v-card-text>
      <v-progress-linear color="primary" indeterminate :active="progressBar"/>
      <v-card-actions class="justify-end pt-0 px-4">
        <v-btn @click="dialog = false" color="primary">Cancel</v-btn>
        <v-btn @click="confirm" color="error">
          <slot name="confirmBtn"/>
        </v-btn>
      </v-card-actions>
    </v-card>
    <template #activator="attrs">
      <slot name="activator" v-bind="attrs"/>
    </template>
  </v-dialog>
</template>

<script setup>
import {ref} from 'vue';

const emit = defineEmits(['confirm']);
const dialog = ref(false);

defineProps({
  title: {
    type: String,
    default: ''
  },
  confirmBtnText: {
    type: String,
    default: 'Delete'
  },
  progressBar: Boolean
});

/**
 */
function confirm() {
  emit('confirm');
  dialog.value = false;
}

</script>

<style scoped>

</style>
