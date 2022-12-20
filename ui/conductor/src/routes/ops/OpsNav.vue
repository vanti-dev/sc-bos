<template>
  <v-list class="pa-0" dense nav>
    <v-list-item-group>
      <v-list-group>
        <template #activator>
          <v-list-item-icon>
            <v-icon>mdi-domain</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Buildings</v-list-item-title>
          </v-list-item-content>
        </template>

        <v-list-group no-action sub-group>
          <template #activator>
            <v-list-item-icon>
              <v-icon>mdi-floor-plan</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>Floor 1</v-list-item-title>
            </v-list-item-content>
          </template>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption">
                Meeting Room 6.01
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption">
                Meeting Room 6.01
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption">
                Meeting Room 6.01
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-group>
        <v-list-group no-action sub-group>
          <template #activator>
            <v-list-item-icon>
              <v-icon>mdi-floor-plan</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>Floor 2</v-list-item-title>
            </v-list-item-content>
          </template>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption">
                Meeting Room 6.01
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption">
                Meeting Room 6.01
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title class="text-caption">
                Meeting Room 6.01
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-group>
      </v-list-group>
    </v-list-item-group>

    <v-list-item to="/ops/notifications">
      <v-list-item-icon>
        <v-icon>mdi-bell-outline</v-icon>
      </v-list-item-icon>
      <v-list-item-content>
        <v-list-item-title>Notifications</v-list-item-title>
      </v-list-item-content>

      <v-chip class="font-weight-bold text primary"> 4</v-chip>
    </v-list-item>
  </v-list>
</template>

<style scoped>
::v-deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
