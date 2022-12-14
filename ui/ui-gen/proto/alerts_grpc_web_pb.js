/**
 * @fileoverview gRPC-Web generated client stub for smartcore.bos
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
proto.smartcore = {};
proto.smartcore.bos = require('./alerts_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.AlertApiClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

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
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.AlertApiPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

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
 *   !proto.smartcore.bos.ListAlertsRequest,
 *   !proto.smartcore.bos.ListAlertsResponse>}
 */
const methodDescriptor_AlertApi_ListAlerts = new grpc.web.MethodDescriptor(
  '/smartcore.bos.AlertApi/ListAlerts',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.ListAlertsRequest,
  proto.smartcore.bos.ListAlertsResponse,
  /**
   * @param {!proto.smartcore.bos.ListAlertsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListAlertsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.ListAlertsRequest,
 *   !proto.smartcore.bos.ListAlertsResponse>}
 */
const methodInfo_AlertApi_ListAlerts = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.ListAlertsResponse,
  /**
   * @param {!proto.smartcore.bos.ListAlertsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListAlertsResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.ListAlertsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.ListAlertsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.ListAlertsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AlertApiClient.prototype.listAlerts =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.AlertApi/ListAlerts',
      request,
      metadata || {},
      methodDescriptor_AlertApi_ListAlerts,
      callback);
};


/**
 * @param {!proto.smartcore.bos.ListAlertsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.ListAlertsResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.AlertApiPromiseClient.prototype.listAlerts =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.AlertApi/ListAlerts',
      request,
      metadata || {},
      methodDescriptor_AlertApi_ListAlerts);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.PullAlertsRequest,
 *   !proto.smartcore.bos.PullAlertsResponse>}
 */
const methodDescriptor_AlertApi_PullAlerts = new grpc.web.MethodDescriptor(
  '/smartcore.bos.AlertApi/PullAlerts',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.smartcore.bos.PullAlertsRequest,
  proto.smartcore.bos.PullAlertsResponse,
  /**
   * @param {!proto.smartcore.bos.PullAlertsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.PullAlertsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.PullAlertsRequest,
 *   !proto.smartcore.bos.PullAlertsResponse>}
 */
const methodInfo_AlertApi_PullAlerts = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.PullAlertsResponse,
  /**
   * @param {!proto.smartcore.bos.PullAlertsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.PullAlertsResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.PullAlertsRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullAlertsResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AlertApiClient.prototype.pullAlerts =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.AlertApi/PullAlerts',
      request,
      metadata || {},
      methodDescriptor_AlertApi_PullAlerts);
};


/**
 * @param {!proto.smartcore.bos.PullAlertsRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullAlertsResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AlertApiPromiseClient.prototype.pullAlerts =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.AlertApi/PullAlerts',
      request,
      metadata || {},
      methodDescriptor_AlertApi_PullAlerts);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.AcknowledgeAlertRequest,
 *   !proto.smartcore.bos.Alert>}
 */
const methodDescriptor_AlertApi_AcknowledgeAlert = new grpc.web.MethodDescriptor(
  '/smartcore.bos.AlertApi/AcknowledgeAlert',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.AcknowledgeAlertRequest,
  proto.smartcore.bos.Alert,
  /**
   * @param {!proto.smartcore.bos.AcknowledgeAlertRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.Alert.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.AcknowledgeAlertRequest,
 *   !proto.smartcore.bos.Alert>}
 */
const methodInfo_AlertApi_AcknowledgeAlert = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.Alert,
  /**
   * @param {!proto.smartcore.bos.AcknowledgeAlertRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.Alert.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.AcknowledgeAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.Alert)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.Alert>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AlertApiClient.prototype.acknowledgeAlert =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.AlertApi/AcknowledgeAlert',
      request,
      metadata || {},
      methodDescriptor_AlertApi_AcknowledgeAlert,
      callback);
};


/**
 * @param {!proto.smartcore.bos.AcknowledgeAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.Alert>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.AlertApiPromiseClient.prototype.acknowledgeAlert =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.AlertApi/AcknowledgeAlert',
      request,
      metadata || {},
      methodDescriptor_AlertApi_AcknowledgeAlert);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.AcknowledgeAlertRequest,
 *   !proto.smartcore.bos.Alert>}
 */
const methodDescriptor_AlertApi_UnacknowledgeAlert = new grpc.web.MethodDescriptor(
  '/smartcore.bos.AlertApi/UnacknowledgeAlert',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.AcknowledgeAlertRequest,
  proto.smartcore.bos.Alert,
  /**
   * @param {!proto.smartcore.bos.AcknowledgeAlertRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.Alert.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.AcknowledgeAlertRequest,
 *   !proto.smartcore.bos.Alert>}
 */
const methodInfo_AlertApi_UnacknowledgeAlert = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.Alert,
  /**
   * @param {!proto.smartcore.bos.AcknowledgeAlertRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.Alert.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.AcknowledgeAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.Alert)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.Alert>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AlertApiClient.prototype.unacknowledgeAlert =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.AlertApi/UnacknowledgeAlert',
      request,
      metadata || {},
      methodDescriptor_AlertApi_UnacknowledgeAlert,
      callback);
};


/**
 * @param {!proto.smartcore.bos.AcknowledgeAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.Alert>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.AlertApiPromiseClient.prototype.unacknowledgeAlert =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.AlertApi/UnacknowledgeAlert',
      request,
      metadata || {},
      methodDescriptor_AlertApi_UnacknowledgeAlert);
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.AlertAdminApiClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

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
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.AlertAdminApiPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

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
 *   !proto.smartcore.bos.CreateAlertRequest,
 *   !proto.smartcore.bos.Alert>}
 */
const methodDescriptor_AlertAdminApi_CreateAlert = new grpc.web.MethodDescriptor(
  '/smartcore.bos.AlertAdminApi/CreateAlert',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.CreateAlertRequest,
  proto.smartcore.bos.Alert,
  /**
   * @param {!proto.smartcore.bos.CreateAlertRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.Alert.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.CreateAlertRequest,
 *   !proto.smartcore.bos.Alert>}
 */
const methodInfo_AlertAdminApi_CreateAlert = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.Alert,
  /**
   * @param {!proto.smartcore.bos.CreateAlertRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.Alert.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.CreateAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.Alert)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.Alert>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AlertAdminApiClient.prototype.createAlert =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.AlertAdminApi/CreateAlert',
      request,
      metadata || {},
      methodDescriptor_AlertAdminApi_CreateAlert,
      callback);
};


/**
 * @param {!proto.smartcore.bos.CreateAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.Alert>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.AlertAdminApiPromiseClient.prototype.createAlert =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.AlertAdminApi/CreateAlert',
      request,
      metadata || {},
      methodDescriptor_AlertAdminApi_CreateAlert);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.UpdateAlertRequest,
 *   !proto.smartcore.bos.Alert>}
 */
const methodDescriptor_AlertAdminApi_UpdateAlert = new grpc.web.MethodDescriptor(
  '/smartcore.bos.AlertAdminApi/UpdateAlert',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.UpdateAlertRequest,
  proto.smartcore.bos.Alert,
  /**
   * @param {!proto.smartcore.bos.UpdateAlertRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.Alert.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.UpdateAlertRequest,
 *   !proto.smartcore.bos.Alert>}
 */
const methodInfo_AlertAdminApi_UpdateAlert = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.Alert,
  /**
   * @param {!proto.smartcore.bos.UpdateAlertRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.Alert.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.UpdateAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.Alert)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.Alert>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AlertAdminApiClient.prototype.updateAlert =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.AlertAdminApi/UpdateAlert',
      request,
      metadata || {},
      methodDescriptor_AlertAdminApi_UpdateAlert,
      callback);
};


/**
 * @param {!proto.smartcore.bos.UpdateAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.Alert>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.AlertAdminApiPromiseClient.prototype.updateAlert =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.AlertAdminApi/UpdateAlert',
      request,
      metadata || {},
      methodDescriptor_AlertAdminApi_UpdateAlert);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.DeleteAlertRequest,
 *   !proto.smartcore.bos.DeleteAlertResponse>}
 */
const methodDescriptor_AlertAdminApi_DeleteAlert = new grpc.web.MethodDescriptor(
  '/smartcore.bos.AlertAdminApi/DeleteAlert',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.DeleteAlertRequest,
  proto.smartcore.bos.DeleteAlertResponse,
  /**
   * @param {!proto.smartcore.bos.DeleteAlertRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.DeleteAlertResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.DeleteAlertRequest,
 *   !proto.smartcore.bos.DeleteAlertResponse>}
 */
const methodInfo_AlertAdminApi_DeleteAlert = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.DeleteAlertResponse,
  /**
   * @param {!proto.smartcore.bos.DeleteAlertRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.DeleteAlertResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.DeleteAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.DeleteAlertResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.DeleteAlertResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AlertAdminApiClient.prototype.deleteAlert =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.AlertAdminApi/DeleteAlert',
      request,
      metadata || {},
      methodDescriptor_AlertAdminApi_DeleteAlert,
      callback);
};


/**
 * @param {!proto.smartcore.bos.DeleteAlertRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.DeleteAlertResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.AlertAdminApiPromiseClient.prototype.deleteAlert =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.AlertAdminApi/DeleteAlert',
      request,
      metadata || {},
      methodDescriptor_AlertAdminApi_DeleteAlert);
};


module.exports = proto.smartcore.bos;

