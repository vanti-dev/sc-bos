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

const proto = {};
proto.smartcore = {};
proto.smartcore.bos = require('./nodes_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.NodeApiClient =
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
proto.smartcore.bos.NodeApiPromiseClient =
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
 *   !proto.smartcore.bos.GetNodeRegistrationRequest,
 *   !proto.smartcore.bos.NodeRegistration>}
 */
const methodDescriptor_NodeApi_GetNodeRegistration = new grpc.web.MethodDescriptor(
  '/smartcore.bos.NodeApi/GetNodeRegistration',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.GetNodeRegistrationRequest,
  proto.smartcore.bos.NodeRegistration,
  /**
   * @param {!proto.smartcore.bos.GetNodeRegistrationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.NodeRegistration.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.GetNodeRegistrationRequest,
 *   !proto.smartcore.bos.NodeRegistration>}
 */
const methodInfo_NodeApi_GetNodeRegistration = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.NodeRegistration,
  /**
   * @param {!proto.smartcore.bos.GetNodeRegistrationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.NodeRegistration.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.GetNodeRegistrationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.NodeRegistration)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.NodeRegistration>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.NodeApiClient.prototype.getNodeRegistration =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.NodeApi/GetNodeRegistration',
      request,
      metadata || {},
      methodDescriptor_NodeApi_GetNodeRegistration,
      callback);
};


/**
 * @param {!proto.smartcore.bos.GetNodeRegistrationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.NodeRegistration>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.NodeApiPromiseClient.prototype.getNodeRegistration =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.NodeApi/GetNodeRegistration',
      request,
      metadata || {},
      methodDescriptor_NodeApi_GetNodeRegistration);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.CreateNodeRegistrationRequest,
 *   !proto.smartcore.bos.NodeRegistration>}
 */
const methodDescriptor_NodeApi_CreateNodeRegistration = new grpc.web.MethodDescriptor(
  '/smartcore.bos.NodeApi/CreateNodeRegistration',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.CreateNodeRegistrationRequest,
  proto.smartcore.bos.NodeRegistration,
  /**
   * @param {!proto.smartcore.bos.CreateNodeRegistrationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.NodeRegistration.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.CreateNodeRegistrationRequest,
 *   !proto.smartcore.bos.NodeRegistration>}
 */
const methodInfo_NodeApi_CreateNodeRegistration = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.NodeRegistration,
  /**
   * @param {!proto.smartcore.bos.CreateNodeRegistrationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.NodeRegistration.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.CreateNodeRegistrationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.NodeRegistration)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.NodeRegistration>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.NodeApiClient.prototype.createNodeRegistration =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.NodeApi/CreateNodeRegistration',
      request,
      metadata || {},
      methodDescriptor_NodeApi_CreateNodeRegistration,
      callback);
};


/**
 * @param {!proto.smartcore.bos.CreateNodeRegistrationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.NodeRegistration>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.NodeApiPromiseClient.prototype.createNodeRegistration =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.NodeApi/CreateNodeRegistration',
      request,
      metadata || {},
      methodDescriptor_NodeApi_CreateNodeRegistration);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.ListNodeRegistrationsRequest,
 *   !proto.smartcore.bos.ListNodeRegistrationsResponse>}
 */
const methodDescriptor_NodeApi_ListNodeRegistrations = new grpc.web.MethodDescriptor(
  '/smartcore.bos.NodeApi/ListNodeRegistrations',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.ListNodeRegistrationsRequest,
  proto.smartcore.bos.ListNodeRegistrationsResponse,
  /**
   * @param {!proto.smartcore.bos.ListNodeRegistrationsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListNodeRegistrationsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.ListNodeRegistrationsRequest,
 *   !proto.smartcore.bos.ListNodeRegistrationsResponse>}
 */
const methodInfo_NodeApi_ListNodeRegistrations = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.ListNodeRegistrationsResponse,
  /**
   * @param {!proto.smartcore.bos.ListNodeRegistrationsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.ListNodeRegistrationsResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.ListNodeRegistrationsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.ListNodeRegistrationsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.ListNodeRegistrationsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.NodeApiClient.prototype.listNodeRegistrations =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.NodeApi/ListNodeRegistrations',
      request,
      metadata || {},
      methodDescriptor_NodeApi_ListNodeRegistrations,
      callback);
};


/**
 * @param {!proto.smartcore.bos.ListNodeRegistrationsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.ListNodeRegistrationsResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.NodeApiPromiseClient.prototype.listNodeRegistrations =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.NodeApi/ListNodeRegistrations',
      request,
      metadata || {},
      methodDescriptor_NodeApi_ListNodeRegistrations);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.TestNodeCommunicationRequest,
 *   !proto.smartcore.bos.TestNodeCommunicationResponse>}
 */
const methodDescriptor_NodeApi_TestNodeCommunication = new grpc.web.MethodDescriptor(
  '/smartcore.bos.NodeApi/TestNodeCommunication',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.TestNodeCommunicationRequest,
  proto.smartcore.bos.TestNodeCommunicationResponse,
  /**
   * @param {!proto.smartcore.bos.TestNodeCommunicationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.TestNodeCommunicationResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.smartcore.bos.TestNodeCommunicationRequest,
 *   !proto.smartcore.bos.TestNodeCommunicationResponse>}
 */
const methodInfo_NodeApi_TestNodeCommunication = new grpc.web.AbstractClientBase.MethodInfo(
  proto.smartcore.bos.TestNodeCommunicationResponse,
  /**
   * @param {!proto.smartcore.bos.TestNodeCommunicationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.TestNodeCommunicationResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.TestNodeCommunicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.smartcore.bos.TestNodeCommunicationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.TestNodeCommunicationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.NodeApiClient.prototype.testNodeCommunication =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.NodeApi/TestNodeCommunication',
      request,
      metadata || {},
      methodDescriptor_NodeApi_TestNodeCommunication,
      callback);
};


/**
 * @param {!proto.smartcore.bos.TestNodeCommunicationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.TestNodeCommunicationResponse>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.NodeApiPromiseClient.prototype.testNodeCommunication =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.NodeApi/TestNodeCommunication',
      request,
      metadata || {},
      methodDescriptor_NodeApi_TestNodeCommunication);
};


module.exports = proto.smartcore.bos;

