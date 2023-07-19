<template>
  <side-bar>
    <v-subheader v-if="!notificationSidebar.length" class="text-title-caps-large neutral--text text--lighten-3">
      No Past Notifications
    </v-subheader>
    <div v-else>
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">
        Past {{ notificationSidebar.length }} {{
          notificationSidebar.length === 1 ? 'Notification' : 'Notifications'
        }}
      </v-subheader>
      <v-card
          v-for="(data, index) in notificationSidebar"
          :key="index"
          class="mt-4"
          elevation="0">
        <span class="d-flex flex-row flex-nowrap px-4 mb-2">
          <v-icon :class="[data.severity.color, 'mt-n2']" size="22">{{ data.severity.icon }}</v-icon>
          <v-spacer/>
          <v-card-subtitle class="text-caption pa-0 pb-2 grey--text">
            {{ data.created }}
          </v-card-subtitle>
        </span>
        <v-card-subtitle class="ma-0 pa-0 px-4 white--text text-capitalize">
          {{ data.description }}
        </v-card-subtitle>
        <v-divider v-if="index < notificationSidebar.length - 1" class="my-3"/>
      </v-card>
    </div>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';

import {useNotifications} from '@/routes/ops/notifications/notifications.js';
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import {computed} from 'vue';


const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);
const notification = useNotifications();

const notificationSidebar = computed(() => {
  // if sidebarData is not an object, or if it is an object but does not have a value property, return an empty array
  if (!sidebarData || !sidebarData.value || typeof sidebarData.value !== 'object') {
    return []; // or handle the unexpected data structure appropriately
  } else {
    // otherwise, continue with the expected data structure
    const icons = {
      info: 'mdi-information-outline',
      warn: 'mdi-alert-circle-outline',
      alert: 'mdi-alert-box-outline',
      danger: 'mdi-close-octagon'
    };

    // filter out the sidebarData entries that are not 'past10'
    const filteredEntries = Object.entries(sidebarData.value).filter(([key, value]) => {
      return key === 'past10';
    });

    // reduce the filtered entries to an array of objects
    const mergedData = filteredEntries.map(([key, value]) => {
      // reduce the value array to an array of objects
      const reducedValue = [];
      value.forEach(item => {
        // if the item has a resolveTime, set the severity color to grey
        let severityColor = '';
        if (item.resolveTime) {
          severityColor = 'grey--text';
          // otherwise, set the severity color to the severity color
        } else {
          severityColor = notification.severityData(item.severity).color;
        }

        // push the reduced item to the reducedValue array
        reducedValue.push({
          severity: {
            icon: icons[notification.severityData(item.severity).text.toLowerCase()],
            color: severityColor
          },
          description: item.description,
          created: new Date(item.createTime).toLocaleString()
        });
      });

      return {
        value: reducedValue
      };
    });

    // if mergedData is not an array, or if it is an array but does not have a length property, return an empty array
    if (!mergedData || !mergedData.length || typeof mergedData[0] !== 'object') {
      return []; // or handle the unexpected data structure appropriately
    } else {
      return mergedData[0].value;
    }
  }
});
</script>

<style scoped>
</style>
