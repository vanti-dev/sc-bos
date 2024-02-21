<template>
  <v-tooltip left>
    <template #activator="{on}">
      <v-icon :class="[grantStates]" right size="20" v-on="on">mdi-door</v-icon>
    </template>
    <span class="text-capitalize">Access: {{ grantStates.split('_').join(' ') }}</span>
  </v-tooltip>
</template>

<script setup>
import {AccessAttempt} from '@sc-bos/ui-gen/proto/access_pb';
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object,
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  },
  showChangeDuration: {
    type: Number,
    default: 30 * 1000
  }
});


const grantId = computed(() => props.value?.grant);
const grantNamesByID = Object.entries(AccessAttempt.Grant).reduce((all, [name, id]) => {
  all[id] = name.toLowerCase();
  return all;
}, {});

const grantStates = computed(() => {
  return grantNamesByID[grantId.value || 0];
});
</script>

<style scoped>
.granted {
  color: var(--v-success-base);
}

.denied {
  color: var(--v-warning-base);
}

.tailgate, .forced, .failed {
  color: var(--v-error-base);
}

.pending, .aborted {
  color: var(--v-info-base);
}

.grant_unknown {
  color: grey;
}
</style>
