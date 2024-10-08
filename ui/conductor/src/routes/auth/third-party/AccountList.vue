<template>
  <content-card>
    <v-data-table
        class="table"
        :headers="headers"
        :items="tenantsList"
        :search="search"
        v-model:sort-by="sortBy"
        :loading="tenantsTracker.loading"
        :row-props="rowProps"
        @click:row="showTenant">
      <template #item.zoneNamesList="{ index, value }">
        <span class="d-inline-flex justify-start" style="gap: 8px">
          <name-chip v-for="zone in value" :key="index + zone" small outlined :name="zone"/>
        </span>
      </template>
      <template #top>
        <v-container fluid style="width: 100%">
          <v-row dense align="center">
            <v-col cols="12" md="5">
              <v-text-field
                  label="Search accounts"
                  variant="outlined"
                  hide-details
                  prepend-inner-icon="mdi-magnify"
                  v-model="search"/>
            </v-col>
            <v-spacer/>
            <new-account-dialog @finished="tenantStore.refreshTenants">
              <template #activator="{ props }">
                <v-btn
                    variant="outlined"
                    v-bind="props"
                    :disabled="blockActions">
                  Add Account
                  <v-icon end>mdi-plus</v-icon>
                </v-btn>
              </template>
            </new-account-dialog>
          </v-row>
        </v-container>
      </template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {useErrorStore} from '@/components/ui-error/error';
import useAuthSetup from '@/composables/useAuthSetup';
import NameChip from '@/routes/auth/third-party/components/NameChip.vue';
import NewAccountDialog from '@/routes/auth/third-party/components/NewAccountDialog.vue';
import {useTenantStore} from '@/routes/auth/third-party/tenantStore';
import {useSidebarStore} from '@/stores/sidebar';
import {storeToRefs} from 'pinia';
import {onMounted, onUnmounted, ref, watch} from 'vue';

const sidebar = useSidebarStore();
const tenantStore = useTenantStore();
const {tenantsList, tenantsTracker} = storeToRefs(tenantStore);
const errorStore = useErrorStore();

const search = ref('');

const headers = [
  {title: 'Name', key: 'title', width: '30%'},
  {title: 'Client ID', key: 'id', width: '28em'},
  // {title: 'Permissions', key: 'permissions'},
  {title: 'Zones', key: 'zoneNamesList', sortable: false}
];

const sortBy = ref([{key: 'title', order: 'asc'}]);

// UI error handling
let unwatchErrors;
onMounted(() => {
  unwatchErrors = errorStore.registerTracker(tenantsTracker);
  tenantStore.refreshTenants();
});
onUnmounted(() => {
  if (unwatchErrors) unwatchErrors();
});

/**
 * @param {PointerEvent} e
 * @param {Tenant.AsObject} item
 */
function showTenant(e, {item}) {
  // router.push(`/auth/third-party/${item.id}`);
  sidebar.visible = true;
  sidebar.title = item.title;
  sidebar.data = item;
}

// update the sidebar data if the tenant list is updated
watch(tenantsList, () => {
  const tenant = tenantsList.value.find(tenant => tenant.id === sidebar.data.id);
  if (!tenant) {
    return;
  }
  sidebar.title = tenant.title;
  sidebar.data = tenant;
}, {deep: true});

/**
 * @param {*} item
 * @return {Record<string,any>}
 */
function rowProps({item}) {
  if (sidebar.visible && sidebar.data?.id === item.id) {
    return {class: 'item-selected'};
  }
  return {};
}

// ------------------------------ //
// ----- Authentication settings ----- //

const {blockActions} = useAuthSetup();
</script>

<style lang="scss" scoped>
@use 'vuetify/settings';

.table :deep(tbody tr) {
  cursor: pointer;
}

.v-data-table :deep(.v-data-footer) {
  background: rgb(var(--v-theme-neutral-lighten-1)) !important;
  border-radius: 0 0 settings.$border-radius-root * 2 settings.$border-radius-root * 2;
  border: none;
  margin: 0 -12px -12px;
}

.v-data-table :deep(.item-selected) {
  background-color: rgb(var(--v-theme-primary-darken-4));
}
</style>
