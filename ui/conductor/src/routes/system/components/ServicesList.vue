<template>
  <content-card>
    <v-row class="pa-4" v-if="configStore.config?.hub">
      <v-combobox
          v-model="userConfig.node"
          :items="nodesListValues"
          label="System Component"
          item-title="name"
          item-value="name"
          hide-details="auto"
          :loading="cohort.loading"
          variant="outlined"/>
      <v-spacer/>
    </v-row>
    <v-data-table
        :headers="headers"
        :items="serviceList"
        item-key="id"
        :search="search"
        :loading="loading"
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
import {useCohortStore} from '@/stores/cohort.js';
import {useUserConfig} from '@/stores/userConfig.js';
import {useUiConfigStore} from '@/stores/ui-config';
import {computed} from 'vue';

const {blockActions} = useAuthSetup();

const configStore = useUiConfigStore();

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
  search,
  serviceList,
  loading,
  nodesListValues,
  showService,
  _startService,
  _stopService
} = useServices(props);
const userConfig = useUserConfig();
const cohort = useCohortStore();

const headers = computed(() => {
  if (props.name === 'drivers') {
    return [
      {title: 'ID', key: 'id'},
      {title: 'Type', key: 'type'},
      {title: 'Status', key: 'active', width: '20em'},
      {key: 'actions', align: 'end', width: '100', sortable: false}
    ];
  } else {
    return [
      {title: 'ID', key: 'id'},
      {title: 'Status', key: 'active', width: '20em'},
      {key: 'actions', align: 'end', width: '100', sortable: false}
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
  background: rgb(var(--v-theme-neutral-lighten-1)) !important;
  border-radius: 0 0 settings.$border-radius-root*2 settings.$border-radius-root*2;
  border: none;
  margin: 0 -12px -12px;
}

.v-data-table :deep(.item-selected) {
  background-color: rgb(var(--v-theme-primary-darken-4));
}


.v-data-table :deep(tr:hover) {
  .automation-device__btn {
    &--red {
      background-color: rgb(var(--v-theme-error));

      .v-btn__content {
        color: white;
      }

      &.v-btn--disabled {
        filter: grayscale(100%);
      }
    }

    &--green {
      background-color: rgb(var(--v-theme-success));

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
