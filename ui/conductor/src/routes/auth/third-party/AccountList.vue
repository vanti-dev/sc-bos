<template>
  <content-card>
    <v-data-table
        class="table"
        :headers="headers"
        :items="tenantsList"
        :search="search"
        sort-by="title"
        :header-props="{ sortIcon: 'mdi-arrow-up-drop-circle-outline' }"
        :loading="tenantsTracker.loading"
        :item-class="rowClass"
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
                  outlined
                  hide-details
                  prepend-inner-icon="mdi-magnify"
                  v-model="search"/>
            </v-col>
            <v-spacer/>
            <new-account-dialog @finished="tenantStore.refreshTenants">
              <template #activator="{ on, attrs }">
                <v-btn
                    outlined
                    v-bind="attrs"
                    v-on="on"
                    :disabled="blockActions">
                  Add Account
                  <v-icon right>mdi-plus</v-icon>
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
  {text: 'Name', value: 'title', width: '30%'},
  {text: 'Client ID', value: 'id', width: '28em'},
  // {text: 'Permissions', value: 'permissions'},
  {text: 'Zones', value: 'zoneNamesList'}
];

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
 * @param {Tenant.AsObject} item
 */
function showTenant(item) {
  // router.push(`/auth/third-party/${item.id}`);
  sidebar.showSidebar = true;
  sidebar.sidebarTitle = item.title;
  sidebar.sidebarData = item;
}

// update the sidebar data if the tenant list is updated
watch(tenantsList, () => {
  const tenant = tenantsList.value.find(tenant => tenant.id === sidebar.sidebarData.id);
  if (!tenant) {
    return;
  }
  sidebar.sidebarTitle = tenant.title;
  sidebar.sidebarData = tenant;
}, {deep: true});

/**
 * @param {*} item
 * @return {string}
 */
function rowClass(item) {
  if (sidebar.showSidebar && sidebar.sidebarData?.id === item.id) {
    return 'item-selected';
  }
  return '';
}

// ------------------------------ //
// ----- Authentication settings ----- //

const {blockActions} = useAuthSetup();
</script>

<style lang="scss" scoped>
:deep(.v-data-table-header__icon) {
  margin-left: 8px;
}

.table :deep(tbody tr) {
  cursor: pointer;
}

.v-data-table :deep(.v-data-footer) {
  background: var(--v-neutral-lighten1) !important;
  border-radius: 0px 0px $border-radius-root * 2 $border-radius-root * 2;
  border: none;
  margin: 0 -12px -12px;
}

.v-data-table :deep(.item-selected) {
  background-color: var(--v-primary-darken4);
}
</style>
