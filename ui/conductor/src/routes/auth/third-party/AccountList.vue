<template>
  <v-container fluid class="pa-0">
    <main-card>
      <v-card-actions>
        <v-text-field
            label="Search tenants"
            outlined
            hide-details
            prepend-inner-icon="mdi-magnify"
            v-model="search"/>
      </v-card-actions>
      <v-data-table
          class="table"
          :headers="headers"
          :items="tenantRows"
          :search="search"
          sort-by="title"
          show-select
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
      </v-data-table>
    </main-card>
  </v-container>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb.js';
import {newActionTracker} from '@/api/resource.js';
import {listTenants} from '@/api/ui/tenant.js';
import MainCard from '@/components/ContentCard.vue';
import {computed, onMounted, reactive, ref} from 'vue';
import {useRouter} from 'vue-router/composables';

const tenantsTracker = reactive(
    /** @type {ActionTracker<ListTenantsResponse.AsObject>} */ newActionTracker()
);

const search = ref('');

const headers = computed(() => {
  return [
    {text: 'Name', value: 'title'},
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

onMounted(() => {
  listTenants(null, tenantsTracker);
});

const router = useRouter();
/**
 *
 * @param {Tenant.AsObject} item
 */
function showTenant(item) {
  router.push(`/auth/third-party/${item.id}`);
}
</script>

<style scoped>
.table {
  background-color: transparent;
}

::v-deep(.v-data-table-header__icon) {
  margin-left: 8px;
}

.table ::v-deep(tbody tr) {
  cursor: pointer;
}

</style>
