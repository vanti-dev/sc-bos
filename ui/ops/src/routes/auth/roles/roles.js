import {usePermissionsStore} from '@/stores/permissions.js';
import {computed, toValue} from 'vue';

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