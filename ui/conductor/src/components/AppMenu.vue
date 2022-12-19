<template>
  <v-menu
    v-bind="showMenu"
    class="-pl-16"
    transition="slide-x-transition"
    bottom
    right
    offset-y
    max-width="40%"
    tile
  >
    <template #activator="{ on, attrs }">
      <v-btn
        icon
        tile
        v-bind="attrs"
        v-on="on"
        style="background: #ffffff1a"
        :class="btnClass"
      >
        <menu-icon color="white" />
      </v-btn>
    </template>

    <v-card tile width="100vw">
      <nav-button
        :to="{ name: 'auth' }"
        subtitle="Edit user accounts, and create API token"
        title="Auth"
        icon="/menu/auth.svg"
      />
      <nav-button
        :to="{ name: 'devices' }"
        subtitle="Add/update/delete devices from the system, view device's status and configuration, and control device settings"
        title="Devices"
        icon="/menu/devices.svg"
      />
      <nav-button
        :to="{ name: 'ops' }"
        subtitle="View status dashboards,check notifications and events"
        title="Operations"
        icon="/menu/operations.svg"
      />
      <nav-button
        :to="{ name: 'settings' }"
        subtitle="Configure site-specific settings and edit zones"
        title="Site Settings"
        icon="/menu/settings.svg"
      />
      <nav-button
        :to="{ name: 'workflows' }"
        subtitle="view automation status, and update settings"
        title="Workflows & Automations"
        icon="/menu/workflows.svg"
      />
    </v-card>
  </v-menu>
</template>

<script setup>
import MenuIcon from '@/components/MenuIcon.vue';
import NavButton from '@/components/NavButton.vue';
import {usePage} from '@/components/page.js';
import {ref} from 'vue';

defineProps({
  btnClass: [String, Object]
});

const {themeColor} = usePage();

const showMenu = ref(false);

</script>

<style scoped>
.icon circle {
  stroke: none;
  fill: v-bind(logoColor);
}

.menu-card {
  display: grid;
  grid-auto-flow: column;
  grid-gap: 30px;
  justify-content: start;
  padding: 30px 42px;
  background: #ECEDED;
}
</style>
