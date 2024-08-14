<template>
  <status-alert v-if="props.streamError" icon="mdi-cancel" :resource="props.streamError"/>

  <v-tooltip v-else left>
    <template #activator="{props}">
      <v-icon :class="[grantClass]" right size="20" v-bind="props">mdi-door</v-icon>
    </template>
    <span class="text-capitalize">Access: {{ grantState.split('_').join(' ') }}</span>
  </v-tooltip>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {useAccessAttempt} from '@/traits/access/access.js';

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
  },
  streamError: {
    type: Object,
    default: null
  }
});

const {grantState, grantClass} = useAccessAttempt(() => props.value);
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
