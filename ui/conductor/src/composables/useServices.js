import {newActionTracker} from '@/api/resource';
import {startService, stopService} from '@/api/ui/services';
import {useErrorStore} from '@/components/ui-error/error';
import {useHubStore} from '@/stores/hub';
import {useSidebarStore} from '@/stores/sidebar';
import {useServicesStore} from '@/stores/services';
import {serviceName} from '@/util/proxy';
import {computed, onMounted, onUnmounted, reactive, ref, toValue, watch, watchEffect} from 'vue';

/**
 * @param {{
 * name: MaybeRefOrGetter<string>,
 * type: MaybeRefOrGetter<string>
 * }} props
 * @return {{
 * serviceList: ComputedRef<Service.AsObject[]>,
 * serviceCollection: Ref<Collection<Service.AsObject>>,
 * showService: function(Service.AsObject): void,
 * startStopTracker: ActionTracker,
 * search: Ref<string>,
 * node: ComputedRef<Node.AsObject>,
 * nodesListValues: ComputedRef<Node.AsObject[]>,
 * _startService: function(Service.AsObject): Promise<void>,
 * _stopService: function(Service.AsObject): Promise<void>
 * }}
 */
export default function(props) {
  const serviceStore = useServicesStore();
  const sidebar = useSidebarStore();
  const errors = useErrorStore();
  const hubStore = useHubStore();

  const startStopTracker = reactive(
      /** @type {ActionTracker<Service.AsObject>} */
      newActionTracker()
  );
  const serviceCollection = ref({});
  const search = ref('');

  // Available services to select from
  const serviceList = computed(() => {
    return Object.values(serviceCollection.value?.resources?.value ?? []).filter(service => {
      const type = toValue(props.type);
      return type === '' || type === 'all' || service.type === type;
    });
  });

  // Available nodes to select from
  const nodesListValues = computed(() => Object.values(hubStore.nodesList));


  // Watch for changes in node.value.name and update it if needed
  watchEffect(() => {
    if (!serviceStore.node?.name) {
      if (nodesListValues.value.length > 0) {
        serviceStore.node = nodesListValues.value[0];
      }
    }
  });


  watch(serviceCollection, () => {
    // todo: this causes us to load all pages, connect with paging logic instead
    if (serviceCollection.value) {
      serviceCollection.value.needsMorePages = true;
    }
  });
  // Watch for changes in the name prop and the name of the node and update the serviceCollection
  watch(
      [() => props.name, () => serviceStore.node],
      async ([newName, newNode]) => {
        if (serviceCollection.value.reset) serviceCollection.value.reset();

        const col = serviceStore.getService(
            newName, await newNode?.commsAddress, await newNode?.commsName
        ).servicesCollection;
        // reinitialize in case this service collection has been previously reset;
        col.init();
        col.query(newName);
        serviceCollection.value = col;
      },
      {immediate: true}
  );

  //
  //
  // SERVICE ACTIONS
  /**
   *
   * @param {Service.AsObject} service
   */
  function showService(service) {
    sidebar.visible = true;
    sidebar.title = service.id;
    sidebar.data = {service, config: JSON.parse(service.configRaw)};
  }

  /**
   *
   * @param {Service.AsObject} service
   */
  async function _startService(service) {
    // Update the data if the sidebar is open and the service is being started
    if (sidebar.visible && sidebar.data?.service?.id !== service.id) {
      sidebar.title = service.id;
      sidebar.data = {service, config: JSON.parse(service.configRaw)};
    }

    await startService({
      name: serviceName(await serviceStore.node.commsName, toValue(props.name)),
      id: service.id
    }, startStopTracker);
  }


  /**
   *
   * @param {Service.AsObject} service
   */
  async function _stopService(service) {
    // Update the data if the sidebar is open and the service is being stopped
    if (sidebar.visible && sidebar.data?.service?.id !== service.id) {
      sidebar.title = service.id;
      sidebar.data = {service, config: JSON.parse(service.configRaw)};
    }

    await stopService({
      name: serviceName(await serviceStore.node.commsName, toValue(props.name)),
      id: service.id
    }, startStopTracker);
  }

  // Watch for changes in the serviceList and update the data if needed
  // This is necessary if we want to update the status details in the sidebar
  // simultaneously with the status details in the service list
  // Mainly when the sidebar is open and the service is being started/stopped
  watch(
      serviceList,
      (newServiceList, oldServiceList) => {
        if (!sidebar.data?.service?.id) return;

        // Find the service in the new list that matches the id in data
        const updatedService = newServiceList.find(s => s.id === sidebar.data.service.id);

        if (updatedService) {
          // Check if the service has been updated by comparing it with the old list
          const oldService = oldServiceList.find(s => s.id === updatedService.id);

          // Perform a deep comparison if necessary, for now, we just check if the old service exists
          if (!oldService || JSON.stringify(updatedService) !== JSON.stringify(oldService)) {
            // Update the data with the new service data
            // Ensuring to parse the config if it's in a raw JSON string format
            sidebar.data = {
              service: updatedService,
              config: updatedService.configRaw ? JSON.parse(updatedService.configRaw) : {}
            };
          }
        }
      },
      {immediate: true, deep: true} // Watch for nested changes within the serviceList array
  );

  //
  //
  // UI error handling
  let unwatchErrors;
  let unwatchStartStopErrors;
  onMounted(() => {
    unwatchErrors = errors.registerCollection(serviceCollection);
    unwatchStartStopErrors = errors.registerTracker(startStopTracker);
  });
  onUnmounted(() => {
    if (unwatchErrors) unwatchErrors();
    if (unwatchStartStopErrors) unwatchStartStopErrors();
    serviceCollection.value.reset();
  });

  return {
    serviceList,
    serviceCollection,
    showService,
    startStopTracker,
    search,
    nodesListValues,
    _startService,
    _stopService
  };
}
