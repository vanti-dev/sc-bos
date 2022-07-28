/**
 * @fileoverview gRPC-Web generated client stub for vanti.ew_auth_poc
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.vanti = {};
proto.vanti.ew_auth_poc = require('./test_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.vanti.ew_auth_poc.TestApiClient =
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
proto.vanti.ew_auth_poc.TestApiPromiseClient =
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
 *   !proto.vanti.ew_auth_poc.GetTestRequest,
 *   !proto.vanti.ew_auth_poc.Test>}
 */
const methodDescriptor_TestApi_GetTest = new grpc.web.MethodDescriptor(
  '/vanti.ew_auth_poc.TestApi/GetTest',
  grpc.web.MethodType.UNARY,
  proto.vanti.ew_auth_poc.GetTestRequest,
  proto.vanti.ew_auth_poc.Test,
  /**
   * @param {!proto.vanti.ew_auth_poc.GetTestRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.vanti.ew_auth_poc.Test.deserializeBinary
);


/**
 * @param {!proto.vanti.ew_auth_poc.GetTestRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.ew_auth_poc.Test)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.ew_auth_poc.Test>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.ew_auth_poc.TestApiClient.prototype.getTest =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/vanti.ew_auth_poc.TestApi/GetTest',
      request,
      metadata || {},
      methodDescriptor_TestApi_GetTest,
      callback);
};


/**
 * @param {!proto.vanti.ew_auth_poc.GetTestRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.ew_auth_poc.Test>}
 *     Promise that resolves to the response
 */
proto.vanti.ew_auth_poc.TestApiPromiseClient.prototype.getTest =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/vanti.ew_auth_poc.TestApi/GetTest',
      request,
      metadata || {},
      methodDescriptor_TestApi_GetTest);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.ew_auth_poc.UpdateTestRequest,
 *   !proto.vanti.ew_auth_poc.Test>}
 */
const methodDescriptor_TestApi_UpdateTest = new grpc.web.MethodDescriptor(
  '/vanti.ew_auth_poc.TestApi/UpdateTest',
  grpc.web.MethodType.UNARY,
  proto.vanti.ew_auth_poc.UpdateTestRequest,
  proto.vanti.ew_auth_poc.Test,
  /**
   * @param {!proto.vanti.ew_auth_poc.UpdateTestRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.vanti.ew_auth_poc.Test.deserializeBinary
);


/**
 * @param {!proto.vanti.ew_auth_poc.UpdateTestRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.ew_auth_poc.Test)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.ew_auth_poc.Test>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.ew_auth_poc.TestApiClient.prototype.updateTest =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/vanti.ew_auth_poc.TestApi/UpdateTest',
      request,
      metadata || {},
      methodDescriptor_TestApi_UpdateTest,
      callback);
};


/**
 * @param {!proto.vanti.ew_auth_poc.UpdateTestRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.ew_auth_poc.Test>}
 *     Promise that resolves to the response
 */
proto.vanti.ew_auth_poc.TestApiPromiseClient.prototype.updateTest =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/vanti.ew_auth_poc.TestApi/UpdateTest',
      request,
      metadata || {},
      methodDescriptor_TestApi_UpdateTest);
};


module.exports = proto.vanti.ew_auth_poc;

