import {computed} from 'vue';
import {useAccountStore} from '@/stores/account';
import {useAppConfigStore} from '@/stores/app-config';
import {roleToPermissions} from '@/assets/roleToPermissions';

/**
 * Initializing the authentication setup
 *
 * @return {{
 * navigate: (function(*, *): void),
 * role: ComputedRef<null|string>,
 * accessLevel: (function(string): boolean),
 * isLoggedIn: ComputedRef<boolean>,
 * blockActions: ComputedRef<boolean>,
 * blockSystemEdit: ComputedRef<boolean>
 * }}
 */
export default function() {
  const appConfig = useAppConfigStore();
  const accountStore = useAccountStore();

  /**
   * @param {string} toPath
   * @param {NavigationGuardNext} next
   */
  const navigate = (toPath, next) => {
    if (toPath === '/') {
      next(appConfig.homePath);
    } else {
      next(appConfig.pathEnabled(toPath));
    }
  };

  // Logged in user's roles
  // replace - in a role to be camelCase
  const roles = computed(() => {
    return accountStore.roles.map((role) => role.replace(/-([a-z])/g, (g) => g[1].toUpperCase()));
  });

  // Logged in user's permissions depending on the role
  const rolePermissions = computed(() => {
    return roles.value.map((role) => ({
      role,
      permissions: roleToPermissions[role] || {fullAccess: [], limitedAccess: [], blockedAccess: []}
    }));
  });

  // The user should have 3 access levels: fullAccess, limitedAccess, blockedAccess
  // Depending on the role and role permissions, we are going to check if the user has
  // permission (and access) to certain pages and functionalities
  const accessLevel = (name) => {
    if (!appConfig.config?.disableAuthentication) {
      // Formatting the name to match the main path (e.g. /site or /devices)
      const formattedName = !name.includes('/') ? `/${name}` : name;

      // If the roles or role permissions are not defined, the user has blocked access
      if (roles.value.length === 0 || !rolePermissions.value.length) {
        return {
          fullAccess: false,
          limitedAccess: false,
          blockedAccess: true
        };
      }

      // Match the main path (e.g. /site or /devices) with the role permissions
      const access = rolePermissions.value.map((rp) => ({
        role: rp.role,
        fullAccess: rp.permissions.fullAccess.includes(formattedName),
        limitedAccess: rp.permissions.limitedAccess.includes(formattedName),
        blockedAccess: rp.permissions.blockedAccess.includes(formattedName)
      }));

      return access;
    } else {
      return {
        fullAccess: true,
        limitedAccess: false,
        blockedAccess: false
      };
    }
  };

  const hasAnyRole = (...targetRoles) => {
    return targetRoles.some(role => roles.value.includes(role));
  };

  // Checking if the user has no access to certain pages and functionalities
  // depending on multiple roles and role permissions
  const hasNoAccess = (name) => {
    const accessLevels = accessLevel(name);

    if (Array.isArray(accessLevels)) {
      return accessLevels.some(access => access.blockedAccess);
    }

    return accessLevels.blockedAccess;
  };

  const allowActions = computed(() => hasAnyRole('admin', 'superAdmin', 'commissioner', 'operator'));
  const allowEdits = computed(() => hasAnyRole('admin', 'superAdmin', 'commissioner'));

  // Blocking actions (e.g. edit, delete, light control etc.)
  const blockActions = computed(() => {
    if (!appConfig.config?.disableAuthentication) {
      return !allowActions.value;
    } else return false;
  });

  // Blocking system edit (e.g. add, edit, delete, restart etc.)
  const blockSystemEdit = computed(() => {
    if (!appConfig.config?.disableAuthentication) {
      return !allowEdits.value;
    } else {
      return false;
    }
  });


  return {
    navigate,

    roles,
    accessLevel,
    hasNoAccess,
    isLoggedIn: computed(() => {
      if (!appConfig.config?.disableAuthentication) {
        return accountStore.isLoggedIn;
      } else return true;
    }),

    blockActions,
    blockSystemEdit
  };
}
