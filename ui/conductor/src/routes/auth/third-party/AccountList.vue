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
      <template #item.zones="{ index, value }">
        <span class="d-inline-flex justify-start" style="gap: 8px">
          <v-chip v-for="zone in value" :key="index + zone" small outlined>{{ zone }}</v-chip>
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
                  Add Account<v-icon right>mdi-plus</v-icon>
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
import {onMounted, onUnmounted, ref} from 'vue';
import NewAccountDialog from '@/routes/auth/third-party/components/NewAccountDialog.vue';
import {usePageStore} from '@/stores/page';
import {useTenantStore} from '@/routes/auth/third-party/tenantStore';
import {storeToRefs} from 'pinia';
import {useErrorStore} from '@/components/ui-error/error';
import useAuthSetup from '@/composables/useAuthSetup';

const pageStore = usePageStore();
const tenantStore = useTenantStore();
const {tenantsList, tenantsTracker} = storeToRefs(tenantStore);
const errorStore = useErrorStore();

const search = ref('');

const headers = [
  {text: 'Name', value: 'title'},
  // {text: 'Permissions', value: 'permissions'},
  {text: 'Zones', value: 'zones'}
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
  pageStore.showSidebar = true;
  pageStore.sidebarTitle = item.title;
  pageStore.sidebarData = item;
}

/**
 * @param {*} item
 * @return {string}
 */
function rowClass(item) {
  if (pageStore.showSidebar && pageStore.sidebarData?.id === item.id) {
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
