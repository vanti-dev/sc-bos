<template>
  <v-alert type="info" tile icon="mdi-key" closable @click:close="emit('close')">
    <template #title>
      Service Account Secret Created
    </template>
    <template #text>
      <p class="mb-3">
        Make sure to save your secret token somewhere private. You won't be able to see it again.
      </p>
      <copy-div class="secret" copied-text="Secret copied to clipboard">
        <pre class="ml-4 my-2 text-h6">{{ props.credential }}</pre>
      </copy-div>
    </template>
  </v-alert>
</template>

<script setup>

import CopyDiv from '@/components/CopyDiv.vue';

const props = defineProps({
  credential: {
    type: String,
    default: undefined,
  }
});
const emit = defineEmits(['close']);
</script>

<style scoped>
.secret {
  border-radius: 0.2rem;
  border: 1px solid rgba(0, 0, 0, 0.12);
  background-color: rgba(0, 0, 0, 0.04);
}
</style>