import {defineStore} from 'pinia';
import {onMounted, onUnmounted, reactive} from 'vue';

import {getEnrollment, testEnrollment} from '@/api/sc/traits/enrollment';
import {closeResource, newActionTracker} from '@/api/resource';

export const useEnrollmentStore = defineStore('enrollment', () => {
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

  // Getting enrollment on mount
  onMounted(async () => {
    await getEnrollmentValue();
    testEnrollment(testEnrollmentValue);
  });

  // Closing enrollment on unmount
  onUnmounted(() => {
    closeResource(enrollmentValue);
    closeResource(testEnrollmentValue);
  });

  return {
    enrollmentValue,
    testEnrollmentValue,

    getEnrollmentValue
  };
});
