import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useLightingStore = defineStore("lightingStore", () => {
  const lights = ref([
    {
      device_id: "LIT-L02_12-001",
      location: "L02_12",
      status: "On",
      battery_status: "100%",
      model: "Philips LED 1245813",
    },
    {
      device_id: "LIT-L02_12-002",
      location: "L02_12",
      status: "Off",
      battery_status: "100%",
      model: "Philips LED 1245812",
    },
    {
      device_id: "LIT-L02_12-003",
      location: "L02_12",
      status: "On",
      battery_status: "-",
      model: "Philips LED 1245814",
    },
    {
      device_id: "LIT-L02_12-004",
      location: "L02_12",
      status: "On",
      battery_status: "80%",
      model: "Philips LED 1245812",
    },
    {
      device_id: "LIT-L02_12-005",
      location: "L02_12",
      status: "On",
      battery_status: "80%",
      model: "Philips LED 1245812",
    },
  ]);

  const headers = ref([
    {
      text: "Device ID",
      align: "start",
      sortable: false,
      value: "device_id",
    },
    { text: "Location", value: "location" },
    { text: "Status", value: "status" },
    { text: "Battery Status", value: "battery_status" },
    { text: "Model", value: "model" },
  ]);

  const selected = ref([]);

  const meetings = ref([
    ["Meeting Room 6.01", "mdi-account-multiple-outline"],
    ["Meeting Room 6.01", "mdi-account-multiple-outline"],
  ]);

  const status = ref("All");

  const model = ref("All");

  const search = ref("");

  //filter for model or status

  const filteredLights = computed(() =>
    lights.value.filter((light) => {
      if (status.value === "All" && model.value === "All") {
        return true;
      } else if (status.value === "All") {
        return light.model === model.value;
      } else if (model.value === "All") {
        return light.status === status.value;
      } else {
        return light.status === status.value && light.model === model.value;
      }
    })
  );

  return {
    lights,
    headers,
    selected,
    meetings,
    status,
    model,
    search,
    filteredLights,
  };
});
