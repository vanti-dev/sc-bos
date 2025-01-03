import {timestampToDate} from '@/api/convpb';
import {newActionTracker} from '@/api/resource';
import {listTenants} from '@/api/ui/tenant';
import {defineStore} from 'pinia';
import {computed, reactive} from 'vue';

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
