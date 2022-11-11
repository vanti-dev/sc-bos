<template>
  <div class="bgColor px-8">
    <div class="ma-4">
      <p class="font-weight-bold text-overline">Site Operations</p>
    </div>
    <v-divider></v-divider>
    <v-list color="#111721" dense>
      <v-subheader class="subheading text-uppercase font-weight-bold">
        Status
      </v-subheader>
      <v-list-item-group>
        <v-list-group>
          <template v-slot:activator>
            <v-list-item-icon>
              <v-icon>mdi-domain</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>Buildings Overview</v-list-item-title>
            </v-list-item-content>
          </template>

          <v-list-group no-action sub-group>
            <template v-slot:activator>
              <v-list-item-icon>
                <v-icon>mdi-selection-multiple-marker</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>Floor 1</v-list-item-title>
              </v-list-item-content>
            </template>
            <v-list-item>
              <v-list-item-icon>
                <v-icon>mdi-laptop-account</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>Meeting Room 6.01</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-group>
          <v-list-item>
            <v-list-item-icon>
              <v-icon>mdi-selection-multiple-marker</v-icon>
            </v-list-item-icon>

            <v-list-item-content>
              <v-list-item-title>Floor 2</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-group>

        <v-list-item>
          <v-list-item-icon>
            <v-icon>mdi-bell-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Notifications</v-list-item-title>
          </v-list-item-content>

          <v-chip color="#33BFD6" class="font-weight-bold text"> 4 </v-chip>
        </v-list-item>
      </v-list-item-group>

      <v-divider></v-divider>

      <v-subheader class="subheading text-uppercase font-weight-bold">
        Systems
      </v-subheader>
      <v-list-item-group>
        <v-list-item>
          <v-list-item-icon>
            <v-icon>mdi-lightbulb-variant-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Lighting</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>

      <v-subheader class="subheading text-uppercase font-weight-bold">
        Analytics
      </v-subheader>
      <v-list-item-group>
        <v-list-item>
          <v-list-item-icon>
            <v-icon>mdi-lightning-bolt-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Energy Usage</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-icon>
            <v-icon>mdi-account-group-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Occupancy</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-icon>
            <v-icon>mdi-tune-vertical</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Environment</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>

      <v-divider></v-divider>

      <v-subheader class="subheading text-uppercase font-weight-bold">
        Reports
        <v-spacer />
        <v-icon>mdi-plus-box</v-icon>
      </v-subheader>
      <v-list-item-group>
        <v-list-item>
          <v-list-item-icon>
            <v-icon>mdi-lightning-bolt-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Energy Usage</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-icon>
            <v-icon>mdi-account-group-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Occupancy</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-icon>
            <v-icon>mdi-tune-vertical</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Environment</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>
    </v-list>
  </div>
</template>

<script>
import { useOperateStore } from "@/stores/operate.js";
import { storeToRefs } from "pinia";

export default {
  setup() {
    const store = useOperateStore();
    const { meetings } = storeToRefs(store);

    return { store, meetings };
  },
};
</script>

<style lang="scss" scoped>
.bgColor {
  background: #111721;
}
</style>
