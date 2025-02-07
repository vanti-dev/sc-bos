/**
 * @fileoverview gRPC-Web generated client stub for smartcore.bos
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.5.0
// 	protoc              v5.29.1
// source: raise_lower.proto


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_field_mask_pb = require('google-protobuf/google/protobuf/field_mask_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var types_info_pb = require('@smart-core-os/sc-api-grpc-web/types/info_pb.js')
const proto = {};
proto.smartcore = {};
proto.smartcore.bos = require('./raise_lower_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.RaiseLowerApiClient =
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
proto.smartcore.bos.RaiseLowerApiPromiseClient =
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
 *   !proto.smartcore.bos.GetBearerStateRequest,
 *   !proto.smartcore.bos.BearerState>}
 */
const methodDescriptor_RaiseLowerApi_GetBearerState = new grpc.web.MethodDescriptor(
  '/smartcore.bos.RaiseLowerApi/GetBearerState',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.GetBearerStateRequest,
  proto.smartcore.bos.BearerState,
  /**
   * @param {!proto.smartcore.bos.GetBearerStateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.BearerState.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.GetBearerStateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.BearerState)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.BearerState>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.RaiseLowerApiClient.prototype.getBearerState =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.RaiseLowerApi/GetBearerState',
      request,
      metadata || {},
      methodDescriptor_RaiseLowerApi_GetBearerState,
      callback);
};


/**
 * @param {!proto.smartcore.bos.GetBearerStateRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.BearerState>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.RaiseLowerApiPromiseClient.prototype.getBearerState =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.RaiseLowerApi/GetBearerState',
      request,
      metadata || {},
      methodDescriptor_RaiseLowerApi_GetBearerState);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.PullBearerStateRequest,
 *   !proto.smartcore.bos.PullBearerStateResponse>}
 */
const methodDescriptor_RaiseLowerApi_PullBearerState = new grpc.web.MethodDescriptor(
  '/smartcore.bos.RaiseLowerApi/PullBearerState',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.smartcore.bos.PullBearerStateRequest,
  proto.smartcore.bos.PullBearerStateResponse,
  /**
   * @param {!proto.smartcore.bos.PullBearerStateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.PullBearerStateResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.PullBearerStateRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullBearerStateResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.RaiseLowerApiClient.prototype.pullBearerState =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.RaiseLowerApi/PullBearerState',
      request,
      metadata || {},
      methodDescriptor_RaiseLowerApi_PullBearerState);
};


/**
 * @param {!proto.smartcore.bos.PullBearerStateRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullBearerStateResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.RaiseLowerApiPromiseClient.prototype.pullBearerState =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.RaiseLowerApi/PullBearerState',
      request,
      metadata || {},
      methodDescriptor_RaiseLowerApi_PullBearerState);
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.RaiseLowerInfoClient =
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
proto.smartcore.bos.RaiseLowerInfoPromiseClient =
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
 *   !proto.smartcore.bos.DescribeBearerRequest,
 *   !proto.smartcore.bos.BearerStateSupport>}
 */
const methodDescriptor_RaiseLowerInfo_DescribeBearerState = new grpc.web.MethodDescriptor(
  '/smartcore.bos.RaiseLowerInfo/DescribeBearerState',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.DescribeBearerRequest,
  proto.smartcore.bos.BearerStateSupport,
  /**
   * @param {!proto.smartcore.bos.DescribeBearerRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.BearerStateSupport.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.DescribeBearerRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.BearerStateSupport)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.BearerStateSupport>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.RaiseLowerInfoClient.prototype.describeBearerState =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.RaiseLowerInfo/DescribeBearerState',
      request,
      metadata || {},
      methodDescriptor_RaiseLowerInfo_DescribeBearerState,
      callback);
};


/**
 * @param {!proto.smartcore.bos.DescribeBearerRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.BearerStateSupport>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.RaiseLowerInfoPromiseClient.prototype.describeBearerState =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.RaiseLowerInfo/DescribeBearerState',
      request,
      metadata || {},
      methodDescriptor_RaiseLowerInfo_DescribeBearerState);
};


module.exports = proto.smartcore.bos;

