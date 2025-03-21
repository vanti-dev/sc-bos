import {getRole, listRoles} from '@/api/ui/account.js';
import {useAction} from '@/composables/action.js';
import useCollection from '@/composables/collection.js';
import {usePermissionsStore} from '@/stores/permissions.js';
import {computed, effectScope, onScopeDispose, reactive, toValue, watch} from 'vue';

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
 * Returns a reactive object with a key per id whose values are the result of useGetRole.
 *
 * @param {import('vue').MaybeRefOrGetter<string>} name
 * @param {import('vue').MaybeRefOrGetter<string[]>} ids
 * @return {import('vue').Reactive<Record<string, ToRefs<UnwrapNestedRefs<UseActionResponse<Role.AsObject>>>>>}
 */
export function useGetRoles(name, ids) {
  // keyed by id
  const res = reactive(
      /** @type {Record<string, ToRefs<UnwrapNestedRefs<UseActionResponse<Role.AsObject>>>>} */
      {}
  );
  // keyed by id
  const closers = /** @type {Record<string, function(): void>} */ {};

  const requests = computed(() => {
    const _name = toValue(name);
    // make sure no duplicate ids exist
    const idSet = new Set(toValue(ids));
    return Array.from(idSet).map((id) => ({name: _name, id}));
  });
  watch(requests, (requests) => {
    const toDelete = new Set(Object.keys(res)); // of id strings
    const toAdd = new Set(); // of requests
    for (const request of requests) {
      const key = request.id;
      toDelete.delete(key);
      if (!Object.hasOwn(closers, key)) {
        toAdd.add(request);
      }
    }
    for (const key of toDelete.values()) {
      closers[key]();
      delete res[key];
      delete closers[key];
    }
    for (const request of toAdd.values()) {
      const scope = effectScope();
      closers[request.id] = () => scope.stop();
      scope.run(() => {
        res[request.id] = useGetRole(request);
      });
    }
  });
  onScopeDispose(() => {
    for (const close of Object.values(closers)) {
      close();
    }
  });

  return res;
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