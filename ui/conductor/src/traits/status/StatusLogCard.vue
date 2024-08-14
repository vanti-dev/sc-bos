<template>
  <v-card elevation="0" tile>
    <v-card-title class="text-title-caps-large text-neutral-lighten-3">
      <span>{{ levelToStr(level, 'Status: ') }}</span>
      <v-spacer/>
      <v-icon right :color="iconColor">{{ iconStr }}</v-icon>
    </v-card-title>
    <v-card-text v-if="description">
      {{ description }}
    </v-card-text>
    <template v-if="hasMoreProblems">
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
                <v-list-item-title>{{ levelToStr(problem.level, 'Status: ') }}</v-list-item-title>
                <v-list-item-subtitle>{{ problem.name }}</v-list-item-subtitle>
                <v-list-item-subtitle>{{ problem.description }}</v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-expand-transition>
    </template>
  </v-card>
</template>

<script setup>
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
  }
});

const {
  level, levelToStr,
  description,
  iconStr, iconColor,
  problems, hasMoreProblems
} = useStatusLog(() => props.value);
const showMore = ref(false);
</script>
