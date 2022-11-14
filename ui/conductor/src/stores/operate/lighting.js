import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useLightingStore = defineStore("lightingStore", () => {
  const lights = ref([
    {
      device_id: "LIT-L02_12-001",
      location: "L02_12",
      status: "On",
      battery_status: "100%",
      model: "Philips LED 1245812",
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
      model: "Philips LED 1245812",
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

  return {
    lights,
    headers,
    selected,
    meetings,
  };
});
