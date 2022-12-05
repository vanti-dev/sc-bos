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
proto.vanti.bsp.ew = require('./nodes_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.vanti.bsp.ew.NodeApiClient =
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
proto.vanti.bsp.ew.NodeApiPromiseClient =
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
 *   !proto.vanti.bsp.ew.GetNodeRegistrationRequest,
 *   !proto.vanti.bsp.ew.NodeRegistration>}
 */
const methodDescriptor_NodeApi_GetNodeRegistration = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.NodeApi/GetNodeRegistration',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.GetNodeRegistrationRequest,
    proto.vanti.bsp.ew.NodeRegistration,
    /**
     * @param {!proto.vanti.bsp.ew.GetNodeRegistrationRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.NodeRegistration.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.GetNodeRegistrationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.NodeRegistration)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.NodeRegistration>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.NodeApiClient.prototype.getNodeRegistration =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.NodeApi/GetNodeRegistration',
          request,
          metadata || {},
          methodDescriptor_NodeApi_GetNodeRegistration,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.GetNodeRegistrationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.NodeRegistration>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.NodeApiPromiseClient.prototype.getNodeRegistration =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.NodeApi/GetNodeRegistration',
          request,
          metadata || {},
          methodDescriptor_NodeApi_GetNodeRegistration);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.CreateNodeRegistrationRequest,
 *   !proto.vanti.bsp.ew.NodeRegistration>}
 */
const methodDescriptor_NodeApi_CreateNodeRegistration = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.NodeApi/CreateNodeRegistration',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.CreateNodeRegistrationRequest,
    proto.vanti.bsp.ew.NodeRegistration,
    /**
     * @param {!proto.vanti.bsp.ew.CreateNodeRegistrationRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.NodeRegistration.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.CreateNodeRegistrationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.NodeRegistration)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.NodeRegistration>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.NodeApiClient.prototype.createNodeRegistration =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.NodeApi/CreateNodeRegistration',
          request,
          metadata || {},
          methodDescriptor_NodeApi_CreateNodeRegistration,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.CreateNodeRegistrationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.NodeRegistration>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.NodeApiPromiseClient.prototype.createNodeRegistration =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.NodeApi/CreateNodeRegistration',
          request,
          metadata || {},
          methodDescriptor_NodeApi_CreateNodeRegistration);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.ListNodeRegistrationsRequest,
 *   !proto.vanti.bsp.ew.ListNodeRegistrationsResponse>}
 */
const methodDescriptor_NodeApi_ListNodeRegistrations = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.NodeApi/ListNodeRegistrations',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.ListNodeRegistrationsRequest,
    proto.vanti.bsp.ew.ListNodeRegistrationsResponse,
    /**
     * @param {!proto.vanti.bsp.ew.ListNodeRegistrationsRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.ListNodeRegistrationsResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.ListNodeRegistrationsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.ListNodeRegistrationsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.ListNodeRegistrationsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.NodeApiClient.prototype.listNodeRegistrations =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.NodeApi/ListNodeRegistrations',
          request,
          metadata || {},
          methodDescriptor_NodeApi_ListNodeRegistrations,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.ListNodeRegistrationsRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.ListNodeRegistrationsResponse>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.NodeApiPromiseClient.prototype.listNodeRegistrations =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.NodeApi/ListNodeRegistrations',
          request,
          metadata || {},
          methodDescriptor_NodeApi_ListNodeRegistrations);
    };


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.vanti.bsp.ew.TestNodeCommunicationRequest,
 *   !proto.vanti.bsp.ew.TestNodeCommunicationResponse>}
 */
const methodDescriptor_NodeApi_TestNodeCommunication = new grpc.web.MethodDescriptor(
    '/vanti.bsp.ew.NodeApi/TestNodeCommunication',
    grpc.web.MethodType.UNARY,
    proto.vanti.bsp.ew.TestNodeCommunicationRequest,
    proto.vanti.bsp.ew.TestNodeCommunicationResponse,
    /**
     * @param {!proto.vanti.bsp.ew.TestNodeCommunicationRequest} request
     * @return {!Uint8Array}
     */
    function(request) {
      return request.serializeBinary();
    },
    proto.vanti.bsp.ew.TestNodeCommunicationResponse.deserializeBinary
);


/**
 * @param {!proto.vanti.bsp.ew.TestNodeCommunicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.vanti.bsp.ew.TestNodeCommunicationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.vanti.bsp.ew.TestNodeCommunicationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.vanti.bsp.ew.NodeApiClient.prototype.testNodeCommunication =
    function(request, metadata, callback) {
      return this.client_.rpcCall(this.hostname_ +
          '/vanti.bsp.ew.NodeApi/TestNodeCommunication',
          request,
          metadata || {},
          methodDescriptor_NodeApi_TestNodeCommunication,
          callback);
    };


/**
 * @param {!proto.vanti.bsp.ew.TestNodeCommunicationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.vanti.bsp.ew.TestNodeCommunicationResponse>}
 *     Promise that resolves to the response
 */
proto.vanti.bsp.ew.NodeApiPromiseClient.prototype.testNodeCommunication =
    function(request, metadata) {
      return this.client_.unaryCall(this.hostname_ +
          '/vanti.bsp.ew.NodeApi/TestNodeCommunication',
          request,
          metadata || {},
          methodDescriptor_NodeApi_TestNodeCommunication);
    };


module.exports = proto.vanti.bsp.ew;

