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


var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var google_protobuf_duration_pb = require('google-protobuf/google/protobuf/duration_pb.js')
const proto = {};
proto.smartcore = {};
proto.smartcore.bos = require('./emergency_lighting_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.EmergencyLightingApiClient =
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
proto.smartcore.bos.EmergencyLightingApiPromiseClient =
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
 *   !proto.smartcore.bos.GetEmergencyLightRequest,
 *   !proto.smartcore.bos.EmergencyLight>}
 */
const methodDescriptor_EmergencyLightingApi_GetEmergencyLight = new grpc.web.MethodDescriptor(
  '/smartcore.bos.EmergencyLightingApi/GetEmergencyLight',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.GetEmergencyLightRequest,
  proto.smartcore.bos.EmergencyLight,
  /**
   * @param {!proto.smartcore.bos.GetEmergencyLightRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.EmergencyLight.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.GetEmergencyLightRequest,
 *   !proto.smartcore.bos.EmergencyLight>}
 */
const methodInfo_EmergencyLightingApi_GetEmergencyLight = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.EmergencyLight,
  /**
   * @param {!proto.smartcore.bos.GetEmergencyLightRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.EmergencyLight.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.GetEmergencyLightRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.EmergencyLight)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.EmergencyLight>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.EmergencyLightingApiClient.prototype.getEmergencyLight =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.EmergencyLightingApi/GetEmergencyLight',
      request,
      metadata || {},
      methodDescriptor_EmergencyLightingApi_GetEmergencyLight,
      callback);
};


/**
 * @param {!proto.smartcore.bos.GetEmergencyLightRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.EmergencyLight>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.EmergencyLightingApiPromiseClient.prototype.getEmergencyLight =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.EmergencyLightingApi/GetEmergencyLight',
      request,
      metadata || {},
      methodDescriptor_EmergencyLightingApi_GetEmergencyLight);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.ListEmergencyLightsRequest,
 *   !proto.smartcore.bos.ListEmergencyLightsResponse>}
 */
const methodDescriptor_EmergencyLightingApi_ListEmergencyLights = new grpc.web.MethodDescriptor(
  '/smartcore.bos.EmergencyLightingApi/ListEmergencyLights',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.ListEmergencyLightsRequest,
  proto.smartcore.bos.ListEmergencyLightsResponse,
  /**
   * @param {!proto.smartcore.bos.ListEmergencyLightsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListEmergencyLightsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.ListEmergencyLightsRequest,
 *   !proto.smartcore.bos.ListEmergencyLightsResponse>}
 */
const methodInfo_EmergencyLightingApi_ListEmergencyLights = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.ListEmergencyLightsResponse,
  /**
   * @param {!proto.smartcore.bos.ListEmergencyLightsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListEmergencyLightsResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.ListEmergencyLightsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.ListEmergencyLightsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.ListEmergencyLightsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.EmergencyLightingApiClient.prototype.listEmergencyLights =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.EmergencyLightingApi/ListEmergencyLights',
      request,
      metadata || {},
      methodDescriptor_EmergencyLightingApi_ListEmergencyLights,
      callback);
};


/**
 * @param {!proto.smartcore.bos.ListEmergencyLightsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.ListEmergencyLightsResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.EmergencyLightingApiPromiseClient.prototype.listEmergencyLights =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.EmergencyLightingApi/ListEmergencyLights',
      request,
      metadata || {},
      methodDescriptor_EmergencyLightingApi_ListEmergencyLights);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.ListEmergencyLightEventsRequest,
 *   !proto.smartcore.bos.ListEmergencyLightEventsResponse>}
 */
const methodDescriptor_EmergencyLightingApi_ListEmergencyLightEvents = new grpc.web.MethodDescriptor(
  '/smartcore.bos.EmergencyLightingApi/ListEmergencyLightEvents',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.ListEmergencyLightEventsRequest,
  proto.smartcore.bos.ListEmergencyLightEventsResponse,
  /**
   * @param {!proto.smartcore.bos.ListEmergencyLightEventsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListEmergencyLightEventsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.ListEmergencyLightEventsRequest,
 *   !proto.smartcore.bos.ListEmergencyLightEventsResponse>}
 */
const methodInfo_EmergencyLightingApi_ListEmergencyLightEvents = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.ListEmergencyLightEventsResponse,
  /**
   * @param {!proto.smartcore.bos.ListEmergencyLightEventsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListEmergencyLightEventsResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.ListEmergencyLightEventsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.ListEmergencyLightEventsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.ListEmergencyLightEventsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.EmergencyLightingApiClient.prototype.listEmergencyLightEvents =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.EmergencyLightingApi/ListEmergencyLightEvents',
      request,
      metadata || {},
      methodDescriptor_EmergencyLightingApi_ListEmergencyLightEvents,
      callback);
};


/**
 * @param {!proto.smartcore.bos.ListEmergencyLightEventsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.ListEmergencyLightEventsResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.EmergencyLightingApiPromiseClient.prototype.listEmergencyLightEvents =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.EmergencyLightingApi/ListEmergencyLightEvents',
      request,
      metadata || {},
      methodDescriptor_EmergencyLightingApi_ListEmergencyLightEvents);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.GetReportCSVRequest,
 *   !proto.smartcore.bos.ReportCSV>}
 */
const methodDescriptor_EmergencyLightingApi_GetReportCSV = new grpc.web.MethodDescriptor(
  '/smartcore.bos.EmergencyLightingApi/GetReportCSV',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.GetReportCSVRequest,
  proto.smartcore.bos.ReportCSV,
  /**
   * @param {!proto.smartcore.bos.GetReportCSVRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ReportCSV.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.GetReportCSVRequest,
 *   !proto.smartcore.bos.ReportCSV>}
 */
const methodInfo_EmergencyLightingApi_GetReportCSV = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.ReportCSV,
  /**
   * @param {!proto.smartcore.bos.GetReportCSVRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ReportCSV.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.GetReportCSVRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.ReportCSV)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.ReportCSV>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.EmergencyLightingApiClient.prototype.getReportCSV =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.EmergencyLightingApi/GetReportCSV',
      request,
      metadata || {},
      methodDescriptor_EmergencyLightingApi_GetReportCSV,
      callback);
};


/**
 * @param {!proto.smartcore.bos.GetReportCSVRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.ReportCSV>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.EmergencyLightingApiPromiseClient.prototype.getReportCSV =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.EmergencyLightingApi/GetReportCSV',
      request,
      metadata || {},
      methodDescriptor_EmergencyLightingApi_GetReportCSV);
};


module.exports = proto.smartcore.bos;

