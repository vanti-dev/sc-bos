import {defineStore} from 'pinia';
import {newActionTracker} from '@/api/resource';
import {computed, reactive} from 'vue';
import {timestampToDate} from '@/api/convpb';
import {listTenants} from '@/api/ui/tenant';

export const useTenantStore = defineStore('tenantStore', () => {
  const tenantsTracker = reactive(
      /** @type {ActionTracker<ListTenantsResponse.AsObject>} */ newActionTracker()
  );

  const tenantsList = computed(() => {
    if (!tenantsTracker.response) return [];
    return tenantsTracker.response.tenantsList.map((t) => ({
      ...t,
      createTime: t.createTime ? timestampToDate(t.createTime) : null
    }));
  });

  /**
   */
  async function refreshTenants() {
    await listTenants(null, tenantsTracker);
  }

  return {
    tenantsList,
    tenantsTracker,
    refreshTenants
  };
});
