<template>
  <status-alert v-if="props.streamError" icon="mdi-cancel" :resource="props.streamError"/>

  <v-tooltip v-else location="left">
    <template #activator="{props: _props}">
      <v-icon :class="[grantClass]" end size="20" v-bind="_props">mdi-door</v-icon>
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
  color: rgb(var(--v-theme-success));
}

.denied {
  color: rgb(var(--v-theme-warning));
}

.tailgate, .forced, .failed {
  color: rgb(var(--v-theme-error));
}

.pending, .aborted {
  color: rgb(var(--v-theme-info));
}

.grant_unknown {
  color: grey;
}
</style>
