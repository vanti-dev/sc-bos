/**
 * @fileoverview gRPC-Web generated client stub for smartcore.bos
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.4.2
// 	protoc              v3.21.12
// source: status.proto


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_field_mask_pb = require('google-protobuf/google/protobuf/field_mask_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var types_time_period_pb = require('@smart-core-os/sc-api-grpc-web/types/time/period_pb.js')
const proto = {};
proto.smartcore = {};
proto.smartcore.bos = require('./status_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.StatusApiClient =
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
proto.smartcore.bos.StatusApiPromiseClient =
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
 *   !proto.smartcore.bos.GetCurrentStatusRequest,
 *   !proto.smartcore.bos.StatusLog>}
 */
const methodDescriptor_StatusApi_GetCurrentStatus = new grpc.web.MethodDescriptor(
  '/smartcore.bos.StatusApi/GetCurrentStatus',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.GetCurrentStatusRequest,
  proto.smartcore.bos.StatusLog,
  /**
   * @param {!proto.smartcore.bos.GetCurrentStatusRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.StatusLog.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.GetCurrentStatusRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.StatusLog)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.StatusLog>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.StatusApiClient.prototype.getCurrentStatus =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.StatusApi/GetCurrentStatus',
      request,
      metadata || {},
      methodDescriptor_StatusApi_GetCurrentStatus,
      callback);
};


/**
 * @param {!proto.smartcore.bos.GetCurrentStatusRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.StatusLog>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.StatusApiPromiseClient.prototype.getCurrentStatus =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.StatusApi/GetCurrentStatus',
      request,
      metadata || {},
      methodDescriptor_StatusApi_GetCurrentStatus);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.PullCurrentStatusRequest,
 *   !proto.smartcore.bos.PullCurrentStatusResponse>}
 */
const methodDescriptor_StatusApi_PullCurrentStatus = new grpc.web.MethodDescriptor(
  '/smartcore.bos.StatusApi/PullCurrentStatus',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.smartcore.bos.PullCurrentStatusRequest,
  proto.smartcore.bos.PullCurrentStatusResponse,
  /**
   * @param {!proto.smartcore.bos.PullCurrentStatusRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.PullCurrentStatusResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.PullCurrentStatusRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullCurrentStatusResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.StatusApiClient.prototype.pullCurrentStatus =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.StatusApi/PullCurrentStatus',
      request,
      metadata || {},
      methodDescriptor_StatusApi_PullCurrentStatus);
};


/**
 * @param {!proto.smartcore.bos.PullCurrentStatusRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullCurrentStatusResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.StatusApiPromiseClient.prototype.pullCurrentStatus =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.StatusApi/PullCurrentStatus',
      request,
      metadata || {},
      methodDescriptor_StatusApi_PullCurrentStatus);
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.StatusHistoryClient =
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
proto.smartcore.bos.StatusHistoryPromiseClient =
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
 *   !proto.smartcore.bos.ListCurrentStatusHistoryRequest,
 *   !proto.smartcore.bos.ListCurrentStatusHistoryResponse>}
 */
const methodDescriptor_StatusHistory_ListCurrentStatusHistory = new grpc.web.MethodDescriptor(
  '/smartcore.bos.StatusHistory/ListCurrentStatusHistory',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.ListCurrentStatusHistoryRequest,
  proto.smartcore.bos.ListCurrentStatusHistoryResponse,
  /**
   * @param {!proto.smartcore.bos.ListCurrentStatusHistoryRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListCurrentStatusHistoryResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.ListCurrentStatusHistoryRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.ListCurrentStatusHistoryResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.ListCurrentStatusHistoryResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.StatusHistoryClient.prototype.listCurrentStatusHistory =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.StatusHistory/ListCurrentStatusHistory',
      request,
      metadata || {},
      methodDescriptor_StatusHistory_ListCurrentStatusHistory,
      callback);
};


/**
 * @param {!proto.smartcore.bos.ListCurrentStatusHistoryRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.ListCurrentStatusHistoryResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.StatusHistoryPromiseClient.prototype.listCurrentStatusHistory =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.StatusHistory/ListCurrentStatusHistory',
      request,
      metadata || {},
      methodDescriptor_StatusHistory_ListCurrentStatusHistory);
};


module.exports = proto.smartcore.bos;

