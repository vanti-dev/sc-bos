<template>
  <status-alert v-if="props.streamError" icon="mdi-connection" :resource="props.streamError"/>

  <v-icon
      v-else-if="ok"
      color="white"
      :style="{visibility: showOk ? 'initial' : 'hidden'}"
      size="20">
    mdi-check
  </v-icon>

  <v-menu
      v-else-if="notOk"
      :close-on-content-click="false"
      location="left"
      max-width="500px"
      min-width="500px">
    <template #activator="{props: menuProps}">
      <v-tooltip location="left">
        <template #activator="{props: tooltipProps}">
          <v-icon
              v-bind="{...menuProps, ...tooltipProps}"
              :color="iconColor"
              size="20">
            {{
              iconStr
            }}
          </v-icon>
        </template>
        Status
      </v-tooltip>
    </template>
    <v-card>
      <v-card-title>
        <v-icon start :color="iconColor">{{ iconStr }}</v-icon>
        <span>{{ levelToStr(level) }}</span>
      </v-card-title>
      <v-card-text>
        {{ description }}
      </v-card-text>
      <template v-if="hasMoreProblems">
        <v-card-actions>
          <v-btn variant="text" @click="showMore = !showMore" block>
            Show details
            <v-spacer/>
            <v-icon end>{{ showMore ? 'mdi-chevron-up' : 'mdi-chevron-down' }}</v-icon>
          </v-btn>
        </v-card-actions>
        <v-expand-transition>
          <v-card-text v-if="showMore" class="py-0">
            <v-list>
              <v-list-item v-for="problem in problems" :key="problem.name" lines="three">
                <v-list-item-title>{{ levelToStr(problem.level) }}</v-list-item-title>
                <v-list-item-subtitle>{{ problem.name }}</v-list-item-subtitle>
                <v-list-item-subtitle>{{ problem.description }}</v-list-item-subtitle>
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
import StatusAlert from '@/components/StatusAlert.vue';
import {useStatusLog} from '@/traits/status/status.js';
import {ref} from 'vue';

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
  },
  streamError: {
    type: Object,
    default: null
  }
});

const {
  level, levelToStr,
  ok, notOk,
  description,
  iconStr, iconColor,
  problems, hasMoreProblems
} = useStatusLog(() => props.value);
const showMore = ref(false);
</script>
