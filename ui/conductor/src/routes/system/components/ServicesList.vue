<template>
  <content-card>
    <v-row class="pa-4" v-if="configStore.config?.hub">
      <v-combobox
          v-model="servicesStore.node"
          :items="nodesListValues"
          label="System Component"
          item-title="name"
          item-value="name"
          hide-details="auto"
          :loading="hubStore.nodesListCollection.loading ?? true"
          variant="outlined"/>
      <v-spacer/>
    </v-row>
    <v-data-table
        :headers="headers"
        :items="serviceList"
        item-key="id"
        :search="search"
        :loading="serviceCollection.loading"
        @click:row="(_, s) => showService(s.item)">
      <template #item.active="{item}">
        <service-status :service="item"/>
      </template>
      <template #item.actions="{item}">
        <v-btn
            v-if="item.active"
            variant="outlined"
            class="automation-device__btn--red"
            color="error"
            :disabled="blockActions"
            width="100%"
            @click.stop="_stopService(item)">
          Stop
        </v-btn>
        <v-btn
            v-else
            variant="outlined"
            class="automation-device__btn--green"
            color="success"
            :disabled="blockActions"
            width="100%"
            @click.stop="_startService(item)">
          Start
        </v-btn>
      </template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import {ServiceNames} from '@/api/ui/services';
import ContentCard from '@/components/ContentCard.vue';
import useAuthSetup from '@/composables/useAuthSetup';
import useServices from '@/composables/useServices';
import ServiceStatus from '@/routes/system/components/ServiceStatus.vue';
import {useHubStore} from '@/stores/hub';
import {useServicesStore} from '@/stores/services.js';
import {useUiConfigStore} from '@/stores/ui-config';
import {computed} from 'vue';

const {blockActions} = useAuthSetup();

const configStore = useUiConfigStore();
const hubStore = useHubStore();

const props = defineProps({
  name: {
    type: String,
    default: ServiceNames.Systems
  },
  // optional type filter for services list
  type: {
    type: String,
    default: ''
  }
});

const {
  serviceCollection,
  search,
  serviceList,
  nodesListValues,
  showService,
  _startService,
  _stopService
} = useServices(props);
const servicesStore = useServicesStore();

const headers = computed(() => {
  if (props.name === 'drivers') {
    return [
      {text: 'ID', value: 'id'},
      {text: 'Type', value: 'type'},
      {text: 'Status', value: 'active', width: '20em'},
      {text: '', value: 'actions', align: 'end', width: '100'}
    ];
  } else {
    return [
      {text: 'ID', value: 'id'},
      {text: 'Status', value: 'active', width: '20em'},
      {text: '', value: 'actions', align: 'end', width: '100'}
    ];
  }
});
</script>

<style lang="scss" scoped>
@use 'vuetify/settings';

:deep(.v-data-table-header__icon) {
  margin-left: 8px;
}

.v-data-table :deep(.v-data-footer) {
  background: var(--v-neutral-lighten1) !important;
  border-radius: 0 0 settings.$border-radius-root*2 settings.$border-radius-root*2;
  border: none;
  margin: 0 -12px -12px;
}

.v-data-table :deep(.item-selected) {
  background-color: var(--v-primary-darken4);
}


.v-data-table :deep(tr:hover) {
  .automation-device__btn {
    &--red {
      background-color: var(--v-error-base);

      .v-btn__content {
        color: white;
      }

      &.v-btn--disabled {
        filter: grayscale(100%);
      }
    }

    &--green {
      background-color: var(--v-success-base);

      .v-btn__content {
        color: white;
      }

      &.v-btn--disabled {
        filter: grayscale(100%);
      }
    }
  }
}
</style>
