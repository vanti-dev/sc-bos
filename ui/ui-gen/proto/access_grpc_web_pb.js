/**
 * @fileoverview gRPC-Web generated client stub for smartcore.bos
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.5.0
// 	protoc              v4.25.0-rc1
// source: access.proto


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_field_mask_pb = require('google-protobuf/google/protobuf/field_mask_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')

var types_image_pb = require('@smart-core-os/sc-api-grpc-web/types/image_pb.js')
const proto = {};
proto.smartcore = {};
proto.smartcore.bos = require('./access_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.AccessApiClient =
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
proto.smartcore.bos.AccessApiPromiseClient =
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
 *   !proto.smartcore.bos.GetLastAccessAttemptRequest,
 *   !proto.smartcore.bos.AccessAttempt>}
 */
const methodDescriptor_AccessApi_GetLastAccessAttempt = new grpc.web.MethodDescriptor(
  '/smartcore.bos.AccessApi/GetLastAccessAttempt',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.GetLastAccessAttemptRequest,
  proto.smartcore.bos.AccessAttempt,
  /**
   * @param {!proto.smartcore.bos.GetLastAccessAttemptRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.AccessAttempt.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.GetLastAccessAttemptRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.AccessAttempt)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.AccessAttempt>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AccessApiClient.prototype.getLastAccessAttempt =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.AccessApi/GetLastAccessAttempt',
      request,
      metadata || {},
      methodDescriptor_AccessApi_GetLastAccessAttempt,
      callback);
};


/**
 * @param {!proto.smartcore.bos.GetLastAccessAttemptRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.AccessAttempt>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.AccessApiPromiseClient.prototype.getLastAccessAttempt =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.AccessApi/GetLastAccessAttempt',
      request,
      metadata || {},
      methodDescriptor_AccessApi_GetLastAccessAttempt);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.PullAccessAttemptsRequest,
 *   !proto.smartcore.bos.PullAccessAttemptsResponse>}
 */
const methodDescriptor_AccessApi_PullAccessAttempts = new grpc.web.MethodDescriptor(
  '/smartcore.bos.AccessApi/PullAccessAttempts',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.smartcore.bos.PullAccessAttemptsRequest,
  proto.smartcore.bos.PullAccessAttemptsResponse,
  /**
   * @param {!proto.smartcore.bos.PullAccessAttemptsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.PullAccessAttemptsResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.PullAccessAttemptsRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullAccessAttemptsResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AccessApiClient.prototype.pullAccessAttempts =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.AccessApi/PullAccessAttempts',
      request,
      metadata || {},
      methodDescriptor_AccessApi_PullAccessAttempts);
};


/**
 * @param {!proto.smartcore.bos.PullAccessAttemptsRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullAccessAttemptsResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AccessApiPromiseClient.prototype.pullAccessAttempts =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.AccessApi/PullAccessAttempts',
      request,
      metadata || {},
      methodDescriptor_AccessApi_PullAccessAttempts);
};


module.exports = proto.smartcore.bos;

