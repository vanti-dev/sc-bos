<template>
  <v-list-item two-line>
    <v-list-item-content>
      <v-list-item-title>
        {{ secret.note }}
        <span v-if="secret.scopeNames" style="opacity: .5" class="font-italic">
          &mdash; {{ secret.scopeNames.join(', ') }}
        </span>
      </v-list-item-title>
      <v-list-item-subtitle v-if="!expires">This secret will not expire</v-list-item-subtitle>
      <v-list-item-subtitle v-else :class="{expired}">
        {{ expired ? 'Expired' : 'Expires' }}
        <relative-date :date="expireTime"/>
      </v-list-item-subtitle>
    </v-list-item-content>
    <div>
      <span v-if="!used">Never used</span>
      <span v-else>Last used <relative-date :date="useTime"/></span>
      <v-btn color="error" outlined class="ml-4">Delete</v-btn>
    </div>
  </v-list-item>
</template>

<script setup>
import RelativeDate from '@/components/RelativeDate.vue';
import {MAX_INT32} from '@/util/number.js';
import {computed, onBeforeUnmount, ref, watch} from 'vue';

const props = defineProps({
  secret: Object,
  showToken: Boolean
});

const expireTime = computed(() => props.secret?.expireTime);
const expires = computed(() => Boolean(expireTime.value));

const expired = ref(false);
let expiredHandle = 0;
watch(expireTime, t => {
  clearTimeout(expiredHandle);
  const delay = t?.getTime() - Date.now();
  if (delay < 0) {
    expired.value = true;
    return;
  }
  expired.value = false;

  // WARNING: Browsers store setTimeout delays as 32-bit integers which means we can't schedule a delay of more than
  // about 25 days.
  if (delay > MAX_INT32) {
    // If the user doesn't interact with their page for more than 25 days then this will break
    return;
  }

  expiredHandle = setTimeout(() => expired.value = true, delay);
}, {immediate: true});
onBeforeUnmount(() => {
  clearTimeout(expiredHandle);
});

const useTime = computed(() => props.secret?.lastUseTime ?? props.secret?.firstUseTime);
const used = computed(() => Boolean(useTime.value));

</script>

<style scoped>
.expired.expired.expired {
  color: var(--v-error-base);
}
</style>
