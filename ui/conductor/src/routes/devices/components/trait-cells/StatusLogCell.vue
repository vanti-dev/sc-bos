<template>
  <v-icon v-if="ok" color="white" v-show="showOk">mdi-check</v-icon>
  <v-menu v-else-if="notOK" :close-on-content-click="false" offset-y left max-width="500px" min-width="500px">
    <template #activator="{ on }">
      <v-icon v-on="on" :color="iconColor">{{ iconStr }}</v-icon>
    </template>
    <v-card>
      <v-card-title>
        <v-icon left :color="iconColor">{{ iconStr }}</v-icon>
        <span>{{ levelStr(level) }}</span>
      </v-card-title>
      <v-card-text>
        {{ description }}
      </v-card-text>
      <template v-if="hasMore">
        <v-card-actions>
          <v-btn text @click="showMore = !showMore" block>
            Show details
            <v-spacer/>
            <v-icon right>{{ showMore ? 'mdi-chevron-up' : 'mdi-chevron-down' }}</v-icon>
          </v-btn>
        </v-card-actions>
        <v-expand-transition>
          <v-card-text v-if="showMore" class="py-0">
            <v-list>
              <v-list-item v-for="problem in problems" :key="problem.name" three-line>
                <v-list-item-content>
                  <v-list-item-title>{{ levelStr(problem.level) }}</v-list-item-title>
                  <v-list-item-subtitle>{{ problem.name }}</v-list-item-subtitle>
                  <v-list-item-subtitle>{{ problem.description }}</v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-expand-transition>
      </template>
    </v-card>
  </v-menu>
  <span v-else/>
</template>

<script setup>
import {StatusLog} from '@sc-bos/ui-gen/proto/status_pb';
import {computed, ref} from 'vue';

const props = defineProps({
  value: {
    type: Object,
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  },
  showOk: {
    type: Boolean,
    default: false
  }
});

// note, these aren't mutually exclusive. props.value === null will be false for both for example
const level = computed(() => props.value?.level || 0);
const ok = computed(() => props.value?.level === StatusLog.Level.NOMINAL);
const notOK = computed(() => props.value?.level > StatusLog.Level.NOMINAL);
const description = computed(() => props.value?.description || '');
const iconColor = computed(() => {
  if (level.value <= StatusLog.Level.NOMINAL) return 'white';
  if (level.value <= StatusLog.Level.NOTICE) return 'info';
  if (level.value <= StatusLog.Level.REDUCED_FUNCTION) return 'warning';
  if (level.value <= StatusLog.Level.NON_FUNCTIONAL) return 'error';
  if (level.value <= StatusLog.Level.OFFLINE) return 'grey';
  return 'white';
});
const iconStr = computed(() => {
  if (level.value <= StatusLog.Level.NOMINAL) return 'mdi-check-circle-outline';
  if (level.value <= StatusLog.Level.NOTICE) return 'mdi-information-outline';
  if (level.value <= StatusLog.Level.REDUCED_FUNCTION) return 'mdi-progress-alert';
  if (level.value <= StatusLog.Level.NON_FUNCTIONAL) return 'mdi-alert-circle-outline';
  if (level.value <= StatusLog.Level.OFFLINE) return 'mdi-connection';
  return '';
});
const problems = computed(() => props.value?.problemsList || []);
const hasMore = computed(() => problems.value.length > 0);

const showMore = ref(false);

/**
 * @param {number} level
 * @return {string}
 */
function levelStr(level) {
  if (level === StatusLog.Level.LEVEL_UNDEFINED) return '';
  if (level <= StatusLog.Level.NOMINAL) return 'Nominal';
  if (level <= StatusLog.Level.NOTICE) return 'Notice';
  if (level <= StatusLog.Level.REDUCED_FUNCTION) return 'Reduced Function';
  if (level <= StatusLog.Level.NON_FUNCTIONAL) return 'Non-Functional';
  if (level <= StatusLog.Level.OFFLINE) return 'Offline';
  return 'Custom Level ' + level;
}

</script>
