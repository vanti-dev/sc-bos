/**
 * @fileoverview gRPC-Web generated client stub for smartcore.bos
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.4.2
// 	protoc              v3.21.12
// source: priority.proto


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.smartcore = {};
proto.smartcore.bos = require('./priority_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.PriorityApiClient =
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
proto.smartcore.bos.PriorityApiPromiseClient =
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
 *   !proto.smartcore.bos.ClearPriorityValueRequest,
 *   !proto.smartcore.bos.ClearPriorityValueResponse>}
 */
const methodDescriptor_PriorityApi_ClearPriorityEntry = new grpc.web.MethodDescriptor(
  '/smartcore.bos.PriorityApi/ClearPriorityEntry',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.ClearPriorityValueRequest,
  proto.smartcore.bos.ClearPriorityValueResponse,
  /**
   * @param {!proto.smartcore.bos.ClearPriorityValueRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ClearPriorityValueResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.ClearPriorityValueRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.ClearPriorityValueResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.ClearPriorityValueResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.PriorityApiClient.prototype.clearPriorityEntry =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.PriorityApi/ClearPriorityEntry',
      request,
      metadata || {},
      methodDescriptor_PriorityApi_ClearPriorityEntry,
      callback);
};


/**
 * @param {!proto.smartcore.bos.ClearPriorityValueRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.ClearPriorityValueResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.PriorityApiPromiseClient.prototype.clearPriorityEntry =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.PriorityApi/ClearPriorityEntry',
      request,
      metadata || {},
      methodDescriptor_PriorityApi_ClearPriorityEntry);
};


module.exports = proto.smartcore.bos;

