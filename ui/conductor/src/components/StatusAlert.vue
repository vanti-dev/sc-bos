<template>
  <v-tooltip
      v-if="props.resource"
      v-model="show"
      bottom
      :color="props.color">
    <template #activator="{ on, attrs }">
      <v-btn
          v-bind="attrs"
          v-on="on"
          class="mb-n1"
          icon
          :loading="props.loading"
          :ripple="props.isClickable"
          x-small
          @click="props.clickAction">
        <v-icon
            :color="props.color"
            size="22">
          {{ props.icon }}
        </v-icon>
      </v-btn>
    </template>
    <span class="error-name">{{ statusDetails.statusName }}</span>
    <span class="error-details">
      {{ statusDetails.statusCode }}
      {{ statusDetails.statusMessage }}
    </span>
  </v-tooltip>
</template>

<script setup>
import {computed, ref} from 'vue';
import {statusCodeToString} from '@/components/ui-error/util';

const props = defineProps({
  clickAction: {
    type: Function,
    default: () => null
  },
  color: {
    type: String,
    default: 'error'
  },
  icon: {
    type: String,
    default: 'mdi-alert-circle-outline'
  },
  isClickable: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  resource: {
    type: Object,
    default: () => null
  },
  type: {
    type: String,
    default: 'error'
  }
});

const show = ref(false);

const statusDetails = computed(() => {
  if (['error', 'warning'].includes(props.type)) {
    return {
      statusCode: statusCodeToString(props.resource?.error?.code),
      statusMessage: ': ' + props.resource?.error?.message,
      statusName: props.resource?.name
    };
  } else if (props.type === 'success') {
    return {
      statusCode: statusCodeToString(props.resource?.status?.code),
      statusMessage: ': ' + props.resource?.status?.message,
      statusName: props.resource?.name
    };
  }

  return {
    statusCode: 'Unknown',
    statusMessage: '',
    statusName: ''
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
