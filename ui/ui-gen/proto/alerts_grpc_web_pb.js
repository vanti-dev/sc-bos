/**
 * @fileoverview gRPC-Web generated client stub for vanti.bsp.ew
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_field_mask_pb = require('google-protobuf/google/protobuf/field_mask_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var types_change_pb = require('@smart-core-os/sc-api-grpc-web/types/change_pb.js')
const proto = {};
proto.vanti = {};
proto.vanti.bsp = {};
proto.vanti.bsp.ew = require('./alerts_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.vanti.bsp.ew.AlertApiClient =
    function(hostname, credentials, options) {
      if (!options) options = {};
      options.format = 'text';

      /**
       * @private @const {!grpc.web.GrpcWebClientBase} The client
       */
      this.client_ = new grpc.web.GrpcWebClientBase(options);

      /**
       * @private @const {string} The hostname
       */
      this.hostname_ = hostname;

    };


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.vanti.bsp.ew.AlertApiPromiseClient =
    function(hostname, credentials, options) {
      if (!options) options = {};
      options.format = 'text';

      /**
       * @private @const {!grpc.web.GrpcWebClientBase} The client
       */
      this.client_ = new grpc.web.GrpcWebClientBase(options);

      /**
       * @private @const {string} The hostname
       */
      this.hostname_ = hostname;

    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.ListAlertsRequest,
 *   !proto.vanti.bsp.ew.ListAlertsResponse>}
 */
const methodDescriptor_AlertApi_ListAlerts = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.AlertApi/ListAlerts',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.ListAlertsRequest,
    proto.vanti.bsp.ew.ListAlertsResponse,
    /**
     * @param {!proto.vanti.bsp.ew.ListAlertsRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.ListAlertsResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.ListAlertsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.ListAlertsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.ListAlertsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.AlertApiClient.prototype.listAlerts =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.AlertApi/ListAlerts',
          request,
          metadata || {},
          methodDescriptor_AlertApi_ListAlerts,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.ListAlertsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.ListAlertsResponse>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.AlertApiPromiseClient.prototype.listAlerts =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.AlertApi/ListAlerts',
          request,
          metadata || {},
          methodDescriptor_AlertApi_ListAlerts);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.PullAlertsRequest,
 *   !proto.vanti.bsp.ew.PullAlertsResponse>}
 */
const methodDescriptor_AlertApi_PullAlerts = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.AlertApi/PullAlerts',
    grpc.web.MethodType.SERVER_STREAMING,
    proto.vanti.bsp.ew.PullAlertsRequest,
    proto.vanti.bsp.ew.PullAlertsResponse,
    /**
     * @param {!proto.vanti.bsp.ew.PullAlertsRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.PullAlertsResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.PullAlertsRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.PullAlertsResponse>}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.AlertApiClient.prototype.pullAlerts =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/vanti.bsp.ew.AlertApi/PullAlerts',
          request,
          metadata || {},
          methodDescriptor_AlertApi_PullAlerts);
    };


/**
 * @param {!proto.vanti.bsp.ew.PullAlertsRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.PullAlertsResponse>}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.AlertApiPromiseClient.prototype.pullAlerts =
    function(request, metadata) {
      return this.client_.serverStreaming(this.hostname_ +
          '/vanti.bsp.ew.AlertApi/PullAlerts',
          request,
          metadata || {},
          methodDescriptor_AlertApi_PullAlerts);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.AcknowledgeAlertRequest,
 *   !proto.vanti.bsp.ew.Alert>}
 */
const methodDescriptor_AlertApi_AcknowledgeAlert = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.AlertApi/AcknowledgeAlert',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.AcknowledgeAlertRequest,
    proto.vanti.bsp.ew.Alert,
    /**
     * @param {!proto.vanti.bsp.ew.AcknowledgeAlertRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.Alert.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.AcknowledgeAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.Alert)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.Alert>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.AlertApiClient.prototype.acknowledgeAlert =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.AlertApi/AcknowledgeAlert',
          request,
          metadata || {},
          methodDescriptor_AlertApi_AcknowledgeAlert,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.AcknowledgeAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.Alert>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.AlertApiPromiseClient.prototype.acknowledgeAlert =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.AlertApi/AcknowledgeAlert',
          request,
          metadata || {},
          methodDescriptor_AlertApi_AcknowledgeAlert);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.AcknowledgeAlertRequest,
 *   !proto.vanti.bsp.ew.Alert>}
 */
