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

var traits_metadata_pb = require('@smart-core-os/sc-api-grpc-web/traits/metadata_pb.js')

var types_change_pb = require('@smart-core-os/sc-api-grpc-web/types/change_pb.js')
const proto = {};
proto.smartcore = {};
proto.smartcore.bos = require('./devices_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.DevicesApiClient =
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
proto.smartcore.bos.DevicesApiPromiseClient =
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
 *   !proto.smartcore.bos.ListDevicesRequest,
 *   !proto.smartcore.bos.ListDevicesResponse>}
 */
const methodDescriptor_DevicesApi_ListDevices = new grpc.web.MethodDescriptor(
  '/smartcore.bos.DevicesApi/ListDevices',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.ListDevicesRequest,
  proto.smartcore.bos.ListDevicesResponse,
  /**
   * @param {!proto.smartcore.bos.ListDevicesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListDevicesResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.ListDevicesRequest,
 *   !proto.smartcore.bos.ListDevicesResponse>}
 */
const methodInfo_DevicesApi_ListDevices = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.ListDevicesResponse,
  /**
   * @param {!proto.smartcore.bos.ListDevicesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListDevicesResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.ListDevicesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.ListDevicesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.ListDevicesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.DevicesApiClient.prototype.listDevices =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.DevicesApi/ListDevices',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_ListDevices,
      callback);
};


/**
 * @param {!proto.smartcore.bos.ListDevicesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.ListDevicesResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.DevicesApiPromiseClient.prototype.listDevices =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.DevicesApi/ListDevices',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_ListDevices);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.PullDevicesRequest,
 *   !proto.smartcore.bos.PullDevicesResponse>}
 */
const methodDescriptor_DevicesApi_PullDevices = new grpc.web.MethodDescriptor(
  '/smartcore.bos.DevicesApi/PullDevices',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.smartcore.bos.PullDevicesRequest,
  proto.smartcore.bos.PullDevicesResponse,
  /**
   * @param {!proto.smartcore.bos.PullDevicesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.PullDevicesResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.PullDevicesRequest,
 *   !proto.smartcore.bos.PullDevicesResponse>}
 */
const methodInfo_DevicesApi_PullDevices = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.PullDevicesResponse,
  /**
   * @param {!proto.smartcore.bos.PullDevicesRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.PullDevicesResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.PullDevicesRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullDevicesResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.DevicesApiClient.prototype.pullDevices =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.DevicesApi/PullDevices',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_PullDevices);
};


/**
 * @param {!proto.smartcore.bos.PullDevicesRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullDevicesResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.DevicesApiPromiseClient.prototype.pullDevices =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.DevicesApi/PullDevices',
      request,
      metadata || {},
      methodDescriptor_DevicesApi_PullDevices);
};


module.exports = proto.smartcore.bos;

