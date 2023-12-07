<template>
  <v-tooltip
      v-if="props.resource"
      v-model="show"
      bottom
      :color="props.color">
    <template #activator="{ on, attrs }">
      <v-icon
          v-bind="attrs"
          v-on="on"
          :color="props.color"
          size="22"
          style="padding-top: 1px;">
        {{ props.icon }}
      </v-icon>
    </template>
    <div v-for="(status, index) in statusDetails" :key="status.statusName + index">
      <span class="error-name">{{ status.statusName }}</span>
      <span class="error-details">
        {{ status.statusCode }}
        {{ status.statusMessage }}
      </span>
      <v-divider
          v-if="!props.single && index !== statusDetails.length - 1"
          class="neutral lighten--4 my-1 mx-auto"
          style="width:4em"/>
    </div>
  </v-tooltip>
</template>

<script setup>
import {computed, ref} from 'vue';
import {statusCodeToString} from '@/components/ui-error/util';

const props = defineProps({
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
  single: {
    type: Boolean,
    default: true
  }
});

const show = ref(false);

const statusDetails = computed(() => {
  if (!props.single) {
    const errors = props.resource.errors || [];

    return errors.map((status) => {
      return {
        statusCode: statusCodeToString(status?.resource?.error?.code),
        statusMessage: ': ' + status?.resource?.error?.message,
        statusName: status.name
      };
    });
  } else {
    return [{
      statusCode: statusCodeToString(props.resource?.error?.code),
      statusMessage: ': ' + props.resource?.error?.message,
      statusName: props.resource?.name
    }];
  }
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