const methodDescriptor_AlertApi_UnacknowledgeAlert = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.AlertApi/UnacknowledgeAlert',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.AcknowledgeAlertRequest,
    proto.vanti.bsp.ew.Alert,
    /**
     * @param {!proto.vanti.bsp.ew.AcknowledgeAlertRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.Alert.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.AcknowledgeAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.Alert)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.Alert>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.AlertApiClient.prototype.unacknowledgeAlert =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.AlertApi/UnacknowledgeAlert',
          request,
          metadata || {},
          methodDescriptor_AlertApi_UnacknowledgeAlert,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.AcknowledgeAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.Alert>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.AlertApiPromiseClient.prototype.unacknowledgeAlert =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.AlertApi/UnacknowledgeAlert',
          request,
          metadata || {},
          methodDescriptor_AlertApi_UnacknowledgeAlert);
    };


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.vanti.bsp.ew.AlertAdminApiClient =
    function(hostname, credentials, options) {
      if (!options) options = {};
      options.format = 'text';

      /**
       * @private @const {!grpc.web.GrpcWebClientBase} The client
       */
      this.client_ = new grpc.web.GrpcWebClientBase(options);

      /**
       * @private @const {string} The hostname
       */
      this.hostname_ = hostname;

    };


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.vanti.bsp.ew.AlertAdminApiPromiseClient =
    function(hostname, credentials, options) {
      if (!options) options = {};
      options.format = 'text';

      /**
       * @private @const {!grpc.web.GrpcWebClientBase} The client
       */
      this.client_ = new grpc.web.GrpcWebClientBase(options);

      /**
       * @private @const {string} The hostname
       */
      this.hostname_ = hostname;

    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.CreateAlertRequest,
 *   !proto.vanti.bsp.ew.Alert>}
 */
const methodDescriptor_AlertAdminApi_CreateAlert = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.AlertAdminApi/CreateAlert',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.CreateAlertRequest,
    proto.vanti.bsp.ew.Alert,
    /**
     * @param {!proto.vanti.bsp.ew.CreateAlertRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.Alert.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.CreateAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.Alert)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.Alert>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.AlertAdminApiClient.prototype.createAlert =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.AlertAdminApi/CreateAlert',
          request,
          metadata || {},
          methodDescriptor_AlertAdminApi_CreateAlert,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.CreateAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.Alert>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.AlertAdminApiPromiseClient.prototype.createAlert =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.AlertAdminApi/CreateAlert',
          request,
          metadata || {},
          methodDescriptor_AlertAdminApi_CreateAlert);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.UpdateAlertRequest,
 *   !proto.vanti.bsp.ew.Alert>}
 */
const methodDescriptor_AlertAdminApi_UpdateAlert = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.AlertAdminApi/UpdateAlert',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.UpdateAlertRequest,
    proto.vanti.bsp.ew.Alert,
    /**
     * @param {!proto.vanti.bsp.ew.UpdateAlertRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.Alert.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.UpdateAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.Alert)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.Alert>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.AlertAdminApiClient.prototype.updateAlert =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.AlertAdminApi/UpdateAlert',
          request,
          metadata || {},
          methodDescriptor_AlertAdminApi_UpdateAlert,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.UpdateAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.Alert>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.AlertAdminApiPromiseClient.prototype.updateAlert =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.AlertAdminApi/UpdateAlert',
          request,
          metadata || {},
          methodDescriptor_AlertAdminApi_UpdateAlert);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.DeleteAlertRequest,
 *   !proto.vanti.bsp.ew.DeleteAlertResponse>}
 */
const methodDescriptor_AlertAdminApi_DeleteAlert = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.AlertAdminApi/DeleteAlert',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.DeleteAlertRequest,
    proto.vanti.bsp.ew.DeleteAlertResponse,
    /**
     * @param {!proto.vanti.bsp.ew.DeleteAlertRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.DeleteAlertResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.DeleteAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.DeleteAlertResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.DeleteAlertResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.AlertAdminApiClient.prototype.deleteAlert =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.AlertAdminApi/DeleteAlert',
          request,
          metadata || {},
          methodDescriptor_AlertAdminApi_DeleteAlert,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.DeleteAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.DeleteAlertResponse>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.AlertAdminApiPromiseClient.prototype.deleteAlert =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.AlertAdminApi/DeleteAlert',
          request,
          metadata || {},
          methodDescriptor_AlertAdminApi_DeleteAlert);
    };


module.exports = proto.vanti.bsp.ew;

