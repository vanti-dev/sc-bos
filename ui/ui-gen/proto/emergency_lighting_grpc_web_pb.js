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


var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var google_protobuf_duration_pb = require('google-protobuf/google/protobuf/duration_pb.js')
const proto = {};
proto.vanti = {};
proto.vanti.bsp = {};
proto.vanti.bsp.ew = require('./emergency_lighting_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.vanti.bsp.ew.EmergencyLightingApiClient =
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
proto.vanti.bsp.ew.EmergencyLightingApiPromiseClient =
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
 *   !proto.vanti.bsp.ew.GetEmergencyLightRequest,
 *   !proto.vanti.bsp.ew.EmergencyLight>}
 */
const methodDescriptor_EmergencyLightingApi_GetEmergencyLight = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.EmergencyLightingApi/GetEmergencyLight',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.GetEmergencyLightRequest,
    proto.vanti.bsp.ew.EmergencyLight,
    /**
     * @param {!proto.vanti.bsp.ew.GetEmergencyLightRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.EmergencyLight.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.GetEmergencyLightRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.EmergencyLight)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.EmergencyLight>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.EmergencyLightingApiClient.prototype.getEmergencyLight =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.EmergencyLightingApi/GetEmergencyLight',
          request,
          metadata || {},
          methodDescriptor_EmergencyLightingApi_GetEmergencyLight,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.GetEmergencyLightRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.EmergencyLight>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.EmergencyLightingApiPromiseClient.prototype.getEmergencyLight =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.EmergencyLightingApi/GetEmergencyLight',
          request,
          metadata || {},
          methodDescriptor_EmergencyLightingApi_GetEmergencyLight);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.ListEmergencyLightsRequest,
 *   !proto.vanti.bsp.ew.ListEmergencyLightsResponse>}
 */
const methodDescriptor_EmergencyLightingApi_ListEmergencyLights = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.EmergencyLightingApi/ListEmergencyLights',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.ListEmergencyLightsRequest,
    proto.vanti.bsp.ew.ListEmergencyLightsResponse,
    /**
     * @param {!proto.vanti.bsp.ew.ListEmergencyLightsRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.ListEmergencyLightsResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.ListEmergencyLightsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.ListEmergencyLightsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.ListEmergencyLightsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.EmergencyLightingApiClient.prototype.listEmergencyLights =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.EmergencyLightingApi/ListEmergencyLights',
          request,
          metadata || {},
          methodDescriptor_EmergencyLightingApi_ListEmergencyLights,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.ListEmergencyLightsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.ListEmergencyLightsResponse>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.EmergencyLightingApiPromiseClient.prototype.listEmergencyLights =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.EmergencyLightingApi/ListEmergencyLights',
          request,
          metadata || {},
          methodDescriptor_EmergencyLightingApi_ListEmergencyLights);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.ListEmergencyLightEventsRequest,
 *   !proto.vanti.bsp.ew.ListEmergencyLightEventsResponse>}
 */
const methodDescriptor_EmergencyLightingApi_ListEmergencyLightEvents = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.EmergencyLightingApi/ListEmergencyLightEvents',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.ListEmergencyLightEventsRequest,
    proto.vanti.bsp.ew.ListEmergencyLightEventsResponse,
    /**
     * @param {!proto.vanti.bsp.ew.ListEmergencyLightEventsRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.ListEmergencyLightEventsResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.ListEmergencyLightEventsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.ListEmergencyLightEventsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.ListEmergencyLightEventsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.EmergencyLightingApiClient.prototype.listEmergencyLightEvents =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.EmergencyLightingApi/ListEmergencyLightEvents',
          request,
          metadata || {},
          methodDescriptor_EmergencyLightingApi_ListEmergencyLightEvents,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.ListEmergencyLightEventsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.ListEmergencyLightEventsResponse>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.EmergencyLightingApiPromiseClient.prototype.listEmergencyLightEvents =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.EmergencyLightingApi/ListEmergencyLightEvents',
          request,
          metadata || {},
          methodDescriptor_EmergencyLightingApi_ListEmergencyLightEvents);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.GetReportCSVRequest,
 *   !proto.vanti.bsp.ew.ReportCSV>}
 */
const methodDescriptor_EmergencyLightingApi_GetReportCSV = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.EmergencyLightingApi/GetReportCSV',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.GetReportCSVRequest,
    proto.vanti.bsp.ew.ReportCSV,
    /**
     * @param {!proto.vanti.bsp.ew.GetReportCSVRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.ReportCSV.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.GetReportCSVRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.ReportCSV)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.ReportCSV>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.EmergencyLightingApiClient.prototype.getReportCSV =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.EmergencyLightingApi/GetReportCSV',
          request,
          metadata || {},
          methodDescriptor_EmergencyLightingApi_GetReportCSV,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.GetReportCSVRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.ReportCSV>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.EmergencyLightingApiPromiseClient.prototype.getReportCSV =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.EmergencyLightingApi/GetReportCSV',
          request,
          metadata || {},
          methodDescriptor_EmergencyLightingApi_GetReportCSV);
    };


module.exports = proto.vanti.bsp.ew;

