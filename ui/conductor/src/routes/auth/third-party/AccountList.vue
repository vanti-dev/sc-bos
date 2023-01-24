<template>
  <content-card>
    <v-data-table
        class="table"
        :headers="headers"
        :items="tenantRows"
        :search="search"
        sort-by="title"
        :header-props="{ sortIcon: 'mdi-arrow-up-drop-circle-outline' }"
        :loading="tenantsTracker.loading"
        @click:row="showTenant">
      <template #item.zones="{ index, value }">
        <span class="d-inline-flex justify-start" style="gap: 8px">
          <v-chip v-for="zone in value" :key="index + zone" small outlined>{{
            zone
          }}</v-chip>
        </span>
      </template>
      <template #top>
        <v-container fluid style="width: 100%">
          <v-row dense align="center">
            <v-col cols="12" md="5">
              <v-text-field
                  label="Search tenants"
                  outlined
                  hide-details
                  prepend-inner-icon="mdi-magnify"
                  v-model="search"/>
            </v-col>
            <v-spacer/>
            <new-account-dialog @finished="refreshAccounts">
              <template #activator="{on, attrs}">
                <v-btn outlined v-bind="attrs" v-on="on">Add Account<v-icon right>mdi-plus</v-icon></v-btn>
              </template>
            </new-account-dialog>
          </v-row>
        </v-container>
      </template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb.js';
import {newActionTracker} from '@/api/resource.js';
import {listTenants} from '@/api/ui/tenant.js';
import ContentCard from '@/components/ContentCard.vue';
import {computed, onMounted, reactive, ref} from 'vue';
import NewAccountDialog from '@/routes/auth/third-party/components/NewAccountDialog.vue';
import {usePageStore} from '@/stores/page';

const pageStore = usePageStore();

const tenantsTracker = reactive(
    /** @type {ActionTracker<ListTenantsResponse.AsObject>} */ newActionTracker()
);

const search = ref('');

const headers = computed(() => {
  return [
    {text: 'Name', value: 'title'},
    {text: 'Permissions', value: 'permissions'},
    {text: 'Zones', value: 'zones'}
  ];
});

const tenantRows = computed(() => {
  if (!tenantsTracker.response) return [];
  return tenantsTracker.response.tenantsList.map((t) => ({
    ...t,
    createTime: t.createTime ? timestampToDate(t.createTime) : null
  }));
});

onMounted(() => refreshAccounts());

/**
 */
function refreshAccounts() {
  listTenants(null, tenantsTracker);
}

/**
 *
 * @param {Tenant.AsObject} item
 */
function showTenant(item) {
  // router.push(`/auth/third-party/${item.id}`);
  pageStore.showSidebar = true;
  pageStore.sidebarTitle = item.title;
  pageStore.sidebarData = item;
}
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
  border-radius: 0px 0px $border-radius-root*2 $border-radius-root*2;
  border: none;
  margin: 0 -12px -12px;
}
</style>
