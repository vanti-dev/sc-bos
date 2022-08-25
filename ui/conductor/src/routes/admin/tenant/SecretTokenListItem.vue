<template>
  <v-list-item>
    <v-list-item-icon class="mr-4">
      <v-icon color="success">mdi-check</v-icon>
    </v-list-item-icon>
    <v-list-item-content>
      <v-list-item-title class="d-flex align-baseline">
        <span>{{ secret.token }}</span>
        <v-btn @click="copyToClipboard" icon v-if="clipboardCopySupported" small class="ms-2" title="Copy to clipboard">
          <v-icon small v-if="clipboardCopyState === 'wait'">mdi-content-copy</v-icon>
          <v-icon small v-else-if="clipboardCopyState === 'ok'" color="success">mdi-check</v-icon>
          <v-icon small v-else color="error">mdi-alert-circle-outline</v-icon>
        </v-btn>
      </v-list-item-title>
    </v-list-item-content>
    <div>
      <v-btn outlined @click="emit('hideToken')">Done</v-btn>
      <v-btn color="error" outlined class="ml-4" @click="emit('delete', secret)">Delete</v-btn>
    </div>
  </v-list-item>
</template>

<script setup>
import {milliseconds} from 'date-fns';
import {onMounted, ref} from 'vue';

const props = defineProps({
  secret: Object
});
const emit = defineEmits(['hideToken', 'delete']);

const clipboardCopySupported = ref(false);
onMounted(() => {
  clipboardCopySupported.value = Boolean(navigator?.clipboard?.writeText);
})

const clipboardCopyState = ref('wait');
let clipboardCopyAgainHandle = 0;

function copyToClipboard() {
  navigator.clipboard.writeText(props.secret.token)
      .then(() => clipboardCopyState.value = 'ok')
      .then(() => {
        clearTimeout(clipboardCopyAgainHandle);
        clipboardCopyAgainHandle = setTimeout(() => {
          clipboardCopyState.value = 'wait';
        }, milliseconds({seconds: 10}))
      })
      .catch(err => clipboardCopyState.value = err);
}
</script>

<style scoped>
.v-list-item:before {
  content: '';
  position: absolute;
  inset: 0;
  background-color: var(--v-success-base);
  opacity: .2;
}

.v-list-item:hover:before {
  opacity: .2;
}
</style>
