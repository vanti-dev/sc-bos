/**
 * @fileoverview gRPC-Web generated client stub for smartcore.bos
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.5.0
// 	protoc              v5.29.1
// source: temperature.proto


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_field_mask_pb = require('google-protobuf/google/protobuf/field_mask_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var types_unit_pb = require('@smart-core-os/sc-api-grpc-web/types/unit_pb.js')

var priority_pb = require('./priority_pb.js')
const proto = {};
proto.smartcore = {};
proto.smartcore.bos = require('./temperature_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.TemperatureApiClient =
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
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.TemperatureApiPromiseClient =
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
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.GetTemperatureRequest,
 *   !proto.smartcore.bos.Temperature>}
 */
const methodDescriptor_TemperatureApi_GetTemperature = new grpc.web.MethodDescriptor(
  '/smartcore.bos.TemperatureApi/GetTemperature',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.GetTemperatureRequest,
  proto.smartcore.bos.Temperature,
  /**
   * @param {!proto.smartcore.bos.GetTemperatureRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.Temperature.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.GetTemperatureRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.Temperature)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.Temperature>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.TemperatureApiClient.prototype.getTemperature =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.TemperatureApi/GetTemperature',
      request,
      metadata || {},
      methodDescriptor_TemperatureApi_GetTemperature,
      callback);
};


/**
 * @param {!proto.smartcore.bos.GetTemperatureRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.Temperature>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.TemperatureApiPromiseClient.prototype.getTemperature =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.TemperatureApi/GetTemperature',
      request,
      metadata || {},
      methodDescriptor_TemperatureApi_GetTemperature);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.PullTemperatureRequest,
 *   !proto.smartcore.bos.PullTemperatureResponse>}
 */
const methodDescriptor_TemperatureApi_PullTemperature = new grpc.web.MethodDescriptor(
  '/smartcore.bos.TemperatureApi/PullTemperature',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.smartcore.bos.PullTemperatureRequest,
  proto.smartcore.bos.PullTemperatureResponse,
  /**
   * @param {!proto.smartcore.bos.PullTemperatureRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.PullTemperatureResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.PullTemperatureRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullTemperatureResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.TemperatureApiClient.prototype.pullTemperature =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.TemperatureApi/PullTemperature',
      request,
      metadata || {},
      methodDescriptor_TemperatureApi_PullTemperature);
};


/**
 * @param {!proto.smartcore.bos.PullTemperatureRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullTemperatureResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.TemperatureApiPromiseClient.prototype.pullTemperature =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.TemperatureApi/PullTemperature',
      request,
      metadata || {},
      methodDescriptor_TemperatureApi_PullTemperature);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.UpdateTemperatureRequest,
 *   !proto.smartcore.bos.Temperature>}
 */
const methodDescriptor_TemperatureApi_UpdateTemperature = new grpc.web.MethodDescriptor(
  '/smartcore.bos.TemperatureApi/UpdateTemperature',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.UpdateTemperatureRequest,
  proto.smartcore.bos.Temperature,
  /**
   * @param {!proto.smartcore.bos.UpdateTemperatureRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.Temperature.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.UpdateTemperatureRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.Temperature)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.Temperature>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.TemperatureApiClient.prototype.updateTemperature =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.TemperatureApi/UpdateTemperature',
      request,
      metadata || {},
      methodDescriptor_TemperatureApi_UpdateTemperature,
      callback);
};


/**
 * @param {!proto.smartcore.bos.UpdateTemperatureRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.Temperature>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.TemperatureApiPromiseClient.prototype.updateTemperature =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.TemperatureApi/UpdateTemperature',
      request,
      metadata || {},
      methodDescriptor_TemperatureApi_UpdateTemperature);
};


module.exports = proto.smartcore.bos;

