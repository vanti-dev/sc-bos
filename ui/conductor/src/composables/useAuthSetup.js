import {computed, reactive, ref, watch} from 'vue';
import {useAccountStore} from '@/stores/account';
import {useAppConfigStore} from '@/stores/app-config';
import {roleToPermissions} from '@/routes/auth/roleToPermissions';
import router from '@/routes/router';

/**
 * Initializing the authentication setup
 *
 * @return {{init: (function(*, *): void), role: ComputedRef<null|string>, accessLevel: ComputedRef<boolean>}}
 */
export default function() {
  const appConfig = useAppConfigStore();
  const accountStore = useAccountStore();

  /**
   * @param {string} to
   * @param {NavigationGuardNext} next
   */
  const init = (to, next) => {
    if (to === '/') {
      next(appConfig.homePath);
    } else {
      next(appConfig.pathEnabled(to));
    }
  };

  // Logged in user's role
  const role = computed(() => accountStore.role);
  // Logged in user's permissions depending on the role
  const rolePermissions = computed(() => roleToPermissions[role.value]);
  // Current route path
  const path = computed(() => router.currentRoute.path);

  // The user should have 3 access levels: fullAccess, limitedAccess, blockedAccess
  // Depending on the role and role permissions, we are going to check if the user has
  // permission (and access) to certain pages and functionalities
  const accessLevel = computed(() => {
    if (!role.value || !rolePermissions.value) {
      return {
        fullAccess: false,
        limitedAccess: false,
        blockedAccess: true
      };
    }

    const fullAccess = rolePermissions.value.fullAccess.includes(path.value);
    const limitedAccess = rolePermissions.value.limitedAccess.includes(path.value);
    const blockedAccess = rolePermissions.value.blockedAccess.includes(path.value);

    return {
      fullAccess,
      limitedAccess,
      blockedAccess
    };
  });

  const blockActions = computed(() => {
    if (role.value === 'viewer') return true;
    return false;
  });

  const blockSystemEdit = computed(() => {
    if (role.value === 'viewer') return true;
    return false;
  });


  return {
    init,

    role,
    accessLevel,

    blockActions,
    blockSystemEdit
  };
}
