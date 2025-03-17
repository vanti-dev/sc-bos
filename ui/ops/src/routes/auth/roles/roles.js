import {getRole, listRoles} from '@/api/ui/account.js';
import {useAction} from '@/composables/action.js';
import useCollection from '@/composables/collection.js';
import {usePermissionsStore} from '@/stores/permissions.js';
import {computed, toValue} from 'vue';

/**
 * @param {import('vue').MaybeRefOrGetter<Partial<ListRolesRequest.AsObject>>} request
 * @param {import('vue').MaybeRefOrGetter<Partial<ListOnlyCollectionOptions<ListRolesRequest.AsObject, Role.AsObject>>>?} options
 * @return {UseCollectionResponse<Role.AsObject>}
 */
export function useRolesCollection(request, options) {
  const normOpts = computed(() => {
    return {
      cmp: (a, b) => a.id.localeCompare(b.id, undefined, {numeric: true}),
      ...toValue(options)
    };
  });
  const client = {
    async listFn(req, tracker) {
      const res = await listRoles(req, tracker);
      return {
        items: res.rolesList,
        nextPageToken: res.nextPageToken,
        totalSize: res.totalSize
      };
    },
    pullFn(req, resource) {
      const opts = toValue(normOpts);
      if (opts.pullFn) {
        opts.pullFn(req, resource);
      }
    }
  };

  return useCollection(request, client, normOpts);
}

/**
 * @param {import('vue').MaybeRefOrGetter<Partial<GetRoleRequest.AsObject>>} request
 * @return {ToRefs<UnwrapNestedRefs<UseActionResponse<Role.AsObject>>>}
 */
export function useGetRole(request) {
  return useAction(request, getRole);
}

/**
 * @typedef {Permission.AsObject} AssignedPermission
 * @property {boolean} assigned
 */

/**
 * @param {import('vue').MaybeRefOrGetter<string[]>} perms
 * @return {{
 *   toggleList: ComputedRef<AssignedPermission[]>,
 *   assignedList: ComputedRef<AssignedPermission[]>
 * }}
 */
export function useAssignedPermissions(perms) {
  const permIndex = computed(() => toValue(perms).reduce((acc, perm) => {
    acc[perm] = true;
    return acc;
  }, {}));
  const permissionsStore = usePermissionsStore();
  const toggleList = computed(() => permissionsStore.permissionsList.map((perm) => {
    return {
      ...perm,
      assigned: permIndex.value[perm.id] ?? false,
      implies: permissionsStore.impliedPermissions(perm.id).map((id) => permissionsStore.permissionsById(id)),
      dependsOn: permissionsStore.dependentPermissions(perm.id).map((id) => permissionsStore.permissionsById(id)),
    };
  }));
  const assignedList = computed(() => toggleList.value.filter((perm) => perm.assigned));

  /**
   * @param {string | AssignedPermission} idOrPerm
   * @return {string[]}
   */
  const addPermission = (idOrPerm) => {
    const perm = typeof idOrPerm === 'string' ? permissionsStore.permissionsById(idOrPerm) : idOrPerm;
    if (!perm) {
      return toValue(perms);
    }
    const idx = {...permIndex.value};
    idx[perm.id] = true;
    for (const impliedPermission of permissionsStore.impliedPermissions(perm.id)) {
      idx[impliedPermission] = true;
    }
    return Object.keys(idx).sort();
  }
  const removePermission = (idOrPerm) => {
    const perm = typeof idOrPerm === 'string' ? permissionsStore.permissionsById(idOrPerm) : idOrPerm;
    if (!perm) {
      return toValue(perms);
    }
    const idx = {...permIndex.value};
    delete idx[perm.id];
    for (const dependentPermission of permissionsStore.dependentPermissions(perm.id)) {
      delete idx[dependentPermission];
    }
    return Object.keys(idx).sort();
  }
  return {
    toggleList,
    assignedList,
    addPermission,
    removePermission
  };
}