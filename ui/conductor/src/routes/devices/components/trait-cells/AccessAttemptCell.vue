<template>
  <v-tooltip left>
    <template #activator="{on}">
      <v-icon :class="[grantStates]" right size="20" v-on="on">mdi-door</v-icon>
    </template>
    <span class="text-capitalize">Access: {{ grantStates.split('_').join(' ') }}</span>
  </v-tooltip>
</template>

<script setup>
import {computed} from 'vue';
import {AccessAttempt} from '@sc-bos/ui-gen/proto/access_pb';

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
  color: green;
}
.denied, .forced, .failed {
  color: red;
}
.pending, .aborted, .tailgate {
  color: orange;
}
.grant_unknown {
  color: grey;
}
</style>
