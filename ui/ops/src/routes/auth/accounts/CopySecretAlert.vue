<template>
  <v-alert type="info" tile icon="mdi-key" closable @click:close="emit('close')">
    <template #title>
      Service Account Created
    </template>
    <template #text>
      <p class="mb-3">
        Make sure to save your secret token somewhere private. You won't be able to see it again.
      </p>
      <div class="secret d-flex align-center">
        <pre class="ml-4 my-2 mr-auto text-h6">{{ props.credential }}</pre>
        <v-btn :icon="secretCopied ? 'mdi-check' : 'mdi-content-copy'" variant="text" @click="onCopySecret">
          <v-icon size="24"/>
          <v-menu activator="parent"
                  :open-on-click="false"
                  :model-value="secretCopied"
                  location="bottom"
                  offset="6">
            <v-card text="Secret copied to clipboard" color="success"/>
          </v-menu>
        </v-btn>
      </div>
    </template>
  </v-alert>
</template>

<script setup>
import {onScopeDispose, ref, watch} from 'vue';

const props = defineProps({
  credential: {
    type: String,
    default: undefined,
  }
});
const emit = defineEmits(["close"]);

const secretCopied = ref(false);
const onCopySecret = () => {
  navigator.clipboard.writeText(props.credential);
  secretCopied.value = true;
}
let secretCopiedTimeout = 0;
watch(secretCopied, (value) => {
  clearTimeout(secretCopiedTimeout);
  if (!value) return;
  secretCopiedTimeout = setTimeout(() => {
    secretCopied.value = false;
  }, 5000);
});
onScopeDispose(() => clearTimeout(secretCopiedTimeout));
</script>

<style scoped>
.secret {
  border-radius: 0.2rem;
  border: 1px solid rgba(0, 0, 0, 0.12);
  background-color: rgba(0, 0, 0, 0.04);
}
</style>