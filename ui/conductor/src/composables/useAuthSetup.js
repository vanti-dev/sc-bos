import {computed} from 'vue';
import {useAccountStore} from '@/stores/account';
import {useAppConfigStore} from '@/stores/app-config';
import {roleToPermissions} from '@/routes/auth/roleToPermissions';

/**
 * Initializing the authentication setup
 *
 * @return {{
 * init: (function(*, *): void),
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
  const init = (toPath, next) => {
    if (toPath === '/') {
      next(appConfig.homePath);
    } else {
      next(appConfig.pathEnabled(toPath));
    }
  };

  // Logged in user's role
  const role = computed(() => accountStore.role);
  // Logged in user's permissions depending on the role
  const rolePermissions = computed(() => roleToPermissions[role.value]);

  // The user should have 3 access levels: fullAccess, limitedAccess, blockedAccess
  // Depending on the role and role permissions, we are going to check if the user has
  // permission (and access) to certain pages and functionalities
  const accessLevel = (name) => {
    if (!appConfig.config?.disableAuthentication) {
      // Formatting the name to match the main path (e.g. /site or /devices)
      const formattedName = !name.includes('/') ? `/${name}` : name;

      // If the role or role permissions are not defined, the user has blocked access
      if (!role.value || !rolePermissions.value) {
        return {
          fullAccess: false,
          limitedAccess: false,
          blockedAccess: true
        };
      }

      // Match the main path (e.g. /site or /devices) with the role permissions
      const fullAccess = rolePermissions.value.fullAccess.includes(formattedName);
      const limitedAccess = rolePermissions.value.limitedAccess.includes(formattedName);
      const blockedAccess = rolePermissions.value.blockedAccess.includes(formattedName);

      return {
        fullAccess,
        limitedAccess,
        blockedAccess
      };
    } else {
      return {
        fullAccess: true,
        limitedAccess: false,
        blockedAccess: false
      };
    }
  };

  // Blocking actions (e.g. edit, delete, light control etc.)
  const blockActions = computed(() => {
    if (!appConfig.config?.disableAuthentication) {
      if (role.value === 'viewer') return true;
      return false;
    } else return false;
  });

  // Blocking system edit (e.g. add, edit, delete, restart etc.)
  const blockSystemEdit = computed(() => {
    if (!appConfig.config?.disableAuthentication) {
      if (role.value === 'viewer' || role.value === 'operator') return true;
      return false;
    } else return false;
  });

  return {
    init,

    role,
    accessLevel,
    isLoggedIn: computed(() => {
      if (!appConfig.config?.disableAuthentication) {
        return accountStore.isLoggedIn;
      } else return true;
    }),

    blockActions,
    blockSystemEdit
  };
}
