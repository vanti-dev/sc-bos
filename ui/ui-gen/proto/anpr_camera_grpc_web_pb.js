/**
 * @fileoverview gRPC-Web generated client stub for smartcore.bos
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.5.0
// 	protoc              v5.29.3
// source: anpr_camera.proto


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_field_mask_pb = require('google-protobuf/google/protobuf/field_mask_pb.js')

var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')
const proto = {};
proto.smartcore = {};
proto.smartcore.bos = require('./anpr_camera_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.smartcore.bos.AnprCameraApiClient =
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
proto.smartcore.bos.AnprCameraApiPromiseClient =
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
 *   !proto.smartcore.bos.GetLastEventRequest,
 *   !proto.smartcore.bos.AnprEvent>}
 */
const methodDescriptor_AnprCameraApi_GetEvent = new grpc.web.MethodDescriptor(
  '/smartcore.bos.AnprCameraApi/GetEvent',
  grpc.web.MethodType.UNARY,
  proto.smartcore.bos.GetLastEventRequest,
  proto.smartcore.bos.AnprEvent,
  /**
   * @param {!proto.smartcore.bos.GetLastEventRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.AnprEvent.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.GetLastEventRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.smartcore.bos.AnprEvent)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.AnprEvent>|undefined}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AnprCameraApiClient.prototype.getEvent =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/smartcore.bos.AnprCameraApi/GetEvent',
      request,
      metadata || {},
      methodDescriptor_AnprCameraApi_GetEvent,
      callback);
};


/**
 * @param {!proto.smartcore.bos.GetLastEventRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.smartcore.bos.AnprEvent>}
 *     Promise that resolves to the response
 */
proto.smartcore.bos.AnprCameraApiPromiseClient.prototype.getEvent =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/smartcore.bos.AnprCameraApi/GetEvent',
      request,
      metadata || {},
      methodDescriptor_AnprCameraApi_GetEvent);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.smartcore.bos.PullEventsRequest,
 *   !proto.smartcore.bos.PullEventsResponse>}
 */
const methodDescriptor_AnprCameraApi_PullEvents = new grpc.web.MethodDescriptor(
  '/smartcore.bos.AnprCameraApi/PullEvents',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.smartcore.bos.PullEventsRequest,
  proto.smartcore.bos.PullEventsResponse,
  /**
   * @param {!proto.smartcore.bos.PullEventsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.smartcore.bos.PullEventsResponse.deserializeBinary
);


/**
 * @param {!proto.smartcore.bos.PullEventsRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullEventsResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AnprCameraApiClient.prototype.pullEvents =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.AnprCameraApi/PullEvents',
      request,
      metadata || {},
      methodDescriptor_AnprCameraApi_PullEvents);
};


/**
 * @param {!proto.smartcore.bos.PullEventsRequest} request The request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.smartcore.bos.PullEventsResponse>}
 *     The XHR Node Readable Stream
 */
proto.smartcore.bos.AnprCameraApiPromiseClient.prototype.pullEvents =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/smartcore.bos.AnprCameraApi/PullEvents',
      request,
      metadata || {},
      methodDescriptor_AnprCameraApi_PullEvents);
};


module.exports = proto.smartcore.bos;

