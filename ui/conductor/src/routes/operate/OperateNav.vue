<template>
  <v-navigation-drawer app clipped floating color="transparent">
    <v-list color="#111721" dense nav>
      <v-list-item to="/operate/summary">
        <v-list-item-content> Site Operations </v-list-item-content>
      </v-list-item>
      <v-divider class="my-2"></v-divider>
      <v-subheader class="text-uppercase font-weight-bold">
        Status
      </v-subheader>
      <v-list-item-group>
        <v-list-group>
          <template v-slot:activator>
            <v-list-item-icon>
              <v-icon>mdi-domain</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>Buildings</v-list-item-title>
            </v-list-item-content>
          </template>

          <v-list-group no-action sub-group>
            <template v-slot:activator>
              <v-list-item-icon>
                <v-icon>mdi-floor-plan</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>Floor 1</v-list-item-title>
              </v-list-item-content>
            </template>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title class="text-caption"
                  >Meeting Room 6.01</v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title class="text-caption"
                  >Meeting Room 6.01</v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title class="text-caption"
                  >Meeting Room 6.01</v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
          </v-list-group>
          <v-list-group no-action sub-group>
            <template v-slot:activator>
              <v-list-item-icon>
                <v-icon>mdi-floor-plan</v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title>Floor 2</v-list-item-title>
              </v-list-item-content>
            </template>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title class="text-caption"
                  >Meeting Room 6.01</v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title class="text-caption"
                  >Meeting Room 6.01</v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title class="text-caption"
                  >Meeting Room 6.01</v-list-item-title
                >
              </v-list-item-content>
            </v-list-item>
          </v-list-group>
        </v-list-group>

        <v-list-item to="/operate/notifications">
          <v-list-item-icon>
            <v-icon>mdi-bell-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Notifications</v-list-item-title>
          </v-list-item-content>

          <v-chip color="#33BFD6" class="font-weight-bold text"> 4 </v-chip>
        </v-list-item>
      </v-list-item-group>

      <v-subheader class="text-uppercase font-weight-bold">
        Systems
      </v-subheader>
      <v-list-item-group>
        <v-list-item to="/operate/systems/lighting">
          <v-list-item-icon>
            <v-icon>mdi-lightbulb-variant-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Lighting</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>

      <v-subheader class="text-uppercase font-weight-bold">
        Analytics
      </v-subheader>
      <v-list-item-group>
        <v-list-item to="/operate/analytics/usage">
          <v-list-item-icon>
            <v-icon>mdi-lightning-bolt-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Energy Usage</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item to="/operate/analytics/occupancy">
          <v-list-item-icon>
            <v-icon>mdi-account-group-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Occupancy</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item to="/operate/analytics/environment">
          <v-list-item-icon>
            <v-icon>mdi-tune-vertical</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Environment</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>

      <v-subheader class="text-uppercase font-weight-bold">
        Reports
        <v-spacer />
        <v-icon>mdi-plus-box</v-icon>
      </v-subheader>
      <v-list-item-group>
        <v-list-item to="/operate/reports/usage">
          <v-list-item-icon>
            <v-icon>mdi-lightning-bolt-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Energy Usage</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item to="/operate/reports/occupancy">
          <v-list-item-icon>
            <v-icon>mdi-account-group-outline</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Occupancy</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item to="/operate/reports/environment">
          <v-list-item-icon>
            <v-icon>mdi-tune-vertical</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Environment</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>
    </v-list>
    <Footer />
  </v-navigation-drawer>
</template>

<script>
import Footer from "@/components/Footer.vue";
export default {
  name: "OperateNav",
};
</script>

<style scoped>
::v-deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
