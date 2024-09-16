import {useUiConfigStore} from '@/stores/uiConfig.js';
import {subPath} from '@/util/path.js';

/**
 * @return {{
 *   toPath: (path: string) => string
 * }}
 */
export function usePathUtils() {
  const uiConfig = useUiConfigStore();
  return {
    toPath(path) {
      if (!path) return path;
      return subPath(path, uiConfig.configUrl);
    }
  };
}
