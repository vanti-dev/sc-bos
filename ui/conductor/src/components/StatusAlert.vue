<template>
  <v-tooltip
      v-if="props.resource.streamError"
      v-model="show"
      bottom
      color="error">
    <template #activator="{ on, attrs }">
      <v-icon
          color="error"
          size="20"
          v-bind="attrs"
          v-on="on">
        mdi-alert-circle-outline
      </v-icon>
    </template>
    <span class="error-name">{{ errorDetails.errorName }}</span>
    <span class="error-details">{{ errorDetails.errorCode }}: {{ errorDetails.errorMessage }}</span>
  </v-tooltip>
</template>

<script setup>
import {computed, ref} from 'vue';
import {statusCodeToString} from '@/components/ui-error/util';

const props = defineProps({
  resource: {
    type: Object,
    default: () => {}
  }
});

const show = ref(false);

const errorDetails = computed(() => {
  return {
    errorCode: statusCodeToString(props.resource.streamError?.error?.code),
    errorMessage: props.resource.streamError?.error?.message,
    errorName: props.resource.streamError?.name
  };
});
</script>

<style lang="scss" scoped>
.v-tooltip__content.menuable__content__active {
  padding: 2px 8px;
  opacity: 1 !important; // reduce tooltip transparency for readability
}

.error-name {
  display: block;
  font-size: .8em;
}
.error-details {
  display: block;
  font-size: .9em;
}
</style>
