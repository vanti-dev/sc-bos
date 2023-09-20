import {onMounted, onUnmounted, reactive} from 'vue';

import {closeResource, newActionTracker} from '@/api/resource';
import {getEnrollment, testEnrollment} from '@/api/sc/traits/enrollment';
import {useHubStore} from '@/stores/hub';

/**
 *
 */
export default function() {
  const {listHubNodesAction} = useHubStore();
  const enrollmentValue = reactive(newActionTracker());
  const testEnrollmentValue = reactive(newActionTracker());

  const getEnrollmentValue = async () => {
    try {
      await getEnrollment(enrollmentValue);
    } catch (e) {
      console.warn('Error fetching enrollment', e);
    }

    return enrollmentValue;
  };

  // Closing enrollment on unmount
  onUnmounted(() => {
    closeResource(enrollmentValue);
    closeResource(testEnrollmentValue);
  });

  onMounted(async () => {
    await listHubNodesAction();
    await getEnrollmentValue();
    testEnrollment(testEnrollmentValue);
  });

  return {};
}
