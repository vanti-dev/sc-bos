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

const proto = {};
proto.vanti = {};
proto.vanti.bsp = {};
proto.vanti.bsp.ew = require('./test_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.vanti.bsp.ew.TestApiClient =
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
proto.vanti.bsp.ew.TestApiPromiseClient =
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
 *   !proto.vanti.bsp.ew.GetTestRequest,
 *   !proto.vanti.bsp.ew.Test>}
 */
const methodDescriptor_TestApi_GetTest = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.TestApi/GetTest',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.GetTestRequest,
    proto.vanti.bsp.ew.Test,
    /**
     * @param {!proto.vanti.bsp.ew.GetTestRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.Test.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.GetTestRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.Test)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.Test>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.TestApiClient.prototype.getTest =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.TestApi/GetTest',
          request,
          metadata || {},
          methodDescriptor_TestApi_GetTest,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.GetTestRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.Test>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.TestApiPromiseClient.prototype.getTest =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.TestApi/GetTest',
          request,
          metadata || {},
          methodDescriptor_TestApi_GetTest);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.UpdateTestRequest,
 *   !proto.vanti.bsp.ew.Test>}
 */
const methodDescriptor_TestApi_UpdateTest = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.TestApi/UpdateTest',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.UpdateTestRequest,
    proto.vanti.bsp.ew.Test,
    /**
     * @param {!proto.vanti.bsp.ew.UpdateTestRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.Test.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.UpdateTestRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.Test)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.Test>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.TestApiClient.prototype.updateTest =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.TestApi/UpdateTest',
          request,
          metadata || {},
          methodDescriptor_TestApi_UpdateTest,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.UpdateTestRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.Test>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.TestApiPromiseClient.prototype.updateTest =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.TestApi/UpdateTest',
          request,
          metadata || {},
          methodDescriptor_TestApi_UpdateTest);
    };


module.exports = proto.vanti.bsp.ew;

